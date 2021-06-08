package sdk

import (
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/linalg"
)

// +------------------------+
// |	  	 System	 	 	|
// +------------------------+

// StartApplication starts the whole application with an option.
func StartApplication(cfg *core.AppConfig) {
	g := core.NewMasterLoop(cfg)
	g.RunNoBlocking()
}

// ScreenResolution get current screen's resolution. It is thread-safe.
func ScreenResolution() linalg.Vector2f32 {
	return graphics.GetScreenResolution()
}

// Exit will quit the whole program.
func Exit() {
	core.GetCoreController().Kill()
}

// Title returns game window's title.
func Title() string {
	return core.GetTitle()
}

// ChangeScene changes current scene TODO
func ChangeScene(sceneName string) {
	core.SetCurrentSceneName(sceneName)
}

// +------------------------+
// |	  	GameObjs	 	|
// +------------------------+

// Create will instantiate an object immediately.
// The object will be put to the global resource pool in next physical tick.
func Create(constructor func() core.IGameObject2D) core.IGameObject2D {
	return core.Create(constructor)
}

// CreateInactive will instantiate an inactive object immediately.
// The object will be put to the global resource pool in next physical tick.
func CreateInactive(constructor func() core.IGameObject2D) core.IGameObject2D {
	return core.CreateInactive(constructor)
}

// Destroy will deconstruct an active/inactive object immediately.
// The object will be truely removed from resource pool in the next physical tick.
func Destroy(iobj core.IGameObject2D) {
	core.Destroy(iobj)
}

// Activate an object from deactive list, if it exists in it.
func Activate(iobj core.IGameObject2D) bool {
	return core.Activate(iobj)
}

// Deactivate an object from active list, if it exists in it.
func Deactivate(iobj core.IGameObject2D) bool {
	return core.Deactivate(iobj)
}
