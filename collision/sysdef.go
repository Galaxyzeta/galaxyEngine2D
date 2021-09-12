package collision

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
)

type ICollisionSystem interface {
	base.ISystem
	QueryNeighborCollidersWithCollider(col component.PolygonCollider, mode QueryMode) []*component.PolygonCollider
	QueryNeighborCollidersWithPosition(pos linalg.Vector2f64, mode QueryMode) []*component.PolygonCollider
	QueryNeighborCollidersWithColliderAndFilter(col component.PolygonCollider, filter func(*component.PolygonCollider) bool, mode QueryMode) []*component.PolygonCollider
	QueryNeighborCollidersWithPositionAndFilter(pos linalg.Vector2f64, filter func(*component.PolygonCollider) bool, mode QueryMode) []*component.PolygonCollider
}
