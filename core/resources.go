package core

import (
	"sync"

	"galaxyzeta.io/engine/base"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/input/keys"
)

// +------------------------+
// |	    Type Def	 	|
// +------------------------+

// resourceAccessRequest is sent to a channel for user's object creation / deconstruction request.
type resourceAccessRequest struct {
	payload  base.IGameObject2D
	isActive *bool
}
type objPool map[base.IGameObject2D]struct{}
type label string

// +------------------------+
// |	    Labels 		 	|
// +------------------------+

const Label_Default = "default"

// +------------------------+
// |	     Pools 		 	|
// +------------------------+

var activePool map[label]objPool
var inactivePool map[label]objPool
var labelPool map[label]struct{}
var sceneMap map[string]*Scene
var renderSortList []*base.GameObject2D // this array is a stash used for depth base layer sorting.

var routinePool *cc.Executor

var systemPriorityList []base.ISystem = make([]base.ISystem, 0, 256)
var gfxSystemPriorityList []base.ISystem = make([]base.ISystem, 0, 256)

var system2Priority map[base.ISystem]int = make(map[base.ISystem]int)
var name2System map[string]base.ISystem = make(map[string]base.ISystem)

var ctorRegistry map[string]func() base.IGameObject2D = make(map[string]func() base.IGameObject2D)

const MaxRenderListSize = 256

// +------------------------+
// |	     Mutex 		 	|
// +------------------------+

type MutexIndex uint8

const (
	Mutex_ScreenResolution MutexIndex = iota
	Mutex_Title
	Mutex_SceneName
	Mutex_Keyboard_Pressed
	Mutex_Keyboard_Held
	Mutex_Keyboard_Released
	Mutex_ActivePool
	Mutex_InactivePool
	Mutex_System
	Mutex_CursorPos
)

var mutexList []*sync.RWMutex

// +------------------------+
// |	    CAS var 	 	|
// +------------------------+

const (
	Cas_True  int32 = 1
	Cas_False int32 = 0
)
const (
	Cas_CoreController = iota
)

var casList []int32

// +------------------------+
// |	  Input Const 	 	|
// +------------------------+

const (
	KeyRelease = iota
	KeyPress
	KeyHold
)

// +------------------------+
// |	  	Cursor 		 	|
// +------------------------+
var cursorX float64
var cursorY float64

// +------------------------+
// |	    Render   	 	|
// +------------------------+

var RenderCmdChan chan func() = make(chan func(), 1024)

// +------------------------+
// |	   Global Var 	 	|
// +------------------------+

var currentSceneName string             // current active scene.
var currentActiveCameraIndex int = 0    // current main camera index.
var title string                        // the name of your window.
var inputBuffer []map[keys.Key]struct{} // inputBuffer stores input event type, and its corresponding key status.
var app *Application                    // app is the very core of your whole application.
var cwd string                          // cwd is a shorthand for current working directory.
