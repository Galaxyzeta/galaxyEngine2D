package core

import (
	"sync"

	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input/keys"
)

func addObjDefault(obj IGameObject2D, isActive bool) {
	var targetPool map[label]objPool
	if isActive {
		targetPool = activePool
	} else {
		targetPool = inactivePool
	}
	targetPool[Label_Default][obj] = struct{}{}
}

func removeObjDefault(obj IGameObject2D, isActive bool) bool {
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

func ContainsActiveDefault(obj IGameObject2D) bool {
	_, ok := activePool[Label_Default][obj]
	return ok
}

func ContainsInactiveDefault(obj IGameObject2D) bool {
	_, ok := inactivePool[Label_Default][obj]
	return ok
}

// GetCoreController retrieves the central game loop controller for you.
func GetCoreController() *MasterLoop {
	return coreController
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

func GetCamera(index int) *graphics.Camera {
	return cameraPool[index]
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
