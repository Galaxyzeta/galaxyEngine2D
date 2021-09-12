package core

import (
	"runtime"
	"sync"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/input/keys"
)

func GlobalInitializer() {
	// must get cwd
	var err error
	// cwd, err = os.Getwd()
	cwd = "D:/Go/go/src/galaxyzeta.io/engine"
	if err != nil {
		panic(err)
	}

	// init pools
	objPoolInit(&activePool)
	objPoolInit(&inactivePool)

	// init mutexList
	mutexList = make([]*sync.RWMutex, 16)
	for idx := range mutexList {
		mutexList[idx] = &sync.RWMutex{}
	}

	// init casList
	casList = make([]int32, 8)
	for idx := range casList {
		casList[idx] = Cas_False
	}

	// init render sorting list
	renderSortList = make([]*base.GameObject2D, 0, 1024)

	// init global stuff
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
