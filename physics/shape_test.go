package physics

import (
	"testing"

	"galaxyzeta.io/engine/linalg"
)

func TestSATCollision(t *testing.T) {
	var vertices1 []linalg.Vector2f64
	var vertices2 []linalg.Vector2f64
	var vertices3 []linalg.Vector2f64
	var vertices4 []linalg.Vector2f64

	vertices1 = append(vertices1,
		linalg.Vector2f64{X: 1, Y: 1},
		linalg.Vector2f64{X: 2, Y: 1},
		linalg.Vector2f64{X: 2, Y: 2},
		linalg.Vector2f64{X: 1, Y: 2})
	poly1 := Polygon{
		Vertices: vertices1,
	}
	vertices2 = append(vertices2,
		linalg.Vector2f64{X: 3, Y: 3},
		linalg.Vector2f64{X: 4, Y: 3},
		linalg.Vector2f64{X: 4, Y: 4},
		linalg.Vector2f64{X: 3, Y: 4})
	poly2 := Polygon{
		Vertices: vertices2,
	}
	t.Log(poly1.Intersect(poly2)) // false
	t.Log(poly2.Intersect(poly1)) // false
	vertices3 = append(vertices3,
		linalg.Vector2f64{X: 1.5, Y: 1.5},
		linalg.Vector2f64{X: 3.5, Y: 1.5},
		linalg.Vector2f64{X: 3.5, Y: 3.5},
		linalg.Vector2f64{X: 1.5, Y: 3.5})
	poly3 := Polygon{
		Vertices: vertices3,
	}
	t.Log(poly3.Intersect(poly1))
	t.Log(poly3.Intersect(poly2))
	t.Log(poly3.Intersect(poly3))
	vertices4 = append(vertices4,
		linalg.Vector2f64{X: 3, Y: 2},
		linalg.Vector2f64{X: 3.5, Y: 3.5},
		linalg.Vector2f64{X: 1.5, Y: 3.5})
	poly4 := Polygon{
		Vertices: vertices4,
	}
	t.Log(poly4.Intersect(poly1)) // false
	t.Log(poly4.Intersect(poly2))
	t.Log(poly4.Intersect(poly3))

}

func TestCircle2Polygon(t *testing.T) {
	circle := Circle{
		Left:      0,
		Top:       0,
		Radius:    1,
		Percision: 6,
	}
	poly := circle.ToPolygon()
	t.Log(poly)
}
