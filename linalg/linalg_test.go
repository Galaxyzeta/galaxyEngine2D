package linalg

import (
	"testing"

	"galaxyzeta.io/engine/infra/require"
)

func TestVectorTheta(t *testing.T) {
	v1 := Vector2f64{X: 5, Y: 0}
	v2 := Vector2f64{X: 0, Y: 5}
	t.Log(v2.Theta(v1))
	t.Log(v2.ThetaDeg(v1))
}

func TestSegmentIntersect(t *testing.T) {
	// test parallel
	s1 := NewSegmentf64(0, 0, 1, 1)
	s2 := NewSegmentf64(3, 3, 4, 4)
	require.EqBool(false, s1.Intersect(s2))
	// test inline
	s1 = NewSegmentf64(0, 0, 2, 2)
	s2 = NewSegmentf64(1, 1, 4, 4)
	require.EqBool(true, s1.Intersect(s2))
	// test no parallel
	s1 = NewSegmentf64(0, 0, 2, 2)
	s2 = NewSegmentf64(0, 1, 2, 0)
	require.EqBool(true, s1.Intersect(s2))
}
