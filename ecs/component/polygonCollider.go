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
}

func NewPolygonCollider(collider physics.Polygon, iobj2d base.IGameObject2D) *PolygonCollider {
	return &PolygonCollider{
		Collider: collider,
		Name:     NamePolygonCollider,
		iobj2d:   iobj2d,
	}
}

func (pc *PolygonCollider) GetName() string {
	return pc.Name
}

func (pc *PolygonCollider) GetIGameObject2D() base.IGameObject2D {
	return pc.iobj2d
}
