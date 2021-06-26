package linalg

import "math"

func Deg2Rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// InvertDeg subtract deg from 360 to invert up and down.
func InvertDeg(deg float64) float64 {
	return 360 - deg
}
