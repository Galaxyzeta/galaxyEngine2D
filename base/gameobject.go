package base

import (
	"fmt"

	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra/concurrency/lock"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
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
	Hitbox     physics.IShape
	Name       string
	Tags       map[string]struct{}
	Sprite     *graphics.SpriteInstance
	Callbacks  *GameObjectFunctions
	mu         lock.SpinLock
	components map[string]IComponent
	systems    map[string]ISystem
	iobj2d     IGameObject2D
	IsVisible  bool
	IsActive   bool
}

func (obj *GameObject2D) GetIGameObject2D() IGameObject2D {
	return obj.iobj2d
}

// Deprecated: this function should be banned in user mode.
// Using this in a gameloop will cause unexpected result.
func (obj *GameObject2D) SetIGameObject2D(iobj2d IGameObject2D) {
	obj.iobj2d = iobj2d
}

// NewGameObject2D creates a new GameObject2D
func NewGameObject2D(name string) *GameObject2D {
	ret := &GameObject2D{
		Hitbox:     nil,
		Name:       name,
		Tags:       make(map[string]struct{}),
		Sprite:     &graphics.SpriteInstance{},
		Callbacks:  &GameObjectFunctions{},
		IsVisible:  true,
		IsActive:   true,
		mu:         lock.SpinLock{},
		components: map[string]IComponent{},
		systems:    map[string]ISystem{},
	}
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

func (o *GameObject2D) AppendTags(tags ...string) {
	for _, tag := range tags {
		o.Tags[tag] = struct{}{}
	}
}

func (o *GameObject2D) RemoveTags(tags ...string) {
	for _, tag := range tags {
		delete(o.Tags, tag)
	}
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
		panic(fmt.Sprintf("no such component: %v", name))
	}
	return ret
}

func (o *GameObject2D) GetAllComponents() map[string]IComponent {
	return o.components
}

func (o *GameObject2D) GetSubscribedSystemMap() map[string]ISystem {
	return o.systems
}

func (o *GameObject2D) AppendSubscribedSystem(sys ISystem) {
	o.systems[sys.GetName()] = sys
}

func (o *GameObject2D) RemoveSubscribedSystem(sys ISystem) {
	delete(o.systems, sys.GetName())
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
