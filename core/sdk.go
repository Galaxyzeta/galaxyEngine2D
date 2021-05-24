package core

import (
	"galaxyzeta.io/engine/infra"
	"galaxyzeta.io/engine/linalg"
)

// +------------------------+
// |	  	 System	 	 	|
// +------------------------+

// StartApplication starts the whole application with an option.
func StartApplication(cfg *AppConfig) {
	g := newMasterLoop(cfg)
	g.runNoBlocking()
	g.wait()
}

// ScreenResolution get current screen's resolution. It is thread-safe.
func ScreenResolution() *linalg.Vector2i {
	mutextList[mutex_ScreenResolution].Lock()
	defer mutextList[mutex_ScreenResolution].Unlock()
	return screenResolution
}

// Exit will quit the whole program.
func Exit() {
	coreController.kill()
}

// Title returns game window's title.
func Title() string {
	mutextList[mutex_Title].Lock()
	defer mutextList[mutex_Title].Unlock()
	return title
}

// ChangeScene changes current scene TODO
func ChangeScene(sceneName string) {
	currentSceneName = sceneName
}

// +------------------------+
// |	  	GameObjs	 	|
// +------------------------+

// Create will instantiate an object immediately.
// The object will be put to the global resource pool in next physical tick.
func Create(constructor func() IGameObject2D) IGameObject2D {
	return doCreate(constructor, infra.BoolPtr_True)
}

// CreateInactive will instantiate an inactive object immediately.
// The object will be put to the global resource pool in next physical tick.
func CreateInactive(constructor func() IGameObject2D) IGameObject2D {
	return doCreate(constructor, infra.BoolPtr_False)
}

// Destroy will deconstruct an active/inactive object immediately.
// The object will be truely removed from resource pool in the next physical tick.
func Destroy(obj IGameObject2D) {
	doDestroy(obj, nil)
}

// Activate an object from deactive list, if it exists in it.
func Activate(obj IGameObject2D) bool {
	if containsInactiveDefault(obj) {
		delete(inactivePool[Label_Default], obj)
		activePool[Label_Default][obj] = struct{}{}
		obj.GetGameObject2D().isActive = true
		return true
	}
	return false
}

// Deactivate an object from active list, if it exists in it.
func Deactivate(obj IGameObject2D) bool {
	if containsActiveDefault(obj) {
		delete(activePool[Label_Default], obj)
		inactivePool[Label_Default][obj] = struct{}{}
		obj.GetGameObject2D().isActive = false
		return true
	}
	return false
}

// +------------------------+
// |	  	 Input	 	 	|
// +------------------------+

func IsKeyPressed(k Key) (b bool) {
	mutextList[mutex_Keyboard_Pressed].RLock()
	b = isKeyRegistered(k, KeyPressed)
	mutextList[mutex_Keyboard_Pressed].RUnlock()
	return b
}

func IsKeyHeld(k Key) (b bool) {
	mutextList[mutex_Keyboard_Held].RLock()
	b = isKeyRegistered(k, KeyHold)
	mutextList[mutex_Keyboard_Held].RUnlock()
	return b
}

func IsKeyReleased(k Key) (b bool) {
	mutextList[mutex_Keyboard_Released].RLock()
	b = isKeyRegistered(k, KeyPressed)
	mutextList[mutex_Keyboard_Released].RUnlock()
	return b
}
