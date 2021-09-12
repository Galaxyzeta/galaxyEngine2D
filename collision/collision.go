package collision

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

const pointEpsHalf = 0.01
const pointEps = 0.02

type CollisionBuilder struct {
	csys     ICollisionSystem
	cType    CollisionType
	cRetType CollisionResultType
}

type CollisionType int8

const (
	CTypePoint CollisionType = iota
	CTypePolygon
	CTypeRay
)

type CollisionResultType int8

const (
	CRetTypeObject CollisionResultType = iota
	CRetTypeObjects
	CRetTypePC
	CRetTypePCs
	CRetTypeBool
)

func Collision(csys ICollisionSystem, cType CollisionType, cRetType CollisionResultType) CollisionBuilder {
	return CollisionBuilder{
		csys:     csys,
		cType:    cType,
		cRetType: cRetType,
	}
}

func (cb CollisionBuilder) Fetch() {

}

func colliderAtWithFilter(sys ICollisionSystem, pos linalg.Vector2f64, fx func(col *component.PolygonCollider) bool, mode QueryMode) *component.PolygonCollider {
	cols := sys.QueryNeighborCollidersWithPosition(pos, mode)
	for _, col := range cols {
		// use a tiny rectangle to test intersect
		if physics.NewRectangle(pos.X-pointEpsHalf, pos.Y-pointEpsHalf, pointEps, pointEps).ToPolygon().Intersect(col.Collider) && fx(col) {
			return col
		}
	}
	return nil
}

func collidersAtWithFilter(sys ICollisionSystem, pos linalg.Vector2f64, fx func(col *component.PolygonCollider) bool, mode QueryMode) []*component.PolygonCollider {
	cols := sys.QueryNeighborCollidersWithPosition(pos, mode)
	ret := make([]*component.PolygonCollider, 0)
	for _, col := range cols {
		// use a tiny rectangle to test intersect
		if physics.NewRectangle(pos.X-pointEpsHalf, pos.Y-pointEpsHalf, pointEps, pointEps).ToPolygon().Intersect(col.Collider) && fx(col) {
			ret = append(ret, col)
		}
	}
	return ret
}

func ColliderAt(sys ICollisionSystem, pos linalg.Vector2f64, mode QueryMode) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, nil, mode)
}

func CollidersAt(sys ICollisionSystem, pos linalg.Vector2f64, mode QueryMode) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, nil, mode)
}

func ColliderAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string, mode QueryMode) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		return col.I().Obj().Name == name
	}, mode)
}

func CollidersAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string, mode QueryMode) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		return col.I().Obj().Name == name
	}, mode)
}

func ColliderAtWithTag(sys ICollisionSystem, pos linalg.Vector2f64, name string, mode QueryMode) *component.PolygonCollider {
	return colliderAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		_, ok := col.I().Obj().Tags[name]
		return ok
	}, mode)
}

func CollidersAtWithTag(sys ICollisionSystem, pos linalg.Vector2f64, tag string, mode QueryMode) []*component.PolygonCollider {
	return collidersAtWithFilter(sys, pos, func(col *component.PolygonCollider) bool {
		_, ok := col.I().Obj().Tags[tag]
		return ok
	}, mode)
}

func ObjectAt(sys ICollisionSystem, pos linalg.Vector2f64, mode QueryMode) base.IGameObject2D {
	if val := ColliderAt(sys, pos, mode); val != nil {
		return val.I()
	}
	return nil
}

func ObjectsAt(sys ICollisionSystem, pos linalg.Vector2f64, mode QueryMode) []base.IGameObject2D {
	cols := CollidersAt(sys, pos, mode)
	ret := make([]base.IGameObject2D, 0)
	for _, col := range cols {
		ret = append(ret, col.I())
	}
	return ret
}

func ObjectAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string, mode QueryMode) base.IGameObject2D {
	for _, iobj := range ObjectsAt(sys, pos, mode) {
		if name == iobj.Obj().Name {
			return iobj
		}
	}
	return nil
}

func HasAnyObjectAt(sys ICollisionSystem, pos linalg.Vector2f64, mode QueryMode) bool {
	return ObjectAt(sys, pos, mode) != nil
}

func HasObjectAtWithName(sys ICollisionSystem, pos linalg.Vector2f64, name string, mode QueryMode) bool {
	for _, iobj := range ObjectsAt(sys, pos, mode) {
		if name == iobj.Obj().Name {
			return true
		}
	}
	return false
}

func HasObjectAtWithTag(sys ICollisionSystem, pos linalg.Vector2f64, tag string, mode QueryMode) bool {
	for _, iobj := range ObjectsAt(sys, pos, mode) {
		if _, ok := iobj.Obj().Tags[tag]; ok {
			return true
		}
	}
	return false
}

func HasObjectAtWithTags(sys ICollisionSystem, pos linalg.Vector2f64, tags []string, query QueryMode) bool {
	for _, iobj := range ObjectsAt(sys, pos, query) {
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

func HasColliderAtPolygonWithAny(sys ICollisionSystem, p physics.Polygon, mode QueryMode) bool {
	pcWrapper := getPcwrapper(p)
	return checkPolygonCollision(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode))
}

func HasColliderAtPolygonWithName(sys ICollisionSystem, p physics.Polygon, name string, mode QueryMode) bool {
	pcWrapper := getPcwrapper(p)
	return checkPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), func(test *component.PolygonCollider) bool {
		return test.I().Obj().Name == name
	})
}

func HasColliderAtPolygonWithTag(sys ICollisionSystem, p physics.Polygon, tag string, mode QueryMode) bool {
	pcWrapper := getPcwrapper(p)
	return checkPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), func(test *component.PolygonCollider) bool {
		_, ok := test.I().Obj().Tags[tag]
		return ok
	})
}

func ColliderAtPolygonWithAny(sys ICollisionSystem, p physics.Polygon, mode QueryMode) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), func(testpc *component.PolygonCollider) bool {
		return true
	})
}

func CollidersAtPolygonWithAny(sys ICollisionSystem, p physics.Polygon, mode QueryMode) []*component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionsWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), func(testpc *component.PolygonCollider) bool {
		return true
	})
}

func ColliderAtPolygonWithFilter(sys ICollisionSystem, p physics.Polygon, fx func(test *component.PolygonCollider) bool, mode QueryMode) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), fx)
}

func ColliderAtPolygonWithName(sys ICollisionSystem, p physics.Polygon, name string, mode QueryMode) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), func(testpc *component.PolygonCollider) bool {
		return testpc.I().Obj().Name == name
	})
}

func ColliderAtPolygonWithTag(sys ICollisionSystem, p physics.Polygon, tag string, mode QueryMode) *component.PolygonCollider {
	pcWrapper := getPcwrapper(p)
	return collectPolygonCollisionWithFilter(pcWrapper, sys.QueryNeighborCollidersWithCollider(pcWrapper, mode), func(testpc *component.PolygonCollider) bool {
		_, ok := testpc.I().Obj().Tags[tag]
		return ok
	})
}
