package core

import (
	"galaxyzeta.io/engine/graphics"
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
	inputPool    map[Key]struct{}
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

func doDestroy(obj IGameObject2D, isActive *bool) {
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
		inputPool:    make(map[Key]struct{}),
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
