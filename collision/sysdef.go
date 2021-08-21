package collision

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
)

type ICollisionSystem interface {
	base.ISystem
	QueryNeighborCollidersWithCollider(col component.PolygonCollider) []*component.PolygonCollider
	QueryNeighborCollidersWithPosition(pos linalg.Vector2f64) []*component.PolygonCollider
	QueryNeighborCollidersWithColliderAndFilter(col component.PolygonCollider, filter func(*component.PolygonCollider) bool) []*component.PolygonCollider
	QueryNeighborCollidersWithPositionAndFilter(pos linalg.Vector2f64, filter func(*component.PolygonCollider) bool) []*component.PolygonCollider
}
