package component

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/physics"
)

const NamePolygonCollider = "PolygonCollider"

type PolygonCollider struct {
	Collider physics.Polygon
	Name     string
	iobj2d   base.IGameObject2D // attached gameobject2D
	Sr       *SpriteRenderer    // if spriteRenderer is not nil, collider will always synchronize with Sr's hit box.
}

func NewPolygonCollider(collider physics.Polygon, iobj2d base.IGameObject2D) *PolygonCollider {
	return &PolygonCollider{
		Collider: collider,
		Name:     NamePolygonCollider,
		iobj2d:   iobj2d,
	}
}

func NewPolygonColliderDynamicHitbox(followSr *SpriteRenderer, iobj2d base.IGameObject2D) *PolygonCollider {
	return &PolygonCollider{
		Collider: followSr.GetHitbox(),
		Name:     NamePolygonCollider,
		iobj2d:   iobj2d,
		Sr:       followSr,
	}
}

func (pc *PolygonCollider) GetName() string {
	return pc.Name
}

// I returns IGameObject2D, the representation and abstraction of a gameObject.
func (pc *PolygonCollider) I() base.IGameObject2D {
	return pc.iobj2d
}
