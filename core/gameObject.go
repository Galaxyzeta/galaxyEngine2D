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

// Activate an object from deactive list, if it exists in it.
func Activate(obj base.IGameObject2D) bool {
	if ContainsInactiveDefault(obj) {
		delete(inactivePool[Label_Default], obj)
		activePool[Label_Default][obj] = struct{}{}
		obj.Obj().IsActive = true
		return true
	}
	return false
}

// Deactivate an object from active list, if it exists in it.
func Deactivate(obj base.IGameObject2D) bool {
	if ContainsActiveDefault(obj) {
		delete(activePool[Label_Default], obj)
		inactivePool[Label_Default][obj] = struct{}{}
		obj.Obj().IsActive = false
		return true
	}
	return false
}
