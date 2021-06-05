package core

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"galaxyzeta.io/engine/infra"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/linalg"
)

// InstantiateFunc receives an IGameObject2D constructor.
type InstantiateFunc func() IGameObject2D
type gameLoopStats int8

const InstantiateChannelSize = 256
const DeconstructionChannelSize = 256
const Debug = false
const (
	GameLoopStats_Created = iota
	GameLoopStats_Initialized
	GameLoopStats_Running
)

var roundRobin = 0

// AppConfig stores all user defined configs.
type AppConfig struct {
	Resolution  *linalg.Vector2i
	PhysicalFps int
	RenderFps   int
	WorkerCount int
	Title       string
	RenderFunc  func() // will be called in openGL loop.
	InitFunc    func() // will be called before we start main thread.
}

// MasterLoop communicates with OpenGL frontend to do rendering jobs, also manages all sub routines for physical updates.
// There is only one MasterLoop in each process.
type MasterLoop struct {
	// --- concurrency control
	workers      []*subLoop
	workersCount int
	initFunc     func()
	status       gameLoopStats // Describes the working status of current gameLoopController.
	// --- timing control
	physicalFPS  time.Duration // Physical update rate.
	renderFPS    time.Duration // Render update rate.
	renderTicker *time.Ticker  // Render update ticker.
	// --- synchronization
	wg      *sync.WaitGroup // wg is used for Wait() method to continue after all loops stoppped.
	sigKill chan struct{}
	running bool
}

// subLoop will handle step updates.
type subLoop struct {
	name                   string
	startTime              time.Time
	physicalTicker         *time.Ticker               // Physical update ticker.
	registerChannel        chan resourceAccessRequest // A pipeline used to register gameObjects to the pool. When calling Create from SDK, load balancing is applied to distribute a create request to this channel.
	unregisterChannel      chan resourceAccessRequest // A pipeline used to unregister gameObjects to the pool When calling Destroy from SDK, load balancing is applied to distribute a destroy request to this channel.
	processingPool         map[IGameObject2D]struct{} // A list that will be re-populate with step jobs to process before each step starts.
	sigKill                chan struct{}              // A channel used for receiving kill signal.
	synergyGates           []*cc.SynergyGate          // A set of barriers that makes goroutines wait for each other to reach a common execution entry to continue.
	shouldResetInputBuffer bool                       // Mark whether this subLoop is responsible for resetting inputBuffer.
}

// NewMasterLoop returns a new masterGameLoopController. SubworkersCount is ensured to have at least 1.
// Not thread safe, no need to do that.
func NewMasterLoop(cfg *AppConfig) *MasterLoop {
	if !atomic.CompareAndSwapInt32(&casList[Cas_CoreController], Cas_False, Cas_True) {
		panic("cannot have two masterGameLoopController in a standalone process")
	}
	if cfg.WorkerCount < 1 {
		cfg.WorkerCount = 1
	}
	main := &MasterLoop{
		status:       GameLoopStats_Created,
		physicalFPS:  time.Duration(cfg.PhysicalFps),
		renderFPS:    time.Duration(cfg.RenderFps),
		renderTicker: time.NewTicker(time.Second / time.Duration(cfg.RenderFps)),
		workers:      make([]*subLoop, 0, cfg.WorkerCount),
		workersCount: cfg.WorkerCount,
		wg:           &sync.WaitGroup{},
		sigKill:      make(chan struct{}, 1),
		initFunc:     cfg.InitFunc,
	}

	mutexList[Mutex_ScreenResolution].Lock()
	screenResolution = cfg.Resolution
	mutexList[Mutex_ScreenResolution].Unlock()

	sg := make([]*cc.SynergyGate, 0, 3)
	for i := 0; i < 3; i++ {
		sg = append(sg, cc.NewSynergyGate(int64(main.workersCount)))
	}

	for i := 0; i < cfg.WorkerCount; i++ {
		sub := main.newSubGameLoopController(sg, fmt.Sprintf("%d", i))
		main.workers = append(main.workers, sub)
	}
	main.workers[0].shouldResetInputBuffer = true

	main.status = GameLoopStats_Initialized

	coreController = main

	return main
}

// RunNoBlocking creates goroutine for each subGameLoopController to work. Not blocking.
// Not thread safe, you have no need, and should not call RunNoBlocking in concurrent execution environment.
func (g *MasterLoop) RunNoBlocking() {
	if g.status == GameLoopStats_Running {
		panic("cannot run a controller twice")
	}

	window := initOpenGL(screenResolution, title)
	g.initFunc()

	for _, worker := range g.workers {
		g.wg.Add(1)
		go worker.runSubWorker()
	}

	g.running = true
	g.status = GameLoopStats_Running

	// --- begin infinite loop
	g.wg.Add(1)
	fmt.Println("render: wg++")
	renderLoop(window, g.doRender, g.sigKill)
	g.wg.Done()
	fmt.Println("render: wg--")
	// --- infinite loop has stopped, maybe sigkill or something else
}

// Kill terminates all sub workers.
func (g *MasterLoop) Kill() {
	fmt.Println("kill")
	for _, worker := range g.workers {
		fmt.Println("emit kill")
		worker.sigKill <- struct{}{}
	}
	g.sigKill <- struct{}{} // kill openGL routine (main routine) (may panic if the channel has been closed)
	g.running = false
}

// Wait MasterLoop and all subLoops to be killed. Blocking.
func (g *MasterLoop) Wait() {
	g.wg.Wait()
}

// roundRobin selects a subLoop by round-robin strategy.
func (g *MasterLoop) roundRobin() *subLoop {
	s := g.workers[roundRobin]
	roundRobin = (roundRobin + 1) % g.workersCount
	return s
}

//____________________________________
//
// 		SubGameLoopController
//____________________________________

// newSubGameLoopController returns a subGameLoopController.
func (m *MasterLoop) newSubGameLoopController(sg []*cc.SynergyGate, name string) *subLoop {
	g := &subLoop{
		name:              name,
		registerChannel:   make(chan resourceAccessRequest, InstantiateChannelSize),
		unregisterChannel: make(chan resourceAccessRequest, DeconstructionChannelSize),
		sigKill:           make(chan struct{}, 1),
		physicalTicker:    time.NewTicker(time.Second / m.physicalFPS),
		synergyGates:      make([]*cc.SynergyGate, 0, 3),
		processingPool:    make(map[IGameObject2D]struct{}),
	}
	for i := 0; i < 3; i++ {
		g.synergyGates = append(g.synergyGates, sg[i])
	}
	return g
}

func (g *subLoop) runSubWorker() {
	g.startTime = time.Now()
	fmt.Println("sub: wg ++")
	for coreController.running {
		select {
		case <-g.physicalTicker.C:
			g.doPhysicalUpdate()
		case <-g.sigKill:
			g.subLoopExit()
			return
		}
	}
	g.subLoopExit()
}

func (g *subLoop) subLoopExit() {
	close(g.sigKill)
	fmt.Println("sub: wg --")
	coreController.wg.Done()
}

//____________________________________
//
// 		  Processor Functions
//____________________________________

func (g *MasterLoop) doRender() {
	for _, pool := range activePool {
		for gameObj := range pool {
			obj2d := gameObj.GetGameObject2D()
			if obj2d.IsVisible && obj2d.Sprite != nil {
				obj2d.Sprite.Render(obj2d.currentStats.Position.X, obj2d.currentStats.Position.Y)
			}
		}
	}
}

func (g *subLoop) doPhysicalUpdate() {
	g.synergyGates[0].Wait()
	// 1. check whether there are items to create
	for len(g.registerChannel) > 0 {
		req := <-g.registerChannel
		if req.isActive == infra.BoolPtr_True {
			fmt.Println("active construction ok ", g.name)
			g.processingPool[req.payload] = struct{}{}
		}
		addObjDefault(req.payload, *req.isActive)
	}
	g.synergyGates[1].Wait()
	// 2. do step
	for iobj2d := range g.processingPool {
		iobj2d.GetGameObject2D().Callbacks.OnStep(iobj2d)
	}
	g.synergyGates[2].Wait()

	// flush input buffer, only one subLoop can do this.
	if g.shouldResetInputBuffer {
		FlushInputBuffer()
	}
	// 3. check whether there are items to unregister
	for len(g.unregisterChannel) > 0 {
		req := <-g.unregisterChannel
		fmt.Println("sub: destroy ", g.name)
		delete(g.processingPool, req.payload)
		removeObjDefault(req.payload, req.payload.GetGameObject2D().isActive)
	}
}
