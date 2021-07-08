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

func NewQTreeNode(parent *QTreeNode) *QTreeNode {
	return &QTreeNode{
		items:       map[*component.PolygonCollider]struct{}{},
		children:    [4]*QTreeNode{},
		parent:      parent,
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

func (qt QTreeNode) SetHasChild(hasChild bool) {
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
			hasChild:    false,
			loadFactor:  qt.loadFactor,
			minDivision: qt.minDivision,
		}
	}
	qt.root.doInsert(collider, qt.area)
}

func (qt *QuadTree) Query(position linalg.Vector2f64) map[*component.PolygonCollider]struct{} {
	res, _ := qt.root.doQuery(position, qt.area)
	return res.items
}

func (qt *QuadTree) QueryByCollider(collider *component.PolygonCollider) map[*component.PolygonCollider]struct{} {
	bb := collider.Collider.GetBoundingBox()
	bbRect := bb.ToRectangle()
	boundaries := []physics.Rectangle{}
	cnt := 0
	colliderSet := map[*component.PolygonCollider]struct{}{}
	for i := 0; i < 4; i++ {
		// if boundary already inside known, no need to search the tree again
		for _, boundary := range boundaries {
			if bbRect.InsideRectangle(&boundary) {
				continue
			}
		}

		node, boundary := qt.root.doQuery(bb[i], qt.area)
		for elem := range node.items {
			colliderSet[elem] = struct{}{}
		}
		boundaries[cnt] = boundary
	}
	return colliderSet
}

func (qt *QTreeNode) doTraverse(fn func(*component.PolygonCollider, physics.Rectangle, *QTreeNode), boundary physics.Rectangle) {
	if qt == nil {
		return
	}
	if qt.hasChild {
		for i, elem := range qt.children {
			elem.doTraverse(fn, boundaryDivision(boundary, i))
		}
	} else {
		for item := range qt.items {
			fn(item, boundary, qt)
		}
	}
}

func (qt *QTreeNode) doQuery(position linalg.Vector2f64, boundary physics.Rectangle) (*QTreeNode, physics.Rectangle) {
	if !qt.hasChild {
		return qt, boundary
	}
	boundaryWidthHalf := boundary.Width / 2
	boundaryHeightHalf := boundary.Height / 2
	center := linalg.Vector2f64{
		X: boundary.Left + boundaryWidthHalf,
		Y: boundary.Width + boundaryHeightHalf,
	}
	xPos := position.X > center.X
	yPos := position.Y > center.Y
	if xPos {
		if yPos {
			return qt.children[sector1].doQuery(position, physics.NewRectangle(center.X, boundary.Top, boundaryWidthHalf, boundaryHeightHalf))
		}
		return qt.children[sector4].doQuery(position, physics.NewRectangle(center.X, center.Y, boundaryWidthHalf, boundaryHeightHalf))
	}
	if yPos {
		return qt.children[sector2].doQuery(position, physics.NewRectangle(boundary.Left, boundary.Top, boundaryWidthHalf, boundaryHeightHalf))
	}
	return qt.children[sector3].doQuery(position, physics.NewRectangle(boundary.Left, center.Y, boundaryWidthHalf, boundaryHeightHalf))
}

func (qt *QTreeNode) doInsert(collider *component.PolygonCollider, boundary physics.Rectangle) {
	if qt == nil {
		return
	}
	if !qt.hasChild {
		qt.items[collider] = struct{}{}
		if len(qt.items) > qt.loadFactor {
			if boundary.Height > qt.minDivision && boundary.Width > qt.minDivision {
				qt.hasChild = true
				// split
				for item := range qt.items {
					delete(qt.items, item)
					qt.insertIntoChild0(item, boundary)
				}
			}
		}
		return
	}
	qt.insertIntoChild0(collider, boundary)
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

func (qt *QTreeNode) insertIntoChild0(collider *component.PolygonCollider, boundary physics.Rectangle) {
	bb := collider.Collider.GetBoundingBox()
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
				qt.insertIntoChild1(collider, sector1, boundaryDivision(boundary, sector1), &insertedFlag)
				return
			}
			qt.insertIntoChild1(collider, sector4, boundaryDivision(boundary, sector4), &insertedFlag)
			return
		}
		if yPos {
			qt.insertIntoChild1(collider, sector2, boundaryDivision(boundary, sector2), &insertedFlag)
			return
		}
		qt.insertIntoChild1(collider, sector3, boundaryDivision(boundary, sector3), &insertedFlag)
	}
}

func (qt *QTreeNode) insertIntoChild1(collider *component.PolygonCollider, sector int, rect physics.Rectangle, insertedFlag *[4]bool) {
	if insertedFlag[sector] {
		return
	}
	insertedFlag[sector] = true
	if qt.children[sector] == nil {
		qt.children[sector] = NewQTreeNode(qt) // decrease loadFactor
	}
	qt.children[sector].doInsert(collider, rect)
}
