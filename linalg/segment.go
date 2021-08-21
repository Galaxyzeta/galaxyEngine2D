package linalg

type Segmentf64 struct {
	Point1 Vector2f64
	Point2 Vector2f64
}

func NewSegmentf64(x1 float64, y1 float64, x2 float64, y2 float64) Segmentf64 {
	return Segmentf64{
		Point1: NewVector2f64(x1, y1),
		Point2: NewVector2f64(x2, y2),
	}
}

// Standardize the segment, ensures point1 always has smaller X value and Y value than point2.
func (s1 Segmentf64) Standardize() Segmentf64 {
	if s1.Point1.X > s1.Point2.X {
		s1.Point1, s1.Point2 = s1.Point2, s1.Point1
	} else if s1.Point1.X == s1.Point2.X {
		if s1.Point1.Y > s1.Point2.Y {
			s1.Point1, s1.Point2 = s1.Point2, s1.Point1
		}
	}
	return s1
}

// Cofficient returns k and b of a line.
func (s Segmentf64) Cofficient() (float64, float64) {
	k := (s.Point1.X - s.Point2.X) / (s.Point1.Y - s.Point2.Y)
	return k, s.Point1.Y - k*s.Point1.X
}

func lineGetY(k float64, b float64, x float64) float64 {
	return k*x + b
}

func sort2(x float64, y float64) (float64, float64) {
	if x < y {
		return x, y
	}
	return y, x
}

// inRange tells whether the target number is in range [left, right)
func inRange(left float64, right float64, target float64) bool {
	return target >= left && target < right
}

// ToVector returns a vector that points from point1 to point2.
func (s Segmentf64) ToVector() Vector2f64 {
	return Vector2f64{s.Point2.X - s.Point1.X, s.Point2.Y - s.Point1.Y}
}

// ToVector returns a vector that points from point2 to point1.
func (s Segmentf64) ToVectorRev() Vector2f64 {
	return Vector2f64{s.Point1.X - s.Point2.X, s.Point1.Y - s.Point2.Y}
}

// Straddle tells whether two points on one segment are on different side of another segment.
func (s1 Segmentf64) Straddle(s2 Segmentf64) bool {
	v1 := NewVector2f64(s2.Point1.X-s1.Point1.X, s2.Point1.Y-s1.Point1.Y)
	v2 := NewVector2f64(s2.Point1.X-s1.Point2.X, s2.Point1.Y-s1.Point2.Y)
	v3 := s2.ToVector()
	return v1.Mult(v3)*v2.Mult(v3) <= 0
}

// Intersect returns true when s1 intersects with another segment.
func (s1 Segmentf64) Intersect(s2 Segmentf64) bool {
	s1 = s1.Standardize()
	s2 = s2.Standardize()
	v1 := s1.ToVector()
	v2 := s2.ToVector()
	if v2.Mult(v1) == 0 {
		// parallel
		refVec1 := NewVector2f64(s2.Point1.X-s1.Point2.X, s2.Point1.Y-s1.Point2.Y)
		refVec2 := NewVector2f64(s1.Point1.X-s2.Point1.X, s1.Point1.Y-s2.Point1.Y)
		if refVec1.Mult(refVec2) == 0 {
			// in one line
			return s2.Point1.X <= s1.Point2.X && s1.Point1.X <= s2.Point1.X && s2.Point1.Y <= s1.Point2.Y && s1.Point1.Y <= s2.Point1.Y
		} else {
			// parallel
			return false
		}
	} else {
		// no parallel
		return s1.Straddle(s2) && s2.Straddle(s1)
	}
}

func (s1 Segmentf64) IsVertical() bool {
	return s1.Point1.Y == s1.Point1.Y
}
