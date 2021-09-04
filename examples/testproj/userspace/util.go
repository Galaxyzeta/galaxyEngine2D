package objs

import (
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/ecs/component"
)

type BasicComponentsBundle struct {
	tf   *component.Transform2D
	rb   *component.RigidBody2D
	pc   *component.PolygonCollider
	sr   *component.SpriteRenderer
	csys collision.ICollisionSystem
}
