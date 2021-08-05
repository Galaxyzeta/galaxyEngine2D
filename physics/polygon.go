package physics

import (
	"math"

	"galaxyzeta.io/engine/linalg"
)

type Polygon struct {
	anchor      *linalg.Vector2f64  // anchor is a base point where all world position was calculated on
	vertices    []linalg.Vector2f64 // original vertices with no rotation
	boundingBox BoundingBox
	pivot       linalg.Vector2f64
	rotationDeg float64
}

func NewPolygon(anchor *linalg.Vector2f64, pivot linalg.Vector2f64, deg float64, vertices []linalg.Vector2f64) *Polygon {
	ret := &Polygon{
		anchor:      anchor,
		vertices:    vertices,
		pivot:       pivot,
		rotationDeg: deg,
	}
	ret.boundingBox = ret.GetBoundingBox()
	return ret
}

func NewStaticPolygon(pivot linalg.Vector2f64, deg float64, vertices []linalg.Vector2f64) *Polygon {
	ret := &Polygon{
		vertices:    vertices,
		pivot:       pivot,
		rotationDeg: deg,
	}
	ret.boundingBox = ret.GetBoundingBox()
	return ret
}

func (poly Polygon) GetAnchor() *linalg.Vector2f64 {
	return poly.anchor
}

// GetWorldVertices converts a polygon to world coordinates system.
func (poly Polygon) GetWorldVertices() []linalg.Vector2f64 {
	verticesReplica := make([]linalg.Vector2f64, len(poly.vertices))
	copy(verticesReplica[:], poly.vertices)
	rotRad := linalg.Deg2Rad(poly.rotationDeg)
	var anchorX float64 = 0
	var anchorY float64 = 0
	if poly.anchor != nil {
		anchorX = poly.anchor.X
		anchorY = poly.anchor.Y
	}
	for idx, vertice := range verticesReplica {
		x := (vertice.X-poly.pivot.X)*math.Cos(rotRad) - (vertice.Y-poly.pivot.Y)*math.Sin(rotRad) + anchorX
		y := (vertice.X-poly.pivot.X)*math.Sin(rotRad) + (vertice.Y-poly.pivot.Y)*math.Cos(rotRad) + anchorY
		verticesReplica[idx].X = x
		verticesReplica[idx].Y = y
	}
	return verticesReplica
}

// Intersect checks whether two polygons overlaps with eachother.
func (poly Polygon) Intersect(poly2 Polygon) bool {
	vertices := poly.GetWorldVertices()
	vertices2 := poly2.GetWorldVertices()
	for i := 1; i < len(vertices); i++ {
		edgeVec := linalg.Vector2f64{X: vertices[i].X - vertices[i-1].X, Y: vertices[i].Y - vertices[i-1].Y}
		axis := edgeVec.NormalVec()
		seg0 := poly2.ProjectOn(axis)
		seg1 := poly.ProjectOn(axis)
		if !overlap(seg0, seg1) {
			return false
		}
	}
	for i := 1; i < len(vertices2); i++ {
		edgeVec := linalg.Vector2f64{X: vertices2[i].X - vertices2[i-1].X, Y: vertices2[i].Y - vertices2[i-1].Y}
		axis := edgeVec.NormalVec()
		seg0 := poly2.ProjectOn(axis)
		seg1 := poly.ProjectOn(axis)
		if !overlap(seg0, seg1) {
			return false
		}
	}
	return true
}

// GetBoundingBox returns a bounding box calculated by current conditions.
func (poly Polygon) GetBoundingBox() BoundingBox {
	currentWorldVetices := poly.GetWorldVertices()
	var minX, minY, maxX, maxY float64
	minX = currentWorldVetices[0].X
	maxX = minY
	minY = currentWorldVetices[0].Y
	maxY = minY
	for i := 1; i < len(currentWorldVetices); i++ {
		if currentWorldVetices[i].X < minX {
			minX = currentWorldVetices[i].X
		} else if currentWorldVetices[i].X > maxX {
			maxX = currentWorldVetices[i].X
		}

		if currentWorldVetices[i].Y < minY {
			minY = currentWorldVetices[i].Y
		} else if currentWorldVetices[i].Y > maxY {
			maxY = currentWorldVetices[i].Y
		}
	}
	return [4]linalg.Vector2f64{
		{X: maxX, Y: minY}, // top right
		{X: minX, Y: minY}, // top left
		{X: minX, Y: maxY}, // bot left
		{X: maxX, Y: maxY}, // bot right
	}
}

// ProjectOn an axis.
func (poly Polygon) ProjectOn(axis linalg.Vector2f64) linalg.Vector2f64 {
	vertices := poly.GetWorldVertices()
	min := vertices[0].Dot(axis)
	max := min
	for i := 1; i < len(vertices); i++ {
		dotProduct := vertices[i].Dot(axis)
		if dotProduct < min {
			min = dotProduct
		} else if dotProduct > max {
			max = dotProduct
		}
	}
	return linalg.Vector2f64{X: min, Y: max}
}

// overlap judges whether two segments on a same axis overlaps.
func overlap(a linalg.Vector2f64, b linalg.Vector2f64) bool {
	leftMost, rightMost := a, b
	if a.X > b.X {
		leftMost, rightMost = b, a
	}
	if rightMost.X >= leftMost.Y {
		return false
	}
	return true
}
