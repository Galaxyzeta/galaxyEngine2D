package system

import (
	"fmt"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/ecs/component"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/infra/logger"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

var NameCollision2Dsystem = "sys_QuadCollision2D"
var quadTreeLogger = logger.New("QuadTreeCollision2D")

func init() {
	// quadTreeLogger.Disable()
}

func NewQuadTreeCollision2DSystem(priority int, maintainanceArea physics.Rectangle, loadFactor int, minDivision float64) *QuadTreeCollision2DSystem {
	return &QuadTreeCollision2DSystem{
		SystemBase: base.NewSystemBase(priority),
		qt:         collision.NewQuadTree(maintainanceArea, loadFactor, minDivision),
	}
}

// QuadTreeCollision2DSystem manages all game colliders with a quad tree.
// It provides ability to quickly locate colliders that might have a chance to collide.
type QuadTreeCollision2DSystem struct {
	*base.SystemBase
	qt *collision.QuadTree
}

func (s *QuadTreeCollision2DSystem) execute(executor *cc.Executor) {
	rmPolygonColliders := []*component.PolygonCollider{}
	rmNodes := []*collision.QTreeNode{}
	s.qt.Traverse(func(pc *component.PolygonCollider, node *collision.QTreeNode, at collision.AreaType, idx int) bool {
		if at == collision.Inline {
			// inline object: checks intersection with its child nodes.
			if val := node.GetIntersectedSection(pc.Collider.GetBoundingBox()); val >= 0 {
				rmNodes = append(rmNodes, node)
				rmPolygonColliders = append(rmPolygonColliders, pc)
			}
		} else {
			// not inline object: checks intersection with its currently related nodes
			pcRect := pc.Collider.GetBoundingBox().ToRectangle()
			if !pcRect.Intersect(node.GetArea()) {
				if pc.I().Obj().Name == "player" {
					fmt.Print("")
				}
				rmPolygonColliders = append(rmPolygonColliders, pc)
				rmNodes = append(rmNodes, node)
			}
		}
		return false
	})
	for idx, elem := range rmPolygonColliders {
		rmNodes[idx].Delete(elem)
	}
	for _, elem := range rmPolygonColliders {
		s.qt.Insert(elem)
	}
}

// ===== debug only =====
func (s *QuadTreeCollision2DSystem) Traverse(needLock bool, f collision.QTreeTraverseFunc) {
	if needLock {
		s.qt.TraverseWithLock(f)
	} else {
		s.qt.Traverse(f)
	}
}

// ===== Functional Implementation =====

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithCollider(col component.PolygonCollider, mode collision.QueryMode) []*component.PolygonCollider {
	col.Collider.GetBoundingBox()
	return s.QueryNeighborCollidersWithRect(col.Collider.GetBoundingBox().ToRectangle(), mode)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithColliderAndFilter(col component.PolygonCollider, filter func(*component.PolygonCollider) bool, mode collision.QueryMode) []*component.PolygonCollider {
	return s.QueryNeighborCollidersWithPositionAndFilter(*col.Collider.GetAnchor(), filter, mode)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithPosition(pos linalg.Vector2f64, mode collision.QueryMode) []*component.PolygonCollider {
	return s.qt.QueryByPoint(pos, mode)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithRect(r physics.Rectangle, mode collision.QueryMode) []*component.PolygonCollider {
	return s.qt.QueryByRect(r, mode)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithRay(r physics.Ray, mode collision.QueryMode) []*component.PolygonCollider {
	return s.qt.QueryByRay(r, mode)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithPositionAndFilter(pos linalg.Vector2f64, filter func(*component.PolygonCollider) bool, mode collision.QueryMode) []*component.PolygonCollider {
	li := s.qt.QueryByPoint(pos, mode)
	var ret []*component.PolygonCollider
	for _, collider := range li {
		if filter(collider) {
			ret = append(ret, collider)
		}
	}
	return ret
}

// ===== IMPLEMENTATION =====

func (s *QuadTreeCollision2DSystem) Execute(executor *cc.Executor) {
	s.execute(executor)
}

func (s *QuadTreeCollision2DSystem) GetSystemBase() *base.SystemBase {
	return s.SystemBase
}

func (s *QuadTreeCollision2DSystem) GetName() string {
	return NameCollision2Dsystem
}

func (s *QuadTreeCollision2DSystem) Register(iobj base.IGameObject2D) {
	ipc := iobj.Obj().GetComponent(component.NamePolygonCollider)
	pc := ipc.(*component.PolygonCollider)
	s.qt.Insert(pc)
}

func (s *QuadTreeCollision2DSystem) Unregister(iobj base.IGameObject2D) {
	testpc := iobj.Obj().GetComponent(component.NamePolygonCollider)
	s.qt.TraverseWithLock(func(pc *component.PolygonCollider, qn *collision.QTreeNode, _ collision.AreaType, _ int) bool {
		if testpc == pc {
			qn.UnsafeDelete(pc)
			return true
		}
		return false
	})
}

func (s *QuadTreeCollision2DSystem) Activate(iobj base.IGameObject2D) {
	ipc := iobj.Obj().GetComponent(component.NamePolygonCollider)
	pc := ipc.(*component.PolygonCollider)
	s.qt.Activate(pc)
}

func (s *QuadTreeCollision2DSystem) Deactivate(iobj base.IGameObject2D) {
	ipc := iobj.Obj().GetComponent(component.NamePolygonCollider)
	pc := ipc.(*component.PolygonCollider)
	s.qt.Deactivate(pc)
}
