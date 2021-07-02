package component

import (
	"galaxyzeta.io/engine/config"
	"galaxyzeta.io/engine/infra/concurrency/lock"
)

const NameTransform2D = "Transform2D"

type Transform2D struct {
	prevX float64
	prevY float64
	X     float64
	Y     float64
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

func (tf *Transform2D) GetPrevX() float64 {
	return tf.prevX
}

func (tf *Transform2D) GetPrevY() float64 {
	return tf.prevY
}

// MemXY memorizes X, Y postion to prevX, prevY.
func (tf *Transform2D) MemXY() {
	tf.prevX = tf.X
	tf.prevY = tf.Y
}

// Transalte a delta distance.
func (tf *Transform2D) Translate(x float64, y float64) {
	tf.X += x
	tf.Y += y
}

// Teleport to a given location.
func (tf *Transform2D) Teleport(x float64, y float64) {
	tf.X = x
	tf.Y = y
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
