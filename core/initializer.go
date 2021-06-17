package core

import (
	"os"
	"runtime"
	"sync"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/linalg"
)

func GlobalInitializer() {
	// must get cwd
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// init pools
	objPoolInit(&activePool)
	objPoolInit(&inactivePool)

	// init mutexList
	mutexList = make([]*sync.RWMutex, 8)
	for idx := range mutexList {
		mutexList[idx] = &sync.RWMutex{}
	}

	// init casList
	casList = make([]int32, 8)
	for idx := range mutexList {
		casList[idx] = Cas_False
	}

	// init camera list
	cameraPool = make([]*graphics.Camera, 1, 4)
	cameraPool[0] = &graphics.Camera{
		Pos: linalg.Point2f32{
			X: 0,
			Y: 0,
		},
		Resolution: linalg.Vector2f32{
			X: 640,
			Y: 480,
		},
	}

	// init render sorting list
	renderSortList = make([]*base.GameObject2D, 0, 1024)

	// init global stuff
	sceneMap = make(map[string]*Scene)
	currentSceneName = ""
	inputBuffer = make([]map[keys.Key]struct{}, 3)
	inputBuffer[KeyPress] = map[keys.Key]struct{}{}
	inputBuffer[KeyHold] = map[keys.Key]struct{}{}
	inputBuffer[KeyRelease] = map[keys.Key]struct{}{}

	// must lock os thread
	runtime.LockOSThread()
}

// objPoolInit inits a map[label]objPool. Reduce duplicated code.
func objPoolInit(target *map[label]objPool) {
	*target = make(map[label]objPool)
	(*target)[Label_Default] = make(objPool)
}
