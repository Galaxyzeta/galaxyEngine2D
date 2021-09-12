package objs

import (
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/ecs/component"
)

type BasicComponentsBundle struct {
	tf   *component.Transform2D     `gxen:"tf"`
	rb   *component.RigidBody2D     `gxen:"rb"`
	pc   *component.PolygonCollider `gxen:"pc"`
	sr   *component.SpriteRenderer  `gxen:"sr"`
	csys collision.ICollisionSystem
}
