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
	obj2data map[base.IGameObject2D]PhysicalComponentWrapper
}

func NewPhysics2DSystem(prioriy int) *Physics2DSystem {
	return &Physics2DSystem{
		obj2data:   make(map[base.IGameObject2D]PhysicalComponentWrapper, 64),
		SystemBase: base.NewSystemBase(prioriy),
	}
}

func (s *Physics2DSystem) execute(item PhysicalComponentWrapper) {
	var rad float64 = float64(linalg.Deg2Rad(linalg.InvertDeg(item.Direction)))
	item.Transform2D.X += item.Speed * math.Cos(rad)
	item.Transform2D.Y += item.Speed * math.Sin(rad)
	// speed attenuation
	if item.Speed > 0 {
		item.Speed -= item.Acceleration
		if item.Speed < 0 {
			item.Speed = 0
		}
	}
}

// ===== IMPLEMENTATION =====
func (s *Physics2DSystem) Execute(executor *cc.Executor) {
	for _, item := range s.obj2data {
		executor.AsyncExecute(func() (interface{}, error) {
			s.execute(item)
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

func (s *Physics2DSystem) Register(iobj base.IGameObject2D) {
	rb := iobj.GetGameObject2D().GetComponent(component.NameRigidBody2D).(*component.RigidBody2D)
	tf := iobj.GetGameObject2D().GetComponent(component.NameTransform2D).(*component.Transform2D)
	s.obj2data[iobj] = PhysicalComponentWrapper{
		RigidBody2D: rb,
		Transform2D: tf,
	}
}

func (s *Physics2DSystem) Unregister(iobj base.IGameObject2D) {
	delete(s.obj2data, iobj)
}
