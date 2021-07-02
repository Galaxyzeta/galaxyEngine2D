package physics

import (
	"math"

	"galaxyzeta.io/engine/linalg"
)

type IShape interface {
	Intersect(shape IShape) bool
}

type Rotation struct {
	Pivot       linalg.Vector2f64
	RotationDeg float64
}

type Polygon struct {
	Vertices []linalg.Vector2f64 // original vertices with no rotation
	Rotation
}

type Rectangle struct {
	Width  float64
	Height float64
	Left   float64
	Top    float64
}

type Circle struct {
	Left      float64
	Top       float64
	Radius    float64
	Percision int
}

func (rect Rectangle) Intersect(shape IShape) bool {
	switch shape := shape.(type) {
	case *Rectangle:
		return rect.IntersectWithRectangle(shape)
	}
	return false
}

func (rect Rectangle) IntersectWithRectangle(anotherRect *Rectangle) bool {
	anotherRight := anotherRect.Left + anotherRect.Width
	anotherBottom := anotherRect.Top + anotherRect.Height
	thisRight := rect.Left + rect.Width
	thisBottom := rect.Top + rect.Height
	return thisRight >= anotherRect.Left && rect.Left <= anotherRight && rect.Top >= anotherBottom && thisBottom <= anotherRect.Top
}

func (circle Circle) ToPolygon() Polygon {
	var dirDelta float64 = float64(360.0 / circle.Percision)
	var vertices []linalg.Vector2f64
	var deg float64 = 0
	for deg < 360.0 {
		rad := linalg.Deg2Rad(deg)
		vertices = append(vertices, linalg.Vector2f64{X: circle.Radius * math.Cos(rad), Y: circle.Radius * math.Sin(rad)})
		deg += dirDelta
	}
	return Polygon{
		Vertices: vertices,
	}
}

// ToPolygon converts a rectangle into polygon.
func (rect Rectangle) ToPolygon() Polygon {
	var vertices []linalg.Vector2f64
	vertices = append(vertices,
		linalg.Vector2f64{X: rect.Left, Y: rect.Top},
		linalg.Vector2f64{X: rect.Left + rect.Width, Y: rect.Top},
		linalg.Vector2f64{X: rect.Left + rect.Width, Y: rect.Top + rect.Height},
		linalg.Vector2f64{X: rect.Left, Y: rect.Top + rect.Height})
	return Polygon{
		Vertices: vertices,
	}
}

func (poly Polygon) GetRotatedVertices() []linalg.Vector2f64 {
	var vertices []linalg.Vector2f64 = poly.Vertices
	for idx, vertice := range vertices {
		vertices[idx].X = (vertice.X-poly.Pivot.X)*math.Cos(poly.RotationDeg) - (vertice.Y-poly.Pivot.Y)*math.Sin(poly.RotationDeg)
		vertices[idx].Y = (vertice.X-poly.Pivot.X)*math.Sin(poly.RotationDeg) - (vertice.Y-poly.Pivot.Y)*math.Cos(poly.RotationDeg)
	}
	return vertices
}

// Intersect checks whether two polygons overlaps with eachother.
func (poly Polygon) Intersect(poly2 Polygon) bool {
	vertices := poly.GetRotatedVertices()
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

func (poly Polygon) ProjectOn(axis linalg.Vector2f64) linalg.Vector2f64 {
	vertices := poly.Vertices
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
