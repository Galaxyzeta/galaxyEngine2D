package system

import (
	"log"
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
	linkedList := item.RigidBody2D.GetSpeedList()
	var dx, dy float64
	for element := linkedList.Front(); element != nil; element = element.Next() {
		val := element.Value.(component.SpeedVector)
		deg := linalg.Deg2Rad(linalg.InvertDeg(val.Direction))
		dx += val.Speed * math.Cos(deg)
		dy += val.Speed * math.Sin(deg)
		// do speed atten
		if val.Speed > 0 {
			val.Speed -= val.Acceleration
			if val.Speed < 0 {
				linkedList.Remove(element)
				continue
			}
		}
		element.Value = val
	}
	// constant gravity effect
	if item.UseGravity {
		log.Println("use gravity")
		gdeg := linalg.Deg2Rad(linalg.InvertDeg(item.GravityVector.Direction))
		dx += item.GravityVector.Speed * math.Cos(gdeg)
		dy += item.GravityVector.Speed * math.Sin(gdeg)
		item.GravityVector.Speed += item.GravityVector.Acceleration
	}

	item.Transform2D.X += dx
	item.Transform2D.Y += dy

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
