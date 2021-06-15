package system

import (
	"math"

	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/linalg"
)

// PhysicalComponentWrapper wraps RigidBody2D and Transform component.
type PhysicalComponentWrapper struct {
	*component.RigidBody2D
	*component.Transform2D
}

type Physics2DSystem struct {
	data []PhysicalComponentWrapper
}

// ===== IMPLEMENTATION =====
func (s *Physics2DSystem) Execute() {
	for _, item := range s.data {
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
		// gravity effect
		// TODO
	}
}
