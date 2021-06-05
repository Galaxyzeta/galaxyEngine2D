package core

import (
	"galaxyzeta.io/engine/input/keys"
	"os"
	"runtime"
	"sync"
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

	// init global stuff
	sceneMap = make(map[string]*scene)
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
