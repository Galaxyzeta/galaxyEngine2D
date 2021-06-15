package cc

import (
	"fmt"
	"sync"

	"galaxyzeta.io/engine/infra/lb"
)

type Executor struct {
	size         int
	jobChannel   []chan Future
	loadBalancer lb.ILoadBalancer
	wg           sync.WaitGroup
	isRunning    bool
}

type Future struct {
	Result interface{}
	Err    error
	Fn     func() (interface{}, error)
	Done   bool
}

func NewExecutor(size int) *Executor {
	return &Executor{
		size:         size,
		jobChannel:   make([]chan Future, size),
		loadBalancer: lb.NewRoundRobin(size),
		wg:           sync.WaitGroup{},
	}
}

// Run the executor.
func (e *Executor) Run() {
	e.isRunning = true
	for i := 0; i < e.size; i++ {
		e.wg.Add(1)
		go e.jobExecutorRoutine(i, e.wg)
	}
}

func (e *Executor) jobExecutorRoutine(id int, wg sync.WaitGroup) {

	defer func() {
		p := recover()
		if p != nil {
			fmt.Println("[Fatal] jobExecutorRoutine catches a panic: %v", p)
			go e.jobExecutorRoutine(id, wg)
		}
	}()

	for e.isRunning {
		executionCtx := <-e.jobChannel[id]
		res, err := executionCtx.Fn()
		executionCtx.Done = true
		executionCtx.Err = err
		executionCtx.Result = res
	}
	wg.Done()
}

// Shutdown executor pool, waiting for all goroutines to finish.
func (e *Executor) Shutdown() {
	e.isRunning = false
	e.wg.Wait()
}

// Execute a job asynchronously, returns a future object.
func (e *Executor) Execute(fn func() (interface{}, error)) *Future {
	future := Future{
		Result: nil,
		Err:    nil,
		Fn:     fn,
		Done:   false,
	}
	e.jobChannel[e.loadBalancer.Choose()] <- future
	return &future
}
