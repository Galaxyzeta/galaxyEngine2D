package core

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

func containsActiveDefault(obj IGameObject2D) bool {
	_, ok := activePool[Label_Default][obj]
	return ok
}

func containsInactiveDefault(obj IGameObject2D) bool {
	_, ok := inactivePool[Label_Default][obj]
	return ok
}
