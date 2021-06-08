package core

import (
	"galaxyzeta.io/engine/graphics"
)

type Scene struct {
	cams          map[*graphics.Camera]struct{}
	currentCamera *graphics.Camera
	initFunction  func()
}

// CreateScene and add it to global resource manager.
// You should not use &Scene{} directly, that will not work at all.
func CreateScene(name string, initFunc func()) *Scene {
	sc := &Scene{
		cams:          make(map[*graphics.Camera]struct{}),
		currentCamera: nil,
		initFunction:  initFunc,
	}
	sceneMap[name] = sc
	return sc
}

// RegisterCamera to the scene.
func (sc *Scene) RegisterCamera(cam *graphics.Camera) {
	sc.cams[cam] = struct{}{}
}

// UnregisterCamera removes a camera from the scene.
// Not recommend to do this.
func (sc *Scene) UnregisterCamera(cam *graphics.Camera) {
	if cam == sc.currentCamera {
		sc.currentCamera = nil
	}
	delete(sc.cams, cam)
}

func (sc *Scene) SetCurrentCamera(cam *graphics.Camera) {
	sc.currentCamera = cam
}

func (sc *Scene) GetCurrentCamera() *graphics.Camera {
	return sc.currentCamera
}
