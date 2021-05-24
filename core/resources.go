package core

import (
	"galaxyzeta.io/engine/graphics"
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

const (
	mutex_ScreenResolution uint8 = iota
	mutex_Title
	mutex_Keyboard_Pressed
	mutex_Keyboard_Held
	mutex_Keyboard_Released
)

var mutextList []sync.RWMutex

// +------------------------+
// |	    CAS var 	 	|
// +------------------------+

const (
	cas_true  int32 = 1
	cas_false int32 = 0
)
const (
	cas_coreController = iota
)

var casList []int32

// +------------------------+
// |	  Input Const 	 	|
// +------------------------+

const (
	KeyRelease = iota
	KeyPressed
	KeyHold
)

// +------------------------+
// |	   Global Var 	 	|
// +------------------------+

var currentSceneName string
var screenResolution *linalg.Vector2i
var title string
var inputBuffer []map[int]struct{}
var coreController *masterLoop
