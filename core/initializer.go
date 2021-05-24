package core

import "sync"

func GlobalInitializer() {
	// init pools
	objPoolInit(&activePool)
	objPoolInit(&inactivePool)

	// init mutexList
	mutextList = make([]sync.RWMutex, 8)
	for idx := range mutextList {
		mutextList[idx] = sync.RWMutex{}
	}

	// init casList
	casList = make([]int32, 8)
	for idx := range mutextList {
		casList[idx] = cas_false
	}

	// init global stuff
	sceneMap = make(map[string]*scene)
	currentSceneName = ""
	inputBuffer = make([]map[int]struct{}, 3)
	inputBuffer[KeyPressed] = map[int]struct{}{}
	inputBuffer[KeyHold] = map[int]struct{}{}
	inputBuffer[KeyRelease] = map[int]struct{}{}
}

// objPoolInit inits a map[label]objPool. Reduce duplicated code.
func objPoolInit(target *map[label]objPool) {
	*target = make(map[label]objPool)
	(*target)[Label_Default] = make(objPool)
}
