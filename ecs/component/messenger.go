package component

import "galaxyzeta.io/engine/base"

const NameMessenger = "Messenger"

// Messenger delivers events.
type Messenger struct {
	Owner         base.IGameObject2D
	Pc            *PolygonCollider
	ShouldTrigger func(pc *PolygonCollider) bool
	Impact        func(pc *PolygonCollider)
}

func (Messenger) GetName() string {
	return NameMessenger
}
