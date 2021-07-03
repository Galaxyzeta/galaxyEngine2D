package collision

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	cc "galaxyzeta.io/engine/infra/concurrency"
)

var NamePhysics2DSystem = "sys_GridCollision2D"

// GridCollision2DSystem manages all game colliders with a grid based hashset.
type GridCollision2DSystem struct {
	*base.SystemBase
	gridWidth  float64
	gridHeight float64
	colliders  [][]map[*component.PolygonCollider]struct{}
}

func (s *GridCollision2DSystem) execute(executor *cc.Executor) {

}

// func (s *GridCollision2DSystem) getOccupiedGrids(bb physics.BoundingBox) (ret []linalg.Vector2i) {
// 	topLeftGrid := s.point2Grid(bb[physics.BB_TopLeft])
// 	bottomRightGrid := s.point2Grid(bb[physics.BB_BotRight])
// 	for i := topLeftGrid.X; i < bottomRightGrid.X; i++ {
// 		for j := topLeftGrid.Y; j < bottomRightGrid.Y; j++ {
// 			ret = append(ret, linalg.Vector2i{X: i, Y: j})
// 		}
// 	}
// }

// func (s *GridCollision2DSystem) point2Grid(point linalg.Vector2f64) linalg.Vector2i {
// 	return linalg.Vector2i{X: int(point.X / s.gridWidth), Y: int(point.Y / s.gridHeight)}
// }

// ===== IMPLEMENTATION =====

func (s *GridCollision2DSystem) Execute(executor *cc.Executor) {
	// for _, item := range s.obj2data {
	executor.AsyncExecute(func() (interface{}, error) {
		// s.execute(item)
		return nil, nil
	})
	// }
}

func (s *GridCollision2DSystem) GetSystemBase() *base.SystemBase {
	return s.SystemBase
}

func (s *GridCollision2DSystem) GetName() string {
	return NamePhysics2DSystem
}

func (s *GridCollision2DSystem) Register(iobj base.IGameObject2D) {
	ipc := iobj.GetGameObject2D().GetComponent(component.NamePolygonCollider)
	if ipc == nil {
		return
	}
	// pc := ipc.(*component.PolygonCollider)
	// bb := pc.Collider.GetBoundingBox()
}

func (s *GridCollision2DSystem) Unregister(iobj base.IGameObject2D) {
}
