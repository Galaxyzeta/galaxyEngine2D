package core

import (
	"sync"
	"time"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/input/keys"
)

func addObjDefault(iobj base.IGameObject2D, isActive bool) {
	var targetPool map[label]objPool
	var muEnum MutexIndex
	if isActive {
		targetPool = activePool
		muEnum = Mutex_ActivePool
	} else {
		targetPool = inactivePool
		muEnum = Mutex_InactivePool
	}
	mutexList[muEnum].Lock()
	targetPool[Label_Default][iobj] = struct{}{}
	mutexList[muEnum].Unlock()
	// -- register system --
	// this must be an delayed execution
	// because when loading from external level file, the position of
	// the object is set after calling constructor.
	// If we register to system immediately in ctor, will cause system
	// to register related component using the obj's zero position.
	obj := iobj.Obj()
	for _, sys := range obj.GetSubscribedSystemMap() {
		sys.Register(iobj)
	}
}

// ===== Render List =====

func removeObjDefault(obj base.IGameObject2D, isActive bool) bool {
	var targetPool map[label]objPool
	if isActive {
		targetPool = activePool
	} else {
		targetPool = inactivePool
	}
	_, ok := targetPool[Label_Default][obj]
	if !ok {
		return false
	}
	delete(targetPool[Label_Default], obj)
	return true
}

func ContainsActiveDefault(obj base.IGameObject2D) bool {
	_, ok := activePool[Label_Default][obj]
	return ok
}

func ContainsInactiveDefault(obj base.IGameObject2D) bool {
	_, ok := inactivePool[Label_Default][obj]
	return ok
}

// GetCoreController retrieves the central game loop controller for you.
func GetCoreController() *Application {
	return app
}

func GetTitle() string {
	mutexList[Mutex_Title].RLock()
	defer mutexList[Mutex_Title].RUnlock()
	return title
}

func SetTitle(alt string) {
	mutexList[Mutex_Title].RLock()
	defer mutexList[Mutex_Title].RUnlock()
	title = alt
}

func GetCurrentSceneName() string {
	mutexList[Mutex_SceneName].RLock()
	defer mutexList[Mutex_SceneName].RUnlock()
	return currentSceneName
}

func SetCurrentSceneName(newName string) {
	mutexList[Mutex_SceneName].RLock()
	defer mutexList[Mutex_SceneName].RUnlock()
	currentSceneName = newName
}

func GetRWMutex(index MutexIndex) *sync.RWMutex {
	return mutexList[index]
}

func mapActionType2Mutex(actionType keys.Action) MutexIndex {
	switch actionType {
	case keys.Action_KeyHold:
		return Mutex_Keyboard_Held
	case keys.Action_KeyPress:
		return Mutex_Keyboard_Pressed
	case keys.Action_KeyRelease:
		return Mutex_Keyboard_Released
	}
	panic("unknown mapping")
}

// SetInputBuffer sets action and key binding to inputBuffer.
// Thread-safe.
func SetInputBuffer(actionType keys.Action, key keys.Key) {
	mu := mutexList[mapActionType2Mutex(actionType)]
	mu.Lock()
	inputBuffer[actionType][key] = struct{}{}
	mu.Unlock()
}

// UnsetInputBuffer removes action and key binding from inputBuffer.
// Thread-safe.
func UnsetInputBuffer(actionType keys.Action, key keys.Key) {
	mu := mutexList[mapActionType2Mutex(actionType)]
	mu.Lock()
	delete(inputBuffer[actionType], key)
	mu.Unlock()
}

// IsSetInputBuffer checks whether the key and action binding has been set.
// Thread-safe.
func IsSetInputBuffer(actionType keys.Action, key keys.Key) bool {
	mu := mutexList[mapActionType2Mutex(actionType)]
	mu.RLock()
	_, ok := inputBuffer[actionType][key]
	mu.RUnlock()
	return ok
}

// autoResetStatusList will be used in flushInputBuffer.
// Only provided buffer field will be erased.
var autoResetStatusList []int = []int{keys.Action_KeyPress, keys.Action_KeyRelease}

// FlushInputBuffer resets input buffer to zero status except for the keyboard held status.
// That status will be automatically cancelled when a KeyRelease callback is hit.
func FlushInputBuffer() {
	for _, actionName := range autoResetStatusList {
		bufferField := inputBuffer[actionName]
		for k := range bufferField {
			delete(bufferField, k)
		}
	}
}

// GetCwd gets current working directory.
func GetCwd() string {
	return cwd
}

func GetCursorPos() (x float64, y float64) {
	mutexList[Mutex_CursorPos].RLock()
	x = cursorX
	y = cursorY
	mutexList[Mutex_CursorPos].RUnlock()
	return
}

// RegisterCtor binds object name with its constructor. This methods
// should be called only in init(), and only by provided generated code.
// Will panic if dupicated entry was found.
func RegisterCtor(name string, ctor func() base.IGameObject2D) {
	if _, exist := ctorRegistry[name]; !exist {
		ctorRegistry[name] = ctor
	} else {
		panic("duplicate constructor entry was found")
	}
}

func GetCtor(name string) func() base.IGameObject2D {
	ctor, ok := ctorRegistry[name]
	if !ok {
		panic("ctor entry with provided name does not exist")
	}
	return ctor
}

func GetPhysicsDeltaTime() (ret time.Duration) {
	mu := mutexList[Mutex_PhysicDeltaTime]
	mu.RLock()
	ret = physicalDeltaTime
	mu.RUnlock()
	return
}

// poolMapReplica
func poolMapReplica(orig map[label]objPool) (ret map[label]objPool) {
	ret = make(map[label]objPool, len(orig))
	for k, v := range orig {
		ret[k] = poolReplica(v)
	}
	return
}

// poolReplica replicates a objPool. Not thread safe.
func poolReplica(orig objPool) (ret objPool) {
	ret = make(objPool, len(orig))
	for k, v := range orig {
		ret[k] = v
	}
	return
}
