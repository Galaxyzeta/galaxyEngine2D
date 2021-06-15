package component

import (
	"galaxyzeta.io/engine/config"
	"galaxyzeta.io/engine/infra/concurrency/lock"
)

const NameTransform2D = "Transform2D"

type Transform2D struct {
	prevX float32
	prevY float32
	X     float32
	Y     float32
	mu    lock.SpinLock
}

func NewTransform2D() *Transform2D {
	return new(Transform2D)
}

// ===== IMPLEMENTATION =====
// GetName is an implementation of IComponent.
func (tf *Transform2D) GetName() string {
	return NameTransform2D
}

// ===== PUBLIC METHOD =====

func (tf *Transform2D) GetPrevX() float32 {
	return tf.prevX
}

func (tf *Transform2D) GetPrevY() float32 {
	return tf.prevY
}

// MemXY memorizes X, Y postion to prevX, prevY.
func (tf *Transform2D) MemXY() {
	tf.prevX = tf.X
	tf.prevY = tf.Y
}

// ===== LOCK METHODS =====

func (tf *Transform2D) Lock() {
	if config.EnableMultithread {
		tf.mu.Lock()
	}
}

func (tf *Transform2D) Unlock() {
	if config.EnableMultithread {
		tf.mu.Unlock()
	}
}
