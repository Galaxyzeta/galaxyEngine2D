package collision

import (
	"sync/atomic"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/infra/concurrency/lock"
	"galaxyzeta.io/engine/infra/logger"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

const (
	Section1 = iota
	Section2
	Section3
	Section4
	Overlap  = -2
	Overflow = -3
)

type AreaType int8

const (
	Normal AreaType = iota
	Inline
)

var idGenerator int64
var qtlogger *logger.Logger = logger.New("QuadTree")

func init() {
	// qtlogger.Disable()
}

type QuadTree struct {
	root        *QTreeNode        // store items inside the area
	area        physics.Rectangle // initial Qtree management area
	loadFactor  int
	minDivision float64
	looseOffset float64 // once a collider entered a cell, in how much offset to determine a collider has left the original cell.
	mu          *lock.SpinLock
	lookup      map[base.IGameObject2D]*QTreeNode
}

type QTreeNode struct {
	id            int64
	items         []*component.PolygonCollider
	inlineItems   []*component.PolygonCollider // inlineItems stores items that is actually on the boundary
	inactiveItems []*component.PolygonCollider // inactiveItems stors deactivated items
	children      []*QTreeNode                 // points to 4 sub dimensions
	parent        *QTreeNode
	area          physics.Rectangle
	minDivision   float64
	loadFactor    int       // how many items can be held at most in this node
	quadTree      *QuadTree // always point to the root
}

type QTreeTraverseFunc func(*component.PolygonCollider, *QTreeNode, AreaType, int) bool

type qtreeNodeCollectorFunc func() []*component.PolygonCollider

type QueryMode int8

const (
	ActiveOnly   QueryMode = 0
	InactiveOnly QueryMode = 1
	All          QueryMode = 2
)

func (qt *QTreeNode) collectorFxCollectActive() []*component.PolygonCollider {
	return qt.GetActiveItems()
}

func (qt *QTreeNode) collectorFxCollectInactive() []*component.PolygonCollider {
	return qt.GetInactiveItems()
}

func (qt *QTreeNode) collectorFxCollectAll() []*component.PolygonCollider {
	return qt.GetAllItems()
}

func (qt *QTreeNode) chooseCollectorFx(mode QueryMode) qtreeNodeCollectorFunc {
	switch mode {
	case ActiveOnly:
		return qt.collectorFxCollectActive
	case InactiveOnly:
		return qt.collectorFxCollectInactive
	case All:
		return qt.collectorFxCollectAll
	default:
		return qt.collectorFxCollectAll
	}
}

func NewQuadTree(maintainanceArea physics.Rectangle, loadFactor int, minDivision float64) *QuadTree {
	if minDivision < 32 {
		panic("cannot have minDivision less than 32")
	}
	return &QuadTree{
		area:        maintainanceArea,
		loadFactor:  loadFactor,
		minDivision: minDivision,
		looseOffset: 0,
		mu:          &lock.SpinLock{},
		lookup:      make(map[base.IGameObject2D]*QTreeNode),
	}
}

func NewQTreeNode(parent *QTreeNode, section int) *QTreeNode {
	return &QTreeNode{
		id:          atomic.AddInt64(&idGenerator, 1),
		items:       []*component.PolygonCollider{},
		parent:      parent,
		area:        boundaryDivision(parent.area, section),
		loadFactor:  parent.loadFactor,
		minDivision: parent.minDivision,
		quadTree:    parent.quadTree,
	}
}

func (qt QTreeNode) GetActiveItems() (ret []*component.PolygonCollider) {
	ret = append(ret, qt.items...)
	ret = append(ret, qt.inlineItems...)
	return ret
}

func (qt QTreeNode) GetInactiveItems() (ret []*component.PolygonCollider) {
	ret = append(ret, qt.inactiveItems...)
	return ret
}

func (qt QTreeNode) GetAllItems() (ret []*component.PolygonCollider) {
	ret = append(ret, qt.items...)
	ret = append(ret, qt.inlineItems...)
	ret = append(ret, qt.inactiveItems...)
	return ret
}

func (qt QTreeNode) GetChildren() []*QTreeNode {
	return qt.children
}

func (qt QTreeNode) GetArea() physics.Rectangle {
	return qt.area
}

func (qt QuadTree) GetLooseOffset() float64 {
	return qt.looseOffset
}

// Traverse is an unsafe operation. Use this can bring better performance, but you need to ensure
// the whole tree is not concurrently visited.
func (qt *QuadTree) Traverse(fn QTreeTraverseFunc) {
	if qt.root == nil {
		return
	}
	qt.root.doTraverse(fn)
}

// TraverseNeedLock is a safer version of traverse.
// When you're doing a traverse in rendering routine, you must use this.
func (qt *QuadTree) TraverseWithLock(fn QTreeTraverseFunc) {
	// qtlogger.Debug("--begin--")
	// timer := time.Now()
	qt.mu.Lock()
	defer qt.mu.Unlock()
	if qt.root == nil {
		return
	}
	qt.root.doTraverse(fn)
	// qtlogger.Debug("--end--")
	// qtlogger.Debugf("cost = %v", time.Since(timer))
}

func (qt *QuadTree) Insert(collider *component.PolygonCollider) {
	if qt.root == nil {
		qt.root = &QTreeNode{
			items:       []*component.PolygonCollider{},
			area:        qt.area,
			loadFactor:  qt.loadFactor,
			minDivision: qt.minDivision,
			quadTree:    qt,
		}
	}
	// judge fully outside of maintainance area.
	qt.mu.Lock()
	if !collider.Collider.GetBoundingBox().ToRectangle().Intersect(&qt.area) {
		qt.root.doInsertInline(collider)
		qt.mu.Unlock()
		return
	}
	qt.root.insertRecursively(collider)
	qt.mu.Unlock()
}

func (qt *QuadTree) Deactivate(collider *component.PolygonCollider) bool {
	node, ok := qt.lookup[collider.I()]
	if !ok {
		panic("failed to find correlated node in look up table")
	}
	node.quadTree.mu.Lock()
	idx := node.searchFromNormal(collider)
	if idx >= 0 {
		node.deleteFromNormal(idx)
		node.doInsertInactive(collider)
		node.quadTree.mu.Unlock()
		return true
	}
	idx = node.searchFromInline(collider)
	if idx >= 0 {
		node.deleteFromInline(idx)
		node.doInsertInactive(collider)
		node.quadTree.mu.Unlock()
		return true
	}
	node.quadTree.mu.Unlock()
	return false
}

func (qt *QuadTree) Activate(collider *component.PolygonCollider) bool {
	node, ok := qt.lookup[collider.I()]
	if !ok {
		panic("failed to find correlated node in look up table")
	}
	node.quadTree.mu.Lock()
	idx := node.searchFromInactive(collider)
	if idx >= 0 {
		node.deleteFromInactive(idx)
		node.quadTree.mu.Unlock()
		qt.Insert(collider)
		return true
	} else {
		node.quadTree.mu.Unlock()
	}
	return false
}

func (qt *QTreeNode) tryNodeMerge() {
	if qt.parent == nil {
		return
	}
	itemsCnt := 0
	parent := qt.parent
	for _, eachChild := range parent.children {
		if eachChild == nil {
			continue
		}
		itemsCnt += len(eachChild.items)
		if eachChild.children != nil {
			// cannot merge 4 nodes which at least 1 of them contains sub tree.
			return
		}
	}
	if itemsCnt < (qt.parent.loadFactor >> 1) {
		qtlogger.Debug("trigger node merge")
		// children' items move to parent
		for _, eachChild := range parent.children {
			if eachChild == nil {
				continue
			}

			parent.items = append(parent.items, eachChild.items...)
			parent.inactiveItems = append(parent.items, eachChild.inactiveItems...)

			quadTree := eachChild.quadTree
			for idx, item := range eachChild.items {
				eachChild.items[idx] = nil
				quadTree.setLookup(item.I(), parent, "try node merge 1")

			}

			for idx, item := range eachChild.inactiveItems {
				eachChild.inactiveItems[idx] = nil
				quadTree.setLookup(item.I(), parent, "try node merge 2")
			}

			// no need to handle inline items, because leaf node has no inline items.
		}
		if parent.parent != nil {
			// parent's inline items move to grandparent.
			quadTree := parent.quadTree
			parent.parent.inlineItems = append(parent.inlineItems, parent.parent.inlineItems...)
			for idx, item := range parent.inlineItems {
				parent.inlineItems[idx] = nil
				quadTree.setLookup(item.I(), parent.parent, "try node merge 3")
			}
		}
		// abandon all childs
		parent.children = nil
	}
}

func (qt *QTreeNode) Delete(collider *component.PolygonCollider) {
	qt.quadTree.mu.Lock()
	qt.UnsafeDelete(collider)
	qt.quadTree.mu.Unlock()
}

// UnsafeDelete deletes item without acquiring lock.
// Don't use this if you don't know what you're doing.
// Use Delete() instead.
func (qt *QTreeNode) UnsafeDelete(collider *component.PolygonCollider) {
	for idx, item := range qt.items {
		if item == collider {
			qt.deleteFromNormal(idx)
			return
		}
	}
	for idx, item := range qt.inlineItems {
		if item == collider {
			qt.deleteFromInline(idx)
			return
		}
	}
}

func (qt *QTreeNode) deleteFromNormal(idx int) {
	qt.quadTree.deleteLookup(qt.items[idx].I(), "deleteFromNormal")
	qt.items = doDeleteFromArray(idx, qt.items)
	qt.tryNodeMerge()
}

func (qt *QTreeNode) deleteFromInline(idx int) {
	qt.quadTree.deleteLookup(qt.inlineItems[idx].I(), "deleteFromInline")
	qt.inlineItems = doDeleteFromArray(idx, qt.inlineItems)
}

func (qt *QTreeNode) deleteFromInactive(idx int) {
	qt.quadTree.deleteLookup(qt.inactiveItems[idx].I(), "deleteFromInactive")
	qt.inactiveItems = doDeleteFromArray(idx, qt.inactiveItems)
}

func (qt *QTreeNode) searchFromNormal(collider *component.PolygonCollider) (index int) {
	return doSearchFromArray(qt.items, collider)
}

func (qt *QTreeNode) searchFromInline(collider *component.PolygonCollider) (index int) {
	return doSearchFromArray(qt.inlineItems, collider)
}

func (qt *QTreeNode) searchFromInactive(collider *component.PolygonCollider) (index int) {
	return doSearchFromArray(qt.inactiveItems, collider)
}

func doSearchFromArray(arr []*component.PolygonCollider, target *component.PolygonCollider) (index int) {
	for index = 0; index < len(arr); index++ {
		if target == arr[index] {
			return index
		}
	}
	return -1
}

func doDeleteFromArray(index int, arr []*component.PolygonCollider) (ret []*component.PolygonCollider) {
	last := len(arr) - 1
	arr[index] = arr[last]
	arr[last] = nil
	ret = arr[:last]
	return ret
}

func (qt *QuadTree) QueryByPoint(position linalg.Vector2f64, mode QueryMode) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	qt.root.doQuery(position, mode, &result)
	return result
}

func (qt *QuadTree) QueryByRect(rect physics.Rectangle, mode QueryMode) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	// query inactive
	qt.root.doQueryByRect(rect, mode, &result)
	return result
}

func (qt *QuadTree) QueryByRay(r physics.Ray, mode QueryMode) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	qt.root.doQueryByRay(r, mode, &result)
	return result
}

func (qt *QTreeNode) doQueryByRect(rect physics.Rectangle, mode QueryMode, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	collectorFunc := qt.chooseCollectorFx(mode)
	if rect.IntersectWithRectangle(qt.GetArea()) {
		*result = append(*result, collectorFunc()...)
	}
	for _, childNode := range qt.children {
		childNode.doQueryByRect(rect, mode, result)
	}
}

func (qt *QTreeNode) doQueryByRay(r physics.Ray, mode QueryMode, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	collectorFunc := qt.chooseCollectorFx(mode)
	if r.IntersectPolygon(qt.area.ToPolygon()) {
		if len(qt.children) > 0 {
			for _, qtnode := range qt.children {
				qtnode.doQueryByRay(r, mode, result)
			}
		} else {
			*result = append(*result, collectorFunc()...)
		}
	}
}

func (qt *QTreeNode) doTraverse(fn QTreeTraverseFunc) {
	if qt == nil {
		return
	}
	if len(qt.children) > 0 {
		for _, elem := range qt.children {
			elem.doTraverse(fn)
		}
	}
	for idx, item := range qt.items {
		// rec := time.Now()
		if fn(item, qt, Normal, idx) {
			return
		}
		// qtlogger.Debugf("[Normal] item = %v deltaTime = %v", item.GetIGameObject2D().GetGameObject2D().Name, time.Since(rec))

	}
	for idx, item := range qt.inlineItems {
		// rec := time.Now()
		if fn(item, qt, Inline, idx) {
			return
		}
		// qtlogger.Debugf("[Inline] item = %v deltaTime = %v", item.GetIGameObject2D().GetGameObject2D().Name, time.Since(rec))

	}
}

func (qt *QTreeNode) doQuery(position linalg.Vector2f64, mode QueryMode, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	// no matter it is a leaf node or not, add all is item into the result.
	// because parental nodes stores colliders that are exactly at boundary edges.
	collectorFunc := qt.chooseCollectorFx(mode)

	*result = append(*result, collectorFunc()...)

	if len(qt.children) == 0 {
		return
	}

	boundaryWidthHalf := qt.area.Width / 2
	boundaryHeightHalf := qt.area.Height / 2
	center := linalg.Vector2f64{
		X: qt.area.Left + boundaryWidthHalf,
		Y: qt.area.Top + boundaryHeightHalf,
	}
	xPos := position.X > center.X
	yPos := position.Y > center.Y
	if xPos {
		if yPos {
			qt.children[Section1].doQuery(position, mode, result)
			return
		}
		qt.children[Section4].doQuery(position, mode, result)
		return
	}
	if yPos {
		qt.children[Section2].doQuery(position, mode, result)
		return
	}
	qt.children[Section3].doQuery(position, mode, result)
}

func (qt *QTreeNode) insertRecursively(collider *component.PolygonCollider) {
	if qt == nil {
		return
	}
	if qt.children == nil {
		qtlogger.Debugf("normal insert, %v, area = %v", collider.Collider.GetWorldVertices(), qt.area)
		qt.doInsertNormal(collider)
		if len(qt.items) > qt.loadFactor {
			if qt.area.Height > qt.minDivision && qt.area.Width > qt.minDivision {
				// split
				qt.children = qt.newChildren()
				qtlogger.Debugf("trigger split, area = %v", qt.area)
				for _, item := range qt.items {
					qt.insertIntoChildRecursively(item)
				}
				qt.items = []*component.PolygonCollider{}
			}
		}
		return
	}
	// insert into its descendants.
	qt.insertIntoChildRecursively(collider)
}

func (qt *QTreeNode) insertIntoChildRecursively(collider *component.PolygonCollider) {

	whichSection := qt.GetIntersectedSection(collider.Collider.GetBoundingBox())

	if whichSection == Overlap {
		// overlap, insert into current node's inlineItem map
		qt.doInsertInline(collider)
		qtlogger.Debugf("overlap, insert into inline, %v, area = %v", collider.Collider.GetWorldVertices(), qt.area)
	} else if whichSection == Overflow {
		// overflow
		if qt.parent != nil {
			panic("should not happen")
		}
		qt.doInsertInline(collider)
		qtlogger.Debugf("overflow, insert into inline, %v, area = %v", collider.Collider.GetWorldVertices(), qt.area)
	} else {
		// normal insert, no overlap, insert into children
		qt.children[whichSection].insertRecursively(collider)
	}
}

func (qt *QTreeNode) doInsertInline(collider *component.PolygonCollider) {
	qt.inlineItems = append(qt.inlineItems, collider)
	qt.quadTree.setLookup(collider.I(), qt, "doInsertInline")
}

func (qt *QTreeNode) doInsertNormal(collider *component.PolygonCollider) {
	qt.items = append(qt.items, collider)
	qt.quadTree.setLookup(collider.I(), qt, "doInsertNormal")
}

func (qt *QTreeNode) doInsertInactive(collider *component.PolygonCollider) {
	qt.inactiveItems = append(qt.inactiveItems, collider)
	qt.quadTree.setLookup(collider.I(), qt, "doInsertInactive")
}

func (qt *QTreeNode) newChildren() []*QTreeNode {
	return []*QTreeNode{
		NewQTreeNode(qt, Section1),
		NewQTreeNode(qt, Section2),
		NewQTreeNode(qt, Section3),
		NewQTreeNode(qt, Section4),
	}
}

func boundaryDivision(boundary physics.Rectangle, sector int) physics.Rectangle {
	boundaryWidthHalf := boundary.Width / 2
	boundaryHeightHalf := boundary.Height / 2
	switch sector {
	case Section1:
		return physics.NewRectangle(boundary.Left+boundaryWidthHalf, boundary.Top+boundaryHeightHalf, boundaryWidthHalf, boundaryHeightHalf)
	case Section2:
		return physics.NewRectangle(boundary.Left, boundary.Top+boundaryHeightHalf, boundaryWidthHalf, boundaryHeightHalf)
	case Section3:
		return physics.NewRectangle(boundary.Left, boundary.Top, boundaryWidthHalf, boundaryHeightHalf)
	case Section4:
		return physics.NewRectangle(boundary.Left+boundaryWidthHalf, boundary.Top, boundaryWidthHalf, boundaryHeightHalf)
	}
	panic("unknown sector index")
}

// GetIntersectedSection returns which section does the given boundary box intersect with after dividing
// a parent boundary into four child boundaries.
func (qt *QTreeNode) GetIntersectedSection(bb physics.BoundingBox) int {
	boundary := qt.area
	boundaryWidthHalf := boundary.Width / 2
	boundaryHeightHalf := boundary.Height / 2
	center := linalg.Vector2f64{
		X: boundary.Left + boundaryWidthHalf,
		Y: boundary.Top + boundaryHeightHalf,
	}
	var insertedFlag [4]bool
	for i := 0; i < 4; i++ {
		xPos := bb[i].X > center.X
		yPos := bb[i].Y > center.Y
		if xPos {
			if yPos {
				insertedFlag[Section1] = true
				continue
			}
			insertedFlag[Section4] = true
			continue
		}
		if yPos {
			insertedFlag[Section2] = true
			continue
		}
		insertedFlag[Section3] = true
	}

	whichSection := -1
	for idx, section := range insertedFlag {
		if section {
			if whichSection == -1 {
				whichSection = idx
			} else {
				whichSection = Overlap
				break
			}
		}
	}
	if whichSection == -1 {
		return Overflow
	}
	return whichSection
}

func (qt *QuadTree) setLookup(iobj base.IGameObject2D, node *QTreeNode, remark ...string) {
	qtlogger.Debugf("setting lookup for %v, extra = %v", iobj.Obj().Name, remark)
	qt.lookup[iobj] = node
}

func (qt *QuadTree) deleteLookup(iobj base.IGameObject2D, remark ...string) {
	qtlogger.Debugf("deleting lookup for %v, extra = %v", iobj.Obj().Name, remark)
	delete(qt.lookup, iobj)
}
