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

func NewRectangle(left float64, top float64, w float64, h float64) Rectangle {
	return Rectangle{
		Width:  w,
		Height: h,
		Left:   left,
		Top:    top,
	}
}

func (rect Rectangle) Intersect(shape IShape) bool {
	switch shape := shape.(type) {
	case *Rectangle:
		return rect.IntersectWithRectangle(shape)
	}
	return false
}

func (rect Rectangle) InsideRectangle(anotherRect *Rectangle) bool {
	anotherRight := anotherRect.Left + anotherRect.Width
	anotherBottom := anotherRect.Top + anotherRect.Height
	thisRight := rect.Left + rect.Width
	thisBottom := rect.Top + rect.Height
	return rect.Left >= anotherRect.Left && thisRight <= anotherRight && rect.Top <= anotherRect.Top && thisBottom >= anotherBottom
}

func (rect Rectangle) IntersectWithRectangle(anotherRect *Rectangle) bool {
	anotherRight := anotherRect.Left + anotherRect.Width
	anotherBottom := anotherRect.Top + anotherRect.Height
	thisRight := rect.Left + rect.Width
	thisBottom := rect.Top + rect.Height
	return thisRight >= anotherRect.Left && rect.Left <= anotherRight && rect.Top <= anotherBottom && thisBottom >= anotherRect.Top
}

// CropOutside extends or deminishes current rectangle by certain width and height.
func (rect Rectangle) CropOutside(w float64, h float64) Rectangle {
	rect.Height += h + h
	rect.Width += w + w
	rect.Left -= w
	rect.Top -= w
	return rect
}

// ToPolygon converts a circle into polygon.
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
		vertices: vertices,
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
		vertices: vertices,
	}
}
