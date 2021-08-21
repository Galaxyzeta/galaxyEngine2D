/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"
	"time"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"galaxyzeta.io/engine/sdk"
)

// TestInputDetection is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.

type TestBlock struct {
	*base.GameObject2D
	tf               *component.Transform2D
	pc               *component.PolygonCollider
	SelfDestructTime time.Time
	csys             collision.ICollisionSystem
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestBlock_OnCreate() base.IGameObject2D {
	fmt.Println("SDK Call onCreate")
	this := &TestBlock{}

	spr := graphics.NewSpriteInstance("spr_block")
	this.tf = component.NewTransform2D()
	this.pc = component.NewPolygonCollider(spr.GetHitbox(&this.tf.Pos, physics.Pivot{
		Option: physics.PivotOption_TopLeft,
	}), this)

	gameObject2D := base.NewGameObject2D("block").
		RegisterRender(__TestBlock_OnRender).
		RegisterStep(__TestBlock_OnStep).
		RegisterDestroy(__TestBlock_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.pc)
	gameObject2D.Sprite = spr
	gameObject2D.Sprite.DisableAnimation()
	gameObject2D.Sprite.Z = 10

	gameObject2D.AppendTags("solid")

	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)
	this.GameObject2D = gameObject2D

	core.SubscribeSystem(this, system.NameCollision2Dsystem)

	return this

}

func __TestBlock_OnStep(iobj base.IGameObject2D) {}

func __TestBlock_OnDestroy(iobj base.IGameObject2D) {
	this := iobj.(*TestBlock)
	this.csys.Unregister(iobj)
}

func __TestBlock_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestBlock)
	this.Sprite.Render(sdk.GetCamera(), linalg.Point2f64(this.tf.Pos))
}

// GetGameObject2D implements IGameObject2D.
func (t TestBlock) Obj() *base.GameObject2D {
	return t.GameObject2D
}
