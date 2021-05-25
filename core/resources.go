package core

import (
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/linalg"
	"sync"
)

// +------------------------+
// |	    Type Def	 	|
// +------------------------+

// resourceAccessRequest is sent to a channel for user's object creation / deconstruction request.
type resourceAccessRequest struct {
	payload  IGameObject2D
	isActive *bool
}
type objPool map[IGameObject2D]struct{}
type renderablePool map[graphics.RenderContext]struct{}
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
var renderPool []renderablePool
var sceneMap map[string]*scene

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
// |	   Global Var 	 	|
// +------------------------+

var currentSceneName string
var screenResolution *linalg.Vector2i
var title string
var inputBuffer []map[keys.Key]struct{}
var coreController *MasterLoop
