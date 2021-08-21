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

func ColliderAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		return col.I().Obj().Name == name
	})
}

func CollidersAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		return col.I().Obj().Name == name
	})
}

func ColliderAtWithTag(sys ICollisionSystem, pos linalg.Vector2f64, name string) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		_, ok := col.I().Obj().Tags[name]
		return ok
	})
}

func CollidersAtWithTag(sys ICollisionSystem, pos linalg.Vector2f64, tag string) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		_, ok := col.I().Obj().Tags[tag]
		return ok
	})
}

func ObjectAt(sys ICollisionSystem, pos linalg.Vector2f64) base.IGameObject2D {
	if val := ColliderAt(sys, pos); val != nil {
		return val.I()
	}
	return nil
}

func ObjectsAt(sys ICollisionSystem, pos linalg.Vector2f64) []base.IGameObject2D {
	cols := CollidersAt(sys, pos)
	ret := make([]base.IGameObject2D, 0)
	for _, col := range cols {
		ret = append(ret, col.I())
	}
	return ret
}

func ObjectAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string) base.IGameObject2D {
	for _, iobj := range ObjectsAt(sys, pos) {
		if name == iobj.Obj().Name {
			return iobj
		}
	}
	return nil
}

func HasAnyObjectAt(sys ICollisionSystem, pos linalg.Vector2f64) bool {
	return ObjectAt(sys, pos) != nil
}

func HasObjectAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string) bool {
	for _, iobj := range ObjectsAt(sys, pos) {
		if name == iobj.Obj().Name {
			return true
		}
	}
	return false
}

func HasObjectAtWithTag(sys ICollisionSystem, pos linalg.Vector2f64, tag string) bool {
	for _, iobj := range ObjectsAt(sys, pos) {
		if _, ok := iobj.Obj().Tags[tag]; ok {
			return true
		}
	}
	return false
}

func HasObjectAtWithTags(sys ICollisionSystem, pos linalg.Vector2f64, tags []string) bool {
	for _, iobj := range ObjectsAt(sys, pos) {
		match := true
		for _, tag := range tags {
			if _, ok := iobj.Obj().Tags[tag]; !ok {
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

// === polygon ===

func checkPolygonCollision(pc component.PolygonCollider, pool []*component.PolygonCollider) bool {
	return checkPolygonCollisionWithFilter(pc, pool, func(test *component.PolygonCollider) bool { return true })
}

func checkPolygonCollisionWithFilter(pc component.PolygonCollider, pool []*component.PolygonCollider, fx func(test *component.PolygonCollider) bool) bool {
	for _, test := range pool {
		if pc.Collider.Intersect(test.Collider) && fx(test) {
			return true
		}
	}
	return false
}

func collectPolygonCollisionWithFilter(pc component.PolygonCollider, pool []*component.PolygonCollider, fx func(testpc *component.PolygonCollider) bool) *component.PolygonCollider {
	for _, testpc := range pool {
		if pc.Collider.Intersect(testpc.Collider) && fx(testpc) {
			return testpc
		}
	}
	return nil
}

func collectPolygonCollisionsWithFilter(pc component.PolygonCollider, pool []*component.PolygonCollider, fx func(testpc *component.PolygonCollider) bool) (result []*component.PolygonCollider) {
	for _, testpc := range pool {
		if pc.Collider.Intersect(testpc.Collider) && fx(testpc) {
			result = append(result, testpc)
		}
	}
	return result
}

func getPcwrapper(p physics.Polygon) component.PolygonCollider {
	return component.PolygonCollider{
		Collider: p,
	}
}

func HasColliderAtPolygonWithAny(sys ICollisionSystem, p physics.Polygon) bool {
	pcWrapper := getPcwrapper(p)
	return checkPolygonCollision(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper))
}

func HasColliderAtPolygonWithName(sys ICollisionSystem, p physics.Polygon, name string) bool {
	pcWrapper := getPcwrapper(p)
	return checkPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper), func(test *component.PolygonCollider) bool {
		return test.I().Obj().Name == name
	})
}

func HasColliderAtPolygonWithTag(sys ICollisionSystem, p physics.Polygon, tag string) bool {
	pcWrapper := getPcwrapper(p)
	return checkPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper), func(test *component.PolygonCollider) bool {
		_, ok := test.I().Obj().Tags[tag]
		return ok
	})
}

func ColliderAtPolygonWithAny(sys ICollisionSystem, p physics.Polygon) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper), func(testpc *component.PolygonCollider) bool {
		return true
	})
}

func ColliderAtPolygonWithFilter(sys ICollisionSystem, p physics.Polygon, fx func(test *component.PolygonCollider) bool) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper), fx)
}

func ColliderAtPolygonWithName(sys ICollisionSystem, p physics.Polygon, name string) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper), func(testpc *component.PolygonCollider) bool {
		return testpc.I().Obj().Name == name
	})
}

func ColliderAtPolygonWithTag(sys ICollisionSystem, p physics.Polygon, tag string) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper), func(testpc *component.PolygonCollider) bool {
		_, ok := testpc.I().Obj().Tags[tag]
		return ok
	})
}
