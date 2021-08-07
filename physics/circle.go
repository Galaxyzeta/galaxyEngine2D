package physics

import (
	"math"

	"galaxyzeta.io/engine/linalg"
)

type Circle struct {
	Left      float64
	Top       float64
	Radius    float64
	Percision int
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
