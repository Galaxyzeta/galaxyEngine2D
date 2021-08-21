package collision

import (
	"testing"

	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

func TestInsertIntoQuadTree(t *testing.T) {
	qt := NewQuadTree(physics.NewRectangle(-128, -128, 256, 256), 2, 64)
	pivot := linalg.NewVector2f64(0, 0)
	vertices := [4]linalg.Vector2f64{
		{X: -1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: -1}, {X: -1, Y: -1},
	}
	cnt := 0
	for i := -128.0; i < 128; i += 70 {
		for j := -128.0; j < 128; j += 70 {
			t.Log(i, j)
			cnt++
			qt.Insert(&component.PolygonCollider{
				Collider: *physics.NewPolygon(&linalg.Vector2f64{X: i, Y: j},
					pivot, 0, vertices[:]),
			})
		}
	}
	qt.Traverse(func(pc *component.PolygonCollider, node *QTreeNode, at AreaType, idx int) bool {
		cnt--
		t.Log(pc.Collider.GetBoundingBox(), node)
		return false
	})
	t.Log(cnt)

	t.Log("Q1==============")
	for _, elem := range qt.QueryByPoint(linalg.NewVector2f64(1, 1)) {
		t.Log(elem.Collider.GetWorldVertices())
	}
	t.Log("Q2==============")
	for _, elem := range qt.QueryByPoint(linalg.NewVector2f64(-1, 1)) {
		t.Log(elem.Collider.GetWorldVertices())
	}
	t.Log("Q3==============")
	for _, elem := range qt.QueryByPoint(linalg.NewVector2f64(-1, -1)) {
		t.Log(elem.Collider.GetWorldVertices())
	}
	t.Log("Q4==============")
	for _, elem := range qt.QueryByPoint(linalg.NewVector2f64(1, -1)) {
		t.Log(elem.Collider.GetWorldVertices())
	}
}
