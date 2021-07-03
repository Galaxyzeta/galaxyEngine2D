package component

import (
	"galaxyzeta.io/engine/physics"
)

const NamePolygonCollider = "PolygonCollider"

type PolygonCollider struct {
	Collider physics.Polygon
	Name     string
}

func (pc *PolygonCollider) GetName() string {
	return pc.Name
}
