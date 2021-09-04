/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"time"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
)

const __TestBlock_Name = "obj_testBlock"

func init() {
	core.RegisterCtor(__TestBlock_Name, TestBlock_OnCreate)
}

type TestBlock struct {
	*base.GameObject2D
	tf               *component.Transform2D
	pc               *component.PolygonCollider
	sr               *component.SpriteRenderer
	SelfDestructTime time.Time
	csys             collision.ICollisionSystem
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestBlock_OnCreate() base.IGameObject2D {
	this := &TestBlock{}

	animator := graphics.NewAnimator(graphics.StateClipPair{
		State: "idle",
		Clip:  graphics.NewSpriteInstance("spr_block"),
	})

	this.tf = component.NewTransform2D()
	this.sr = component.NewSpriteRenderer(animator, this.tf, false)
	this.pc = component.NewPolygonCollider(this.sr.GetHitbox(), this)
	this.sr.Spr().DisableAnimation()

	gameObject2D := base.NewGameObject2D("block").
		RegisterRender(__TestBlock_OnRender).
		RegisterStep(__TestBlock_OnStep).
		RegisterDestroy(__TestBlock_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.pc).
		RegisterComponentIfAbsent(this.sr)

	gameObject2D.AppendTags("solid")

	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)
	this.GameObject2D = gameObject2D

	core.SubscribeSystem(this, system.NameCollision2Dsystem)
	core.SubscribeSystem(this, system.NameRenderer2DSystem)

	return this

}

func __TestBlock_OnStep(iobj base.IGameObject2D) {}

func __TestBlock_OnDestroy(iobj base.IGameObject2D) {
	// this := iobj.(*TestBlock)
	// this.csys.Unregister(iobj)
}

func __TestBlock_OnRender(obj base.IGameObject2D) {
}

// GetGameObject2D implements IGameObject2D.
func (t TestBlock) Obj() *base.GameObject2D {
	return t.GameObject2D
}
