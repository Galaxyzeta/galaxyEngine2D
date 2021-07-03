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
	ret.boundingBox = ret.CalcAndGetBoundingBox()
	return ret
}

func (poly Polygon) GetAnchor() *linalg.Vector2f64 {
	return poly.anchor
}

// GetWorldVertices converts a polygon to world coordinates system.
func (poly Polygon) GetWorldVertices() []linalg.Vector2f64 {
	var vertices []linalg.Vector2f64 = poly.vertices
	for idx, vertice := range vertices {
		vertices[idx].X = (vertice.X-poly.pivot.X)*math.Cos(poly.rotationDeg) - (vertice.Y-poly.pivot.Y)*math.Sin(poly.rotationDeg) + poly.anchor.X
		vertices[idx].Y = (vertice.X-poly.pivot.X)*math.Sin(poly.rotationDeg) - (vertice.Y-poly.pivot.Y)*math.Cos(poly.rotationDeg) + poly.anchor.Y
	}
	return vertices
}

// GetBoundingBox gets the bounding box in absolute world coordinates.
func (poly Polygon) GetBoundingBox() BoundingBox {
	ret := poly.boundingBox
	for _, elem := range ret {
		elem.X += poly.anchor.X
		elem.Y += poly.anchor.Y
	}
	return ret
}

// Intersect checks whether two polygons overlaps with eachother.
func (poly Polygon) Intersect(poly2 Polygon) bool {
	vertices := poly.GetWorldVertices()
	for i := 1; i < len(vertices); i++ {
		edgeVec := linalg.Vector2f64{X: vertices[i].X - vertices[i-1].X, Y: vertices[i].Y - vertices[i-1].Y}
		axis := edgeVec.NormalVec()
		seg0 := poly2.ProjectOn(axis)
		seg1 := poly.ProjectOn(axis)
		if !overlap(seg0, seg1) {
			return false
		}
	}
	for i := 1; i < len(vertices); i++ {
		edgeVec := linalg.Vector2f64{X: vertices[i].X - vertices[i-1].X, Y: vertices[i].Y - vertices[i-1].Y}
		axis := edgeVec.NormalVec()
		seg0 := poly2.ProjectOn(axis)
		seg1 := poly.ProjectOn(axis)
		if !overlap(seg0, seg1) {
			return false
		}
	}
	return true
}

// CalcAndGetBoundingBox returns a bounding box calculated by current
func (poly Polygon) CalcAndGetBoundingBox() BoundingBox {
	var minX, minY, maxX, maxY float64
	minX = poly.vertices[0].X
	maxX = minY
	minY = poly.vertices[0].Y
	maxY = minY
	for i := 1; i < len(poly.vertices); i++ {
		if poly.vertices[i].X < minX {
			minX = poly.vertices[i].X
		} else if poly.vertices[i].X > maxX {
			maxX = poly.vertices[i].X
		}

		if poly.vertices[i].Y < minY {
			minY = poly.vertices[i].Y
		} else if poly.vertices[i].Y > maxY {
			maxY = poly.vertices[i].Y
		}
	}
	return [4]linalg.Vector2f64{
		{X: minX, Y: minY},
		{X: maxX, Y: minY},
		{X: maxX, Y: maxY},
		{X: minX, Y: maxY},
	}
}

// ProjectOn an axis.
func (poly Polygon) ProjectOn(axis linalg.Vector2f64) linalg.Vector2f64 {
	vertices := poly.vertices
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
