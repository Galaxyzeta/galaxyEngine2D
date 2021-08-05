package collision

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

const pointEpsHalf = 0.01
const pointEps = 0.02

func colliderAtWithFilter(sys ICollisionSystem, pos linalg.Vector2f64, fx func(col *component.PolygonCollider) bool) *component.PolygonCollider {
	cols := sys.QueryNeighborCollidersWithPosition(pos)
	for _, col := range cols {
		// use a tiny rectangle to test intersect
		if physics.NewRectangle(pos.X-pointEpsHalf, pos.Y-pointEpsHalf, pointEps, pointEps).ToPolygon().Intersect(col.Collider) && fx(col) {
			return col
		}
	}
	return nil
}

func collidersAtWithFilter(sys ICollisionSystem, pos linalg.Vector2f64, fx func(col *component.PolygonCollider) bool) []*component.PolygonCollider {
	cols := sys.QueryNeighborCollidersWithPosition(pos)
	ret := make([]*component.PolygonCollider, 0)
	for _, col := range cols {
		// use a tiny rectangle to test intersect
		if physics.NewRectangle(pos.X-pointEpsHalf, pos.Y-pointEpsHalf, pointEps, pointEps).ToPolygon().Intersect(col.Collider) && fx(col) {
			ret = append(ret, col)
		}
	}
	return ret
}

func ColliderAt(sys ICollisionSystem, pos linalg.Vector2f64) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, nil)
}

func CollidersAt(sys ICollisionSystem, pos linalg.Vector2f64) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, nil)
}

func ColliderAtWithName(sys ICollisionSystem, name string, pos linalg.Vector2f64) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		return col.GetIGameObject2D().GetGameObject2D().Name == name
	})
}

func CollidersAtWithName(sys ICollisionSystem, name string, pos linalg.Vector2f64) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		return col.GetIGameObject2D().GetGameObject2D().Name == name
	})
}

func ColliderAtWithTag(sys ICollisionSystem, name string, pos linalg.Vector2f64) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		_, ok := col.GetIGameObject2D().GetGameObject2D().Tags[name]
		return ok
	})
}

func ObjectAt(sys ICollisionSystem, pos linalg.Vector2f64) base.IGameObject2D {
	if val := ColliderAt(sys, pos); val != nil {
		return val.GetIGameObject2D()
	}
	return nil
}

func ObjectsAt(sys ICollisionSystem, pos linalg.Vector2f64) []base.IGameObject2D {
	cols := CollidersAt(sys, pos)
	ret := make([]base.IGameObject2D, 0)
	for _, col := range cols {
		ret = append(ret, col.GetIGameObject2D())
	}
	return ret
}

func NamedObjectAt(sys ICollisionSystem, name string, pos linalg.Vector2f64) base.IGameObject2D {
	for _, iobj := range ObjectsAt(sys, pos) {
		if name == iobj.GetGameObject2D().Name {
			return iobj
		}
	}
	return nil
}

func HasAnyObjectAt(sys ICollisionSystem, pos linalg.Vector2f64) bool {
	return ObjectAt(sys, pos) != nil
}

func HasNamedObjectAt(sys ICollisionSystem, name string, pos linalg.Vector2f64) bool {
	for _, iobj := range ObjectsAt(sys, pos) {
		if name == iobj.GetGameObject2D().Name {
			return true
		}
	}
	return false
}

func HasObjectAtWithTag(sys ICollisionSystem, tag string, pos linalg.Vector2f64) bool {
	for _, iobj := range ObjectsAt(sys, pos) {
		if _, ok := iobj.GetGameObject2D().Tags[tag]; ok {
			return true
		}
	}
	return false
}

func HasObjectAtWithTags(sys ICollisionSystem, tags []string, pos linalg.Vector2f64) bool {
	for _, iobj := range ObjectsAt(sys, pos) {
		match := true
		for _, tag := range tags {
			if _, ok := iobj.GetGameObject2D().Tags[tag]; !ok {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
