package chrono

import (
	"time"

	"galaxyzeta.io/engine/core"
)

// Ticker is a manually controlled time accumulator that is not thread safe.
type Ticker struct {
	time time.Duration
	tick int64
}

func NewTicker() *Ticker {
	return &Ticker{
		time: 0,
		tick: 0,
	}
}

func (t *Ticker) Tick() {
	t.tick++
	t.time += core.GetPhysicsDeltaTime()
}

func (t Ticker) TickElapsed() int64 {
	return t.tick
}

func (t Ticker) TimeElapsed() time.Duration {
	return t.time
}
