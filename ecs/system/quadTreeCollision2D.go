package system

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/infra/logger"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"galaxyzeta.io/engine/physics/collision"
)

var NameCollision2Dsystem = "sys_QuadCollision2D"
var quadTreeLogger = logger.New("QuadTreeCollision2D")

func init() {
	quadTreeLogger.Disable()
}

func NewQuadTreeCollision2DSystem(priority int, maintainanceArea physics.Rectangle, loadFactor int, minDivision float64) *QuadTreeCollision2DSystem {
	return &QuadTreeCollision2DSystem{
		SystemBase: base.NewSystemBase(priority),
		qt:         collision.NewQuadTree(maintainanceArea, loadFactor, minDivision),
	}
}

// QuadTreeCollision2DSystem manages all game colliders with a grid based hashset.
type QuadTreeCollision2DSystem struct {
	*base.SystemBase
	qt *collision.QuadTree
}

func (s *QuadTreeCollision2DSystem) execute(executor *cc.Executor) {
	offset := s.qt.GetLooseOffset()
	removeQueue := []*component.PolygonCollider{}
	removeNode := []*collision.QTreeNode{}
	s.qt.Traverse(func(pc *component.PolygonCollider, currentGrid physics.Rectangle, node *collision.QTreeNode) {
		if !pc.Collider.GetBoundingBox().ToRectangle().CropOutside(offset, offset).Intersect(currentGrid) {
			removeQueue = append(removeQueue, pc)
			removeNode = append(removeNode, node)
		}
		quadTreeLogger.Infof("pc-anchor = %v currentGrid = %v", pc.Collider.GetAnchor(), currentGrid)
	})
	for idx, elem := range removeQueue {
		items := removeNode[idx].GetItem()
		delete(items, elem)
		if len(items) == 0 {
			removeNode[idx].SetHasChild(true)
		}
	}
	for _, elem := range removeQueue {
		s.qt.Insert(elem)
	}
}

// ===== debug only =====
func (s *QuadTreeCollision2DSystem) Traverse(f func(pc *component.PolygonCollider, r physics.Rectangle, qn *collision.QTreeNode)) {
	s.qt.Traverse(f)
}

// ===== Functional Implementation =====

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithCollider(col component.PolygonCollider) []*component.PolygonCollider {
	col.Collider.GetBoundingBox()
	return s.QueryNeighborCollidersWithPosition(*col.Collider.GetAnchor())
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithColliderAndFilter(col component.PolygonCollider, filter func(*component.PolygonCollider) bool) []*component.PolygonCollider {
	return s.QueryNeighborCollidersWithPositionAndFilter(*col.Collider.GetAnchor(), filter)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithPosition(pos linalg.Vector2f64) []*component.PolygonCollider {
	return s.qt.Query(pos)
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithPositionAndFilter(pos linalg.Vector2f64, filter func(*component.PolygonCollider) bool) []*component.PolygonCollider {
	li := s.qt.Query(pos)
	var ret []*component.PolygonCollider
	for _, collider := range li {
		if filter(collider) {
			ret = append(ret, collider)
		}
	}
	return ret
}

func (s *QuadTreeCollision2DSystem) QueryNeighborCollidersWithRay(r physics.Ray) {

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
	ipc := iobj.GetGameObject2D().GetComponent(component.NamePolygonCollider)
	pc := ipc.(*component.PolygonCollider)
	s.qt.Insert(pc)
}

func (s *QuadTreeCollision2DSystem) Unregister(iobj base.IGameObject2D) {
}
