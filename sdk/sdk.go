package sdk

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/linalg"
)

// +------------------------+
// |	  	 System	 	 	|
// +------------------------+

// StartApplication starts the whole application with an option.
func StartApplication(cfg *core.AppConfig) {
	g := core.NewApplication(cfg)
	g.Start()
}

// ScreenResolution get current screen's resolution. It is thread-safe.
func ScreenResolution() linalg.Vector2f64 {
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
func Create(constructor func() base.IGameObject2D) base.IGameObject2D {
	return core.Create(constructor)
}

// CreateInactive will instantiate an inactive object immediately.
// The object will be put to the global resource pool in next physical tick.
func CreateInactive(constructor func() base.IGameObject2D) base.IGameObject2D {
	return core.CreateInactive(constructor)
}

// Destroy will deconstruct an active/inactive object immediately.
// The object will be truely removed from resource pool in the next physical tick.
func Destroy(iobj base.IGameObject2D) {
	core.Destroy(iobj)
}

// Activate an object from deactive list, if it exists in it.
func Activate(iobj base.IGameObject2D) bool {
	return core.Activate(iobj)
}

// Deactivate an object from active list, if it exists in it.
func Deactivate(iobj base.IGameObject2D) bool {
	return core.Deactivate(iobj)
}

func GetCamera() *graphics.Camera {
	return graphics.GetCurrentCamera()
}
