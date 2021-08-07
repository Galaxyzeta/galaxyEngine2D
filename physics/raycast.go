package physics

import (
	"math"

	"galaxyzeta.io/engine/linalg"
)

type Ray struct {
	Vec    linalg.Vector2f64
	Origin linalg.Vector2f64
}

type RaycastHit struct {
	Hits []linalg.Point2f64
}

// Coefficient calculates the k and b value.
func (r Ray) Coefficient() (float64, float64) {
	k := r.Vec.Y / r.Vec.X
	return k, r.Origin.Y - k*r.Origin.X
}

// IsVertical reports whether the ray points straight along the y axis.
func (r Ray) IsVertical() bool {
	return r.Vec.X == 0
}

func (r Ray) IntersectPolygon(p Polygon) bool {
	points := p.GetWorldVertices()
	refVec := linalg.NewVector2f64(points[0].X-r.Origin.X, points[0].Y-r.Origin.Y)
	sgn := math.Signbit(r.Vec.Mult(refVec))
	for i := 1; i < len(points); i++ {
		refVec = linalg.NewVector2f64(points[i].X-r.Origin.X, points[i].Y-r.Origin.Y)
		if sgn != math.Signbit(r.Vec.Mult(refVec)) {
			return true
		}
	}
	return false
}

func (r Ray) IntersectSegment(s linalg.Segmentf64) bool {
	s = s.Standardize()
	refVec1 := linalg.NewVector2f64(s.Point1.X-r.Origin.X, s.Point1.Y-r.Origin.Y)
	if refVec1.Dot(r.Vec) < 0 {
		return false
	}
	refVec2 := linalg.NewVector2f64(s.Point2.X-r.Origin.X, s.Point2.Y-r.Origin.Y)
	return refVec1.Mult(refVec2)*refVec2.Mult(refVec1) <= 0
}

func (r Ray) IntersectSegmentDetail(s linalg.Segmentf64) (bool, linalg.Vector2f64) {
	s = s.Standardize()
	refVec1 := linalg.NewVector2f64(s.Point1.X-r.Origin.X, s.Point1.Y-r.Origin.Y)
	if refVec1.Dot(r.Vec) < 0 {
		return false, linalg.Vector2f64{}
	}
	refVec2 := linalg.NewVector2f64(s.Point2.X-r.Origin.X, s.Point2.Y-r.Origin.Y)
	if refVec1.Mult(refVec2)*refVec2.Mult(refVec1) > 0 {
		return false, linalg.Vector2f64{}
	}

	// calculate intersection
	if r.IsVertical() {
		k1, b1 := s.Cofficient()
		return true, linalg.NewVector2f64(r.Origin.X, k1*r.Origin.X+b1)
	} else if s.IsVertical() {
		k2, b2 := r.Coefficient()
		return true, linalg.NewVector2f64(s.Point1.X, k2*r.Origin.X+b2)
	}
	k1, b1 := s.Cofficient()
	k2, b2 := r.Coefficient()
	intx := (b2 - b1) / (k1 - k2)
	return true, linalg.NewVector2f64(intx, k1*intx+b1)
}
