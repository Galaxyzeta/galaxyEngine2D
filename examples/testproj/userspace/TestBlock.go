/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra/constdef"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"galaxyzeta.io/engine/sdk"
)

// TestInputDetection is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.

type TestBlock struct {
	*base.GameObject2D
	tf *component.Transform2D
	pc *component.PolygonCollider
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

	gameObject2D := base.NewGameObject2D().
		RegisterRender(__TestBlock_OnRender).
		RegisterStep(constdef.DefaultGameFunction).
		RegisterDestroy(constdef.DefaultGameFunction).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.pc)
	gameObject2D.Sprite = spr
	gameObject2D.Sprite.DisableAnimation()
	gameObject2D.Sprite.Z = 10

	this.GameObject2D = gameObject2D

	core.SubscribeSystem(this, system.NameCollision2Dsystem)

	return this

}

func __TestBlock_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestBlock)
	this.Sprite.Render(sdk.GetCamera(), linalg.Point2f64(this.tf.Pos))

	this.Sprite.RenderWire(sdk.GetCamera(), linalg.Point2f64(this.tf.Pos), linalg.NewRgbaF64(1, 0, 0, 1))
}

// GetGameObject2D implements IGameObject2D.
func (t TestBlock) GetGameObject2D() *base.GameObject2D {
	return t.GameObject2D
}
