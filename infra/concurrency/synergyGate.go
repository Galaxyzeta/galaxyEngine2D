// cc is an abbr. for concurrency-control.
package cc

import (
	"sync"
)

// SynergyGate
type SynergyGate struct {
	max     int64
	current int64
	watch   chan struct{}
	mu      sync.RWMutex
}

// NewSynergyGate creates a synergy gate that allows goroutine to wait for each other to reach a synchronized break point.
func NewSynergyGate(max int64) *SynergyGate {
	return &SynergyGate{
		max:     max,
		current: 0,
		watch:   make(chan struct{}, 1),
		mu:      sync.RWMutex{},
	}
}

// Wait at a SynergyGate for all companion gouroutine to converge.
func (sg *SynergyGate) Wait() {

	sg.mu.Lock()
	sg.current++
	watch := sg.watch // save to local to prevent data race
	cnt := sg.current // save to local to prevent data race
	sg.mu.Unlock()

	if cnt > sg.max {
		panic("synergyGate wait called more than max count !")
	} else if cnt < sg.max {
		<-watch
		return
	}
	// else, breakpoint reached
	sg.mu.Lock()

	sg.current = 0
	close(sg.watch)
	sg.watch = make(chan struct{}, 1)

	sg.mu.Unlock()
}
