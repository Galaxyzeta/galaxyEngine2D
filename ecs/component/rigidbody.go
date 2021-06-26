package component

import (
	"galaxyzeta.io/engine/config"
	"galaxyzeta.io/engine/infra/concurrency/lock"
)

const NameRigidBody2D = "RigidBody"

type RigidBody2D struct {
	Gravity      float64
	GravityDir   float64
	Speed        float64
	Direction    float64
	Acceleration float64
	mu           lock.SpinLock
}

func NewRigidBody2D() *RigidBody2D {
	return new(RigidBody2D)
}

// GetName is an implementation of IComponent.
func (rb *RigidBody2D) GetName() string {
	return NameRigidBody2D
}

func (rb *RigidBody2D) Lock() {
	if config.EnableMultithread {
		rb.mu.Lock()
	}
}

func (rb *RigidBody2D) Unlock() {
	if config.EnableMultithread {
		rb.mu.Unlock()
	}
}
