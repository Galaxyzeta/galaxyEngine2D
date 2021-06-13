package core

import (
	"sync"

	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input/keys"
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
var cameraPool []*graphics.Camera
var renderSortList []*GameObject2D

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
var currentActiveCameraIndex int = 0
var title string
var inputBuffer []map[keys.Key]struct{}
var coreController *MasterLoop
var cwd string
