package collision

import (
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

const (
	sector1 = iota
	sector2
	sector3
	sector4
	overlap = -2
)

type QuadTree struct {
	root         *QTreeNode                              // store items inside the area
	overflowPool map[*component.PolygonCollider]struct{} // store items outside the area
	area         physics.Rectangle                       // initial Qtree management area
	loadFactor   int
	minDivision  float64
	looseOffset  float64 // once a collider entered a cell, in how much offset to determine a collider has left the original cell.
}

type QTreeNode struct {
	items       map[*component.PolygonCollider]struct{}
	children    [4]*QTreeNode // points to 4 sub dimensions
	parent      *QTreeNode
	area        physics.Rectangle
	minDivision float64
	loadFactor  int  // how many items can be held at most in this node
	hasChild    bool // marks whether should insert element into this node, or in its children nodes.
}

func NewQuadTree(maintainanceArea physics.Rectangle, loadFactor int, minDivision float64) *QuadTree {
	if minDivision < 32 {
		panic("cannot have minDivision less than 32")
	}
	return &QuadTree{
		overflowPool: map[*component.PolygonCollider]struct{}{},
		area:         maintainanceArea,
		loadFactor:   loadFactor,
		minDivision:  minDivision,
		looseOffset:  16,
	}
}

func NewQTreeNode(parent *QTreeNode, section int) *QTreeNode {
	return &QTreeNode{
		items:       map[*component.PolygonCollider]struct{}{},
		children:    [4]*QTreeNode{},
		parent:      parent,
		area:        boundaryDivision(parent.area, section),
		hasChild:    false,
		loadFactor:  parent.loadFactor,
		minDivision: parent.minDivision,
	}
}

func (qt *QTreeNode) GetItem() map[*component.PolygonCollider]struct{} {
	return qt.items
}

func (qt QuadTree) GetLooseOffset() float64 {
	return qt.looseOffset
}

func (qt QTreeNode) GetHasChild() bool {
	return qt.hasChild
}

func (qt *QTreeNode) SetHasChild(hasChild bool) {
	qt.hasChild = hasChild
}

func (qt *QuadTree) Traverse(fn func(*component.PolygonCollider, physics.Rectangle, *QTreeNode)) {
	if qt.root == nil {
		return
	}
	qt.root.doTraverse(fn, qt.area)
}

func (qt *QuadTree) Insert(collider *component.PolygonCollider) {
	// judge fully outside of maintainance area.
	if !collider.Collider.GetBoundingBox().ToRectangle().Intersect(&qt.area) {
		qt.overflowPool[collider] = struct{}{}
		return
	}
	if qt.root == nil {
		qt.root = &QTreeNode{
			items:       map[*component.PolygonCollider]struct{}{},
			children:    [4]*QTreeNode{},
			parent:      nil,
			area:        qt.area,
			hasChild:    false,
			loadFactor:  qt.loadFactor,
			minDivision: qt.minDivision,
		}
	}
	qt.root.doInsert(collider)
}

func (qt *QuadTree) Query(position linalg.Vector2f64) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	qt.root.doQuery(position, &result)
	return result
}

func (qt *QuadTree) QueryByRay(r physics.Ray) []*component.PolygonCollider {
	result := make([]*component.PolygonCollider, 0)
	for idx, qtnode := range qt.root.children {
		qtnode.doQueryByRay(r, boundaryDivision(qt.area, idx), &result)
	}
	return result
}

func (qt *QTreeNode) doQueryByRay(r physics.Ray, area physics.Rectangle, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	if r.IntersectPolygon(area.ToPolygon()) {
		if qt.hasChild {
			for idx, qtnode := range qt.children {
				qtnode.doQueryByRay(r, boundaryDivision(area, idx), result)
			}
		} else {
			for item, _ := range qt.items {
				*result = append(*result, item)
			}
		}
	}
}

func (qt *QTreeNode) doTraverse(fn func(*component.PolygonCollider, physics.Rectangle, *QTreeNode), boundary physics.Rectangle) {
	if qt == nil {
		return
	}
	if qt.hasChild {
		for i, elem := range qt.children {
			elem.doTraverse(fn, boundaryDivision(boundary, i))
		}
	}
	for item := range qt.items {
		fn(item, boundary, qt)
	}
}

func (qt *QTreeNode) doQuery(position linalg.Vector2f64, result *[]*component.PolygonCollider) {
	if qt == nil {
		return
	}
	// no matter it is a leaf node or not, add all is item into the result.
	// because parental nodes stores colliders that are exactly at boundary edges.
	for item, _ := range qt.items {
		*result = append(*result, item)
	}

	if !qt.hasChild {
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
			qt.children[sector1].doQuery(position, result)
			return
		}
		qt.children[sector4].doQuery(position, result)
		return
	}
	if yPos {
		qt.children[sector2].doQuery(position, result)
		return
	}
	qt.children[sector3].doQuery(position, result)
}

func (qt *QTreeNode) doInsert(collider *component.PolygonCollider) {
	if qt == nil {
		return
	}
	if !qt.hasChild {
		qt.items[collider] = struct{}{}
		if len(qt.items) > qt.loadFactor {
			if qt.area.Height > qt.minDivision && qt.area.Width > qt.minDivision {
				qt.hasChild = true
				// split
				for item := range qt.items {
					delete(qt.items, item)
					qt.insertIntoChild0(item)
				}
			}
		}
		return
	}
	qt.insertIntoChild0(collider)
}

func (qt *QTreeNode) insertIntoChild0(collider *component.PolygonCollider) {

	whichSection := getIntersectedSection(collider.Collider.GetBoundingBox(), qt.area)
	if whichSection == -2 {
		// overlap, insert into current node
		qt.items[collider] = struct{}{}
	} else {
		// normal insert, no overlap, insert into children
		if qt.children[whichSection] == nil {
			qt.children[whichSection] = NewQTreeNode(qt, whichSection)
		}
		qt.children[whichSection].doInsert(collider)
	}
}

func boundaryDivision(boundary physics.Rectangle, sector int) physics.Rectangle {
	boundaryWidthHalf := boundary.Width / 2
	boundaryHeightHalf := boundary.Height / 2
	switch sector {
	case sector1:
		return physics.NewRectangle(boundary.Left+boundaryWidthHalf, boundary.Top+boundaryHeightHalf, boundaryWidthHalf, boundaryHeightHalf)
	case sector2:
		return physics.NewRectangle(boundary.Left, boundary.Top+boundaryHeightHalf, boundaryWidthHalf, boundaryHeightHalf)
	case sector3:
		return physics.NewRectangle(boundary.Left, boundary.Top, boundaryWidthHalf, boundaryHeightHalf)
	case sector4:
		return physics.NewRectangle(boundary.Left+boundaryWidthHalf, boundary.Top, boundaryWidthHalf, boundaryHeightHalf)
	}
	panic("unknown sector index")
}

// getIntersectedSection returns which section does the given boundary box intersect with after dividing
// a parent boundary into four child boundaries.
func getIntersectedSection(bb physics.BoundingBox, boundary physics.Rectangle) int {
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
				insertedFlag[sector1] = true
				continue
			}
			insertedFlag[sector4] = true
			continue
		}
		if yPos {
			insertedFlag[sector2] = true
			continue
		}
		insertedFlag[sector3] = true
	}

	whichSection := -1
	for idx, section := range insertedFlag {
		if section {
			if whichSection == -1 {
				whichSection = idx
			} else {
				whichSection = -2
				break
			}
		}
	}
	if whichSection == -1 {
		panic("failed to insert")
	}
	return whichSection
}
