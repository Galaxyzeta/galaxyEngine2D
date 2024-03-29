package cc

import (
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
	wg     *sync.WaitGroup
}

func NewExecutor(size int) *Executor {
	jobChannel := make([]chan Future, size) // TODO turn this into user defined config
	for idx, _ := range jobChannel {
		jobChannel[idx] = make(chan Future, 256)
	}

	return &Executor{
		size:         size,
		jobChannel:   jobChannel,
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

	// defer func() {
	// 	p := recover()
	// 	if p != nil {
	// 		fmt.Println("[Fatal] jobExecutorRoutine catches a panic: %v", p)
	// 		if e.isRunning == false {
	// 			return
	// 		}
	// 		go e.jobExecutorRoutine(id, wg)
	// 	}
	// }()

	for e.isRunning {
		executionCtx := <-e.jobChannel[id]
		res, err := executionCtx.Fn()
		executionCtx.Done = true
		executionCtx.Err = err
		executionCtx.Result = res
		executionCtx.wg.Done()
	}
	wg.Done()
}

// Shutdown executor pool, waiting for all goroutines to finish.
func (e *Executor) Shutdown() {
	e.isRunning = false
	for _, subChan := range e.jobChannel {
		close(subChan)
	}
	e.wg.Wait()
}

// AsyncExecute a job asynchronously, returns a future object.
func (e *Executor) AsyncExecute(fn func() (interface{}, error), wg *sync.WaitGroup) *Future {
	wg.Add(1)
	future := Future{
		Result: nil,
		Err:    nil,
		Fn:     fn,
		Done:   false,
		wg:     wg,
	}
	e.jobChannel[e.loadBalancer.Choose()] <- future
	return &future
}
