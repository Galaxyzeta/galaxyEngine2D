package core

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/infra"
)

// doCreate does actual creation.
func doCreate(constructor func() base.IGameObject2D, isActive *bool) base.IGameObject2D {
	obj := constructor()
	obj.Obj().SetIGameObject2D(obj)
	app.registerChannel <- resourceAccessRequest{
		payload:  obj,
		isActive: isActive,
	}
	return obj
}

func doDestroy(obj base.IGameObject2D, isActive *bool) {
	obj2d := obj.Obj()
	if obj2d.Callbacks.OnDestroy != nil {
		obj2d.Callbacks.OnDestroy(obj)
	}
	app.unregisterChannel <- resourceAccessRequest{
		payload:  obj,
		isActive: isActive,
	}
}

// +------------------------+
// |	  	GameObjs	 	|
// +------------------------+

// Create will instantiate an object immediately.
// The object will be put to the global resource pool in next physical tick.
func Create(constructor func() base.IGameObject2D) base.IGameObject2D {
	return doCreate(constructor, infra.BoolPtr_True)
}

// CreateInactive will instantiate an inactive object immediately.
// The object will be put to the global resource pool in next physical tick.
func CreateInactive(constructor func() base.IGameObject2D) base.IGameObject2D {
	return doCreate(constructor, infra.BoolPtr_False)
}

// Destroy will deconstruct an active/inactive object immediately.
// The object will be truely removed from resource pool in the next physical tick.
func Destroy(obj base.IGameObject2D) {
	doDestroy(obj, nil)
}

//
func GetIGameobjects() (ret []base.IGameObject2D) {
	mu := mutexList[Mutex_ActivePool]
	mu.RLock()
	for _, pool := range activePool {
		for iobj, _ := range pool {
			ret = append(ret, iobj)
		}
	}
	mu.RUnlock()
	return ret
}

func GetIGameobjectsMap() (ret map[base.IGameObject2D]struct{}) {
	mu := mutexList[Mutex_ActivePool]
	ret = make(map[base.IGameObject2D]struct{})
	mu.RLock()
	for _, pool := range activePool {
		for iobj, _ := range pool {
			ret[iobj] = struct{}{}
		}
	}
	mu.RUnlock()
	return ret
}

// Activate an object from deactive list, if it exists in it.
func Activate(iobj base.IGameObject2D) bool {
	if ContainsInactiveDefault(iobj) {
		doActivate(iobj)
		delete(inactivePool[Label_Default], iobj)
		activePool[Label_Default][iobj] = struct{}{}
		iobj.Obj().IsActive = true
		systemLogger.Debugf("activate %v", iobj.Obj().Name)
		return true
	}
	return false
}

// Deactivate an object from active list, if it exists in it.
func Deactivate(iobj base.IGameObject2D) bool {
	if ContainsActiveDefault(iobj) {
		doDeactivate(iobj)
		delete(activePool[Label_Default], iobj)
		inactivePool[Label_Default][iobj] = struct{}{}
		iobj.Obj().IsActive = false
		systemLogger.Debugf("deactivate %v", iobj.Obj().Name)
		return true
	}
	return false
}

func doActivate(iobj base.IGameObject2D) {
	for _, sys := range iobj.Obj().GetSubscribedSystemMap() {
		sys.Activate(iobj)
	}
}

func doDeactivate(iobj base.IGameObject2D) {
	for _, sys := range iobj.Obj().GetSubscribedSystemMap() {
		sys.Deactivate(iobj)
	}
}
