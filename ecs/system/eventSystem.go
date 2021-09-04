package system

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/ecs/component"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/infra/ds"
)

const NameEventSystem = "sys_eventSystem"

type EventSystem struct {
	*base.SystemBase

	csys       collision.ICollisionSystem
	messengers ds.IndexedList
}

func (s *EventSystem) execute(_ *cc.Executor) {
	totalLen := s.messengers.Len()
	// iterate over the list
	for i, cur := 0, s.messengers.Front(); i < totalLen; i, cur = i+1, cur.Next() {
		messenger := cur.Value.(*component.Messenger)
		if col := collision.ColliderAtPolygonWithFilter(s.csys, messenger.Pc.Collider, messenger.ShouldTrigger); col != nil {
			messenger.Impact(col)
		}
		cur = cur.Next()
	}
}

// ===== IMPLEMENTATION =====

func (s *EventSystem) Execute(executor *cc.Executor) {
	s.execute(executor)
}

func (s *EventSystem) GetSystemBase() *base.SystemBase {
	return s.SystemBase
}

func (s *EventSystem) GetName() string {
	return NameEventSystem
}

func (s *EventSystem) Register(iobj base.IGameObject2D) {
	messenger := iobj.Obj().GetComponent(component.NameMessenger).(*component.Messenger)
	s.messengers.PushBack(messenger)
}

func (s *EventSystem) Unregister(iobj base.IGameObject2D) {
	messenger := iobj.Obj().GetComponent(component.NameMessenger).(*component.Messenger)
	s.messengers.Remove(messenger)
}
