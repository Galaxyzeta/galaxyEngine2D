package collision

import (
	"math"

	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

const WorldCoordMin float64 = math.MinInt64
const WorldCoordMax float64 = math.MaxInt64
const (
	sector1 = iota
	sector2
	sector3
	sector4
)
const loadFactor = 64

var rootRectangle = physics.Rectangle{
	Width:  WorldCoordMax * 2,
	Height: WorldCoordMax * 2,
	Left:   WorldCoordMin,
	Top:    WorldCoordMin,
}

type QuadTree struct {
	root *QTreeNode
}

type QTreeNode struct {
	items      map[*component.PolygonCollider]struct{}
	children   [4]*QTreeNode // points to 4 sub dimensions
	parent     *QTreeNode
	hasChild   bool // marks whether should insert element into this node, or in its children nodes.
	loadFactor int  // how many items can be held at most in this node
}

func NewQuadTree() *QuadTree {
	return &QuadTree{}
}

func NewQTreeNode(parent *QTreeNode) *QTreeNode {
	return &QTreeNode{
		items:      map[*component.PolygonCollider]struct{}{},
		children:   [4]*QTreeNode{},
		parent:     parent,
		hasChild:   false,
		loadFactor: loadFactor,
	}
}

func (qt *QuadTree) Insert(collider *component.PolygonCollider) {
	if qt.root == nil {
		qt.root = NewQTreeNode(nil)
	}
	qt.root.doInsert(collider, rootRectangle)
}

func (qt *QuadTree) Query(position linalg.Vector2f64) map[*component.PolygonCollider]struct{} {
	return qt.root.doQuery(position, rootRectangle)
}

func (qt *QTreeNode) doQuery(position linalg.Vector2f64, boundary physics.Rectangle) map[*component.PolygonCollider]struct{} {
	if !qt.hasChild {
		return qt.items
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
			return qt.children[sector1].doQuery(position, physics.NewRectangle(boundaryWidthHalf, boundaryHeightHalf, center.X, boundary.Top))
		}
		return qt.children[sector4].doQuery(position, physics.NewRectangle(boundaryWidthHalf, boundaryHeightHalf, center.X, center.Y))
	}
	if yPos {
		return qt.children[sector2].doQuery(position, physics.NewRectangle(boundaryWidthHalf, boundaryHeightHalf, boundary.Left, boundary.Top))
	}
	return qt.children[sector3].doQuery(position, physics.NewRectangle(boundaryWidthHalf, boundaryHeightHalf, boundary.Left, center.Y))
}

func (qt *QTreeNode) doInsert(collider *component.PolygonCollider, boundary physics.Rectangle) {
	if qt == nil {
		return
	}
	if qt.hasChild {
		qt.items[collider] = struct{}{}
		if len(qt.items) > qt.loadFactor {
			// split
			for item, _ := range qt.items {
				delete(qt.items, item)
				qt.insertIntoChild0(collider, boundary)
			}
			qt.hasChild = true
		}
		return
	}
	qt.insertIntoChild0(collider, boundary)
}

func (qt *QTreeNode) insertIntoChild0(collider *component.PolygonCollider, boundary physics.Rectangle) {
	anchor := collider.Collider.GetAnchor()
	boundaryWidthHalf := boundary.Width / 2
	boundaryHeightHalf := boundary.Height / 2
	center := linalg.Vector2f64{
		X: boundary.Left + boundaryWidthHalf,
		Y: boundary.Width + boundaryHeightHalf,
	}
	xPos := anchor.X > center.X
	yPos := anchor.Y > center.Y
	if xPos {
		if yPos {
			qt.insertIntoChild1(collider, sector1, boundaryWidthHalf, boundaryHeightHalf, center.X, boundary.Top)
			return
		}
		qt.insertIntoChild1(collider, sector4, boundaryWidthHalf, boundaryHeightHalf, center.X, center.Y)
		return
	}
	if yPos {
		qt.insertIntoChild1(collider, sector2, boundaryWidthHalf, boundaryHeightHalf, boundary.Left, boundary.Top)
		return
	}
	qt.insertIntoChild1(collider, sector3, boundaryWidthHalf, boundaryHeightHalf, boundary.Left, center.Y)
}

func (qt *QTreeNode) insertIntoChild1(collider *component.PolygonCollider, sector int, w float64, h float64, left float64, top float64) {
	if qt.children[sector] == nil {
		qt.children[sector] = NewQTreeNode(qt) // decrease loadFactor
	}
	qt.children[sector].doInsert(collider, physics.NewRectangle(w, h, left, top))
}
