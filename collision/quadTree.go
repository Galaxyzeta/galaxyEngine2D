package collision

import (
	"sync/atomic"

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

type QuadTree struct {
	root        *QTreeNode        // store items inside the area
	area        physics.Rectangle // initial Qtree management area
	loadFactor  int
	minDivision float64
	looseOffset float64 // once a collider entered a cell, in how much offset to determine a collider has left the original cell.
	mu          *lock.SpinLock
}

type QTreeNode struct {
	id          int64
	items       []*component.PolygonCollider
	inlineItems []*component.PolygonCollider // inlineItems stores items that is actually on the boundary
	children    []*QTreeNode                 // points to 4 sub dimensions
	parent      *QTreeNode
	area        physics.Rectangle
	minDivision float64
	loadFactor  int // how many items can be held at most in this node
}

type QTreeTraverseFunc func(*component.PolygonCollider, *QTreeNode, AreaType, int) bool

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
	}
}

func (qt QTreeNode) GetItems() (ret []*component.PolygonCollider) {
	ret = append(ret, qt.items...)
	ret = append(ret, qt.inlineItems...)
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
	qt.root.doTraverse(fn, qt.area)
}

// TraverseNeedLock is a safer version of traverse.
// When you're doing a traverse in rendering routine, you must use this.
func (qt *QuadTree) TraverseWithLock(fn QTreeTraverseFunc) {
	qt.mu.Lock()
	defer qt.mu.Unlock()
	if qt.root == nil {
		return
	}
	qt.root.doTraverse(fn, qt.area)
}

func (qt *QuadTree) Insert(collider *component.PolygonCollider) {
	if qt.root == nil {
		qt.root = &QTreeNode{
			items:       []*component.PolygonCollider{},
			area:        qt.area,
			loadFactor:  qt.loadFactor,
			minDivision: qt.minDivision,
		}
	}
	// judge fully outside of maintainance area.
	if !collider.Collider.GetBoundingBox().ToRectangle().Intersect(&qt.area) {
		qt.root.inlineItems = append(qt.root.inlineItems, collider)
		return
	}
	qt.root.doInsert(collider)
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
		logger.GlobalLogger.Debug("trigger node merge")
		// children' items move to parent
		for _, eachChild := range parent.children {
			if eachChild == nil {
				continue
			}
			parent.items = append(parent.items, eachChild.items...)
			// no need to handle inline items, because leaf node has no inline items.
		}
		if parent.parent != nil {
			// parent's inline items move to grandparent.
			parent.parent.inlineItems = append(parent.inlineItems, parent.parent.inlineItems...)
		}
		// abandon all childs
		parent.children = nil
	}
}

func (qt *QTreeNode) Delete(collider *component.PolygonCollider) {
	for idx, item := range qt.items {
		if item == collider {
			qt.items = doDeleteFromArray(idx, qt.items)
			qt.tryNodeMerge()
			return
		}
	}
	for idx, item := range qt.inlineItems {
		if item == collider {
			qt.inlineItems = doDeleteFromArray(idx, qt.inlineItems)
			return
		}
	}
}

func doDeleteFromArray(index int, arr []*component.PolygonCollider) (ret []*component.PolygonCollider) {
	ret = append(arr[:index], arr[index+1:]...)
	return ret
}

func (qt *QuadTree) QueryByPoint(position linalg.Vector2f64) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	qt.root.doQuery(position, &result)
	return result
}

func (qt *QuadTree) QueryByRect(rect physics.Rectangle) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	for _, qtnode := range qt.root.children {
		qtnode.doQueryByRect(rect, &result)
	}
	return result
}

func (qt *QuadTree) QueryByRay(r physics.Ray) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	for idx, qtnode := range qt.root.children {
		qtnode.doQueryByRay(r, boundaryDivision(qt.area, idx), &result)
	}
	return result
}

func (qt *QTreeNode) doQueryByRect(rect physics.Rectangle, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	if rect.IntersectWithRectangle(qt.GetArea()) {
		*result = append(*result, qt.GetItems()...)
	}
	for _, childNode := range qt.children {
		childNode.doQueryByRect(rect, result)
	}
}

func (qt *QTreeNode) doQueryByRay(r physics.Ray, area physics.Rectangle, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	if r.IntersectPolygon(area.ToPolygon()) {
		if len(qt.children) > 0 {
			for idx, qtnode := range qt.children {
				qtnode.doQueryByRay(r, boundaryDivision(area, idx), result)
			}
		} else {
			*result = append(*result, qt.GetItems()...)
		}
	}
}

func (qt *QTreeNode) doTraverse(fn QTreeTraverseFunc, boundary physics.Rectangle) {
	if qt == nil {
		return
	}
	if len(qt.children) > 0 {
		for i, elem := range qt.children {
			elem.doTraverse(fn, boundaryDivision(boundary, i))
		}
	}
	for idx, item := range qt.items {
		if fn(item, qt, Normal, idx) {
			return
		}
	}
	for idx, item := range qt.inlineItems {
		if fn(item, qt, Inline, idx) {
			return
		}
	}
}

func (qt *QTreeNode) doQuery(position linalg.Vector2f64, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	// no matter it is a leaf node or not, add all is item into the result.
	// because parental nodes stores colliders that are exactly at boundary edges.
	*result = append(*result, qt.GetItems()...)

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
			qt.children[Section1].doQuery(position, result)
			return
		}
		qt.children[Section4].doQuery(position, result)
		return
	}
	if yPos {
		qt.children[Section2].doQuery(position, result)
		return
	}
	qt.children[Section3].doQuery(position, result)
}

func (qt *QTreeNode) doInsert(collider *component.PolygonCollider) {
	if qt == nil {
		return
	}
	if qt.children == nil {
		logger.GlobalLogger.Debugf("normal insert, %v, area = %v", collider.Collider.GetWorldVertices(), qt.area)
		qt.items = append(qt.items, collider)
		if len(qt.items) > qt.loadFactor {
			if qt.area.Height > qt.minDivision && qt.area.Width > qt.minDivision {
				// split
				qt.children = []*QTreeNode{
					NewQTreeNode(qt, Section1),
					NewQTreeNode(qt, Section2),
					NewQTreeNode(qt, Section3),
					NewQTreeNode(qt, Section4),
				}
				logger.GlobalLogger.Debugf("trigger split, area = %v", qt.area)
				for _, item := range qt.items {
					qt.insertIntoChild0(item)
				}
				qt.items = []*component.PolygonCollider{}
			}
		}
		return
	}
	// insert into its descendants.
	qt.insertIntoChild0(collider)
}

func (qt *QTreeNode) insertIntoChild0(collider *component.PolygonCollider) {

	whichSection := qt.GetIntersectedSection(collider.Collider.GetBoundingBox())
	if whichSection == Overlap {
		// overlap, insert into current node's inlineItem map
		qt.inlineItems = append(qt.inlineItems, collider)
		logger.GlobalLogger.Debugf("overlap, insert into inline, %v, area = %v", collider.Collider.GetWorldVertices(), qt.area)
	} else if whichSection == Overflow {
		// overflow
		if qt.parent != nil {
			panic("should not happen")
		}
		qt.inlineItems = append(qt.inlineItems, collider)
		logger.GlobalLogger.Debugf("overflow, insert into inline, %v, area = %v", collider.Collider.GetWorldVertices(), qt.area)
	} else {
		// normal insert, no overlap, insert into children
		logger.GlobalLogger.Debugf("searching, area = %v", qt.area)
		qt.children[whichSection].doInsert(collider)
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
