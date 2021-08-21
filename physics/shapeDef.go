package physics

import (
	"galaxyzeta.io/engine/linalg"
)

type IShape interface {
	Intersect(shape IShape) bool
}

type Rotation struct {
	Pivot       linalg.Vector2f64
	RotationDeg float64
}

type Point struct {
	X float64
	Y float64
}
