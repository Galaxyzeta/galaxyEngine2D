package base

import (
	cc "galaxyzeta.io/engine/infra/concurrency"
)

type IComponent interface {
	GetName() string
}

type ISystem interface {
	Execute(*cc.Executor)
	GetSystemBase() *SystemBase // GetSystemBase gets the basic element of a system.
	GetName() string
	Register(IGameObject2D)
	Unregister(IGameObject2D)
	Activate(IGameObject2D)
	Deactivate(IGameObject2D)
}

type SystemBase struct {
	priority  int
	isEnabled bool
}

func (s *SystemBase) GetPriority() int {
	return s.priority
}

func (s *SystemBase) Enable() {
	s.isEnabled = true
}

func (s *SystemBase) Disable() {
	s.isEnabled = false
}

func (s *SystemBase) IsEnabled() bool {
	return s.isEnabled
}

func NewSystemBase(priority int) *SystemBase {
	return &SystemBase{
		priority: priority,
	}
}

func NewGraphicalSystemBase(priority int) *SystemBase {
	return &SystemBase{
		priority: priority,
	}
}
