package lb

import (
	"galaxyzeta.io/engine/infra/concurrency/lock"
)

type ILoadBalancer interface {
	Choose() int
}

type RoundRobin struct {
	robin int
	max   int
	mu    *lock.SpinLock
}

func NewRoundRobin(max int) *RoundRobin {
	return &RoundRobin{
		robin: 0,
		max:   max,
		mu:    &lock.SpinLock{},
	}
}

func (r *RoundRobin) Choose() int {
	res := r.robin
	r.mu.Lock()
	r.robin++
	if r.robin >= r.max {
		r.robin = 0
	}
	r.mu.Unlock()
	return res
}
