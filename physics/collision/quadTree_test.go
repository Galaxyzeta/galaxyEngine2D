package collision

import (
	"fmt"
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
	for i := -128.0; i < 128; i += 64 {
		for j := -128.0; j < 128; j += 64 {
			qt.Insert(&component.PolygonCollider{
				Collider: *physics.NewPolygon(&linalg.Vector2f64{X: i, Y: j},
					pivot, 0, vertices[:]),
			})
		}
	}
	qt.Traverse(func(pc *component.PolygonCollider, r physics.Rectangle, node *QTreeNode) {
		fmt.Println(pc.Collider.GetBoundingBox(), node)
	})
}
