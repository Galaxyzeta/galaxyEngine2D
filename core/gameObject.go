package core

import (
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/linalg"
)

type GameObject2DProperties struct {
	Position linalg.Point2f
}

type GameObjectFunctions struct {
	OnCreate  func()
	OnStep    func(self IGameObject2D)
	OnRender  func(self IGameObject2D)
	OnDestroy func(self IGameObject2D)
}

type IGameObject2D interface {
	GetGameObject2D() *GameObject2D
}

type GameObject2D struct {
	prevStats    GameObject2DProperties
	currentStats GameObject2DProperties
	processor    *subLoop // on which subloop is this gameObject2D being processed.
	Hitbox       graphics.IShape
	Sprite       *graphics.Sprite
	Callbacks    *GameObjectFunctions
	inputPool    map[keys.Key]struct{}
	IsVisible    bool
	isActive     bool
}

// doCreate does actual creation.
func doCreate(constructor func() IGameObject2D, isActive *bool) IGameObject2D {
	obj := constructor()
	processor := coreController.roundRobin()
	obj.GetGameObject2D().processor = processor
	processor.registerChannel <- resourceAccessRequest{
		payload:  obj,
		isActive: isActive,
	}
	return obj
}

func DoDestroy(obj IGameObject2D, isActive *bool) {
	obj2d := obj.GetGameObject2D()
	if obj2d.Callbacks.OnDestroy != nil {
		obj2d.Callbacks.OnDestroy(obj)
	}
	obj2d.processor.unregisterChannel <- resourceAccessRequest{
		payload:  obj,
		isActive: isActive,
	}
}

// NewGameObject2D creates a new GameObject2D
func NewGameObject2D() *GameObject2D {
	return &GameObject2D{
		prevStats:    GameObject2DProperties{},
		currentStats: GameObject2DProperties{},
		processor:    &subLoop{},
		Hitbox:       nil,
		Sprite:       &graphics.Sprite{},
		Callbacks:    &GameObjectFunctions{},
		inputPool:    make(map[keys.Key]struct{}),
		IsVisible:    true,
		isActive:     true,
	}
}

func (o *GameObject2D) RegisterStep(method func(IGameObject2D)) *GameObject2D {
	o.Callbacks.OnStep = method
	return o
}

func (o *GameObject2D) RegisterRender(method func(IGameObject2D)) *GameObject2D {
	o.Callbacks.OnRender = method
	return o
}

func (o *GameObject2D) RegisterDestroy(method func(IGameObject2D)) *GameObject2D {
	o.Callbacks.OnDestroy = method
	return o
}

// +------------------------+
// |	  	GameObjs	 	|
// +------------------------+

// Create will instantiate an object immediately.
// The object will be put to the global resource pool in next physical tick.
func Create(constructor func() IGameObject2D) IGameObject2D {
	return doCreate(constructor, infra.BoolPtr_True)
}

// CreateInactive will instantiate an inactive object immediately.
// The object will be put to the global resource pool in next physical tick.
func CreateInactive(constructor func() IGameObject2D) IGameObject2D {
	return doCreate(constructor, infra.BoolPtr_False)
}

// Destroy will deconstruct an active/inactive object immediately.
// The object will be truely removed from resource pool in the next physical tick.
func Destroy(obj IGameObject2D) {
	DoDestroy(obj, nil)
}

// Activate an object from deactive list, if it exists in it.
func Activate(obj IGameObject2D) bool {
	if ContainsInactiveDefault(obj) {
		delete(inactivePool[Label_Default], obj)
		activePool[Label_Default][obj] = struct{}{}
		obj.GetGameObject2D().isActive = true
		return true
	}
	return false
}

// Deactivate an object from active list, if it exists in it.
func Deactivate(obj IGameObject2D) bool {
	if ContainsActiveDefault(obj) {
		delete(activePool[Label_Default], obj)
		inactivePool[Label_Default][obj] = struct{}{}
		obj.GetGameObject2D().isActive = false
		return true
	}
	return false
}
