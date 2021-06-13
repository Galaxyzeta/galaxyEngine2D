package core

import (
	"galaxyzeta.io/engine/component"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra"
	cc "galaxyzeta.io/engine/infra/concurrency"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/linalg"
)

type DefaultComponentWrapper struct {
	Position linalg.Point2f32
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
	processor  *subLoop // on which subloop is this gameObject2D being processed.
	Hitbox     graphics.IShape
	Sprite     *graphics.SpriteInstance
	Callbacks  *GameObjectFunctions
	inputPool  map[keys.Key]struct{}
	mu         cc.SpinLock
	components map[string]IComponent
	iobj2d     IGameObject2D
	IsVisible  bool
	isActive   bool
}

// doCreate does actual creation.
func doCreate(constructor func() IGameObject2D, isActive *bool) IGameObject2D {
	obj := constructor()
	processor := coreController.roundRobin()
	obj.GetGameObject2D().iobj2d = obj
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
	ret := &GameObject2D{
		processor:  &subLoop{},
		Hitbox:     nil,
		Sprite:     &graphics.SpriteInstance{},
		Callbacks:  &GameObjectFunctions{},
		inputPool:  make(map[keys.Key]struct{}),
		IsVisible:  true,
		isActive:   true,
		mu:         cc.SpinLock{},
		components: map[string]IComponent{},
	}
	// register default component
	ret.RegisterComponentIfAbsent(component.NewTransform2D())
	return ret
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

func (o *GameObject2D) RegisterComponent(com IComponent) *GameObject2D {
	o.components[com.GetName()] = com
	return o
}

func (o *GameObject2D) RegisterComponentIfAbsent(com IComponent) *GameObject2D {
	_, ok := o.components[com.GetName()]
	if !ok {
		o.RegisterComponent(com)
	}
	return o
}

func (o *GameObject2D) GetComponent(name string) IComponent {
	ret, ok := o.components[name]
	if !ok {
		panic("no such component")
	}
	return ret
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
	doDestroy(obj, nil)
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

// +------------------------+
// |	 Lock Properties	|
// +------------------------+

func (obj2d *GameObject2D) Lock() {
	obj2d.mu.Lock()
}

func (obj2d *GameObject2D) Unlock() {
	obj2d.mu.Unlock()
}
