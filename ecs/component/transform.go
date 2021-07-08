package component

import (
	"galaxyzeta.io/engine/config"
	"galaxyzeta.io/engine/infra/concurrency/lock"
	"galaxyzeta.io/engine/linalg"
)

const NameTransform2D = "Transform2D"

type Transform2D struct {
	prevPos linalg.Vector2f64
	Pos     linalg.Vector2f64
	mu      lock.SpinLock
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
	return tf.prevPos.X
}

func (tf *Transform2D) GetPrevY() float64 {
	return tf.prevPos.Y
}

// MemXY memorizes X, Y postion to prevX, prevY.
func (tf *Transform2D) MemXY() {
	tf.prevPos = tf.Pos
}

// Transalte a delta distance.
func (tf *Transform2D) Translate(x float64, y float64) {
	tf.Pos.X += x
	tf.Pos.Y += y
}

// Teleport to a given location.
func (tf *Transform2D) Teleport(x float64, y float64) {
	tf.Pos.X = x
	tf.Pos.Y = y
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
