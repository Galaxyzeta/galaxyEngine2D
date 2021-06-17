package system

import (
	"math"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/linalg"
)

var NamePhysics2DSystem = "sys_Physics2D"

// PhysicalComponentWrapper wraps RigidBody2D and Transform component.
type PhysicalComponentWrapper struct {
	*component.RigidBody2D
	*component.Transform2D
}

type Physics2DSystem struct {
	*base.SystemBase
	data []PhysicalComponentWrapper
}

func NewPhysics2DSystem(prioriy int) *Physics2DSystem {
	return &Physics2DSystem{
		data:       make([]PhysicalComponentWrapper, 64),
		SystemBase: base.NewSystemBase(prioriy),
	}
}

func JoinSystem(iobj base.IGameObject2D) {
	// rb := iobj.GetGameObject2D().GetComponent(component.NameRigidBody2D).(*component.RigidBody2D)
	// tf := iobj.GetGameObject2D().GetComponent(component.NameTransform2D).(*component.Transform2D)

}

// ===== IMPLEMENTATION =====
func (s *Physics2DSystem) Execute(executor *cc.Executor) {
	for _, item := range s.data {
		executor.AsyncExecute(func() (interface{}, error) {
			var rad float64 = float64(linalg.Deg2Rad(item.Direction))
			item.Transform2D.X += item.Speed * float32(math.Cos(rad))
			item.Transform2D.Y += item.Speed * float32(math.Sin(rad))
			// speed attenuation
			if item.Speed > 0 {
				item.Speed -= item.Acceleration
				if item.Speed < 0 {
					item.Speed = 0
				}
			}
			return nil, nil
		})
	}
}

func (s *Physics2DSystem) GetSystemBase() *base.SystemBase {
	return s.SystemBase
}

func (s *Physics2DSystem) GetName() string {
	return NamePhysics2DSystem
}
