package component

import (
	"container/list"

	"galaxyzeta.io/engine/config"
	"galaxyzeta.io/engine/infra/concurrency/lock"
)

const NameRigidBody2D = "RigidBody"

type SpeedVector struct {
	Acceleration float64
	Direction    float64
	Speed        float64
}

type RigidBody2D struct {
	UseGravity    bool
	GravityVector SpeedVector
	speed         *list.List
	Vspeed        float64
	Hspeed        float64

	mu lock.SpinLock
}

func NewRigidBody2D() *RigidBody2D {
	return &RigidBody2D{
		speed: &list.List{},
		mu:    lock.SpinLock{},
	}
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

func (rb *RigidBody2D) AddForce(sv SpeedVector) *list.Element {
	return rb.speed.PushBack(sv)
}

// RemoveForce removes the force which is a representation of list node. If force node is nil, will do nothing.
func (rb *RigidBody2D) RemoveForce(forceNode *list.Element) {
	if forceNode == nil {
		return
	}
	rb.speed.Remove(forceNode)
}

func (rb *RigidBody2D) SetGravity(dir float64, g float64) {
	rb.GravityVector.Direction = dir
	rb.GravityVector.Acceleration = g
}

func (rb *RigidBody2D) GetSpeedList() *list.List {
	return rb.speed
}

func (rb *RigidBody2D) GetHspeed() float64 {
	return rb.Hspeed
}

func (rb *RigidBody2D) GetVspeed() float64 {
	return rb.Vspeed
}

// SetHspeed should only be called by system because it is a calculation property.
func (rb *RigidBody2D) SetHspeed(hspeed float64) {
	rb.Hspeed = hspeed
}

// SetVspeed should only be called by system because it is a calculation property.
func (rb *RigidBody2D) SetVspeed(vspeed float64) {
	rb.Vspeed = vspeed
}
