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
	"galaxyzeta.io/engine/infra/logger"
	"galaxyzeta.io/engine/input"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"galaxyzeta.io/engine/sdk"
)

// TestPlayer is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.
type TestPlayer struct {
	*base.GameObject2D
	tf     *component.Transform2D
	rb     *component.RigidBody2D
	pc     *component.PolygonCollider
	logger *logger.Logger
}

//TestPlayer_OnCreate is a public constructor.
func TestPlayer_OnCreate() base.IGameObject2D {
	fmt.Println("SDK Call onCreate")

	this := &TestPlayer{}

	spr := graphics.NewSpriteInstance("spr_megaman")
	this.tf = component.NewTransform2D()
	this.rb = component.NewRigidBody2D()
	this.pc = component.NewPolygonCollider(spr.GetHitbox(&this.tf.Pos, physics.Pivot{
		Option: physics.PivotOption_TopLeft,
	}), this)
	this.GameObject2D = base.NewGameObject2D().
		RegisterRender(__TestPlayer_OnRender).
		RegisterStep(__TestPlayer_OnStep).
		RegisterDestroy(__TestPlayer_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.rb).
		RegisterComponentIfAbsent(this.pc)
	this.GameObject2D.Sprite = spr

	// Enable gravity
	this.rb.UseGravity = true
	this.rb.SetGravity(270, 0.001)

	this.logger = logger.New("Player")

	core.SubscribeSystem(this, system.NamePhysics2DSystem)
	core.SubscribeSystem(this, system.NameCollision2Dsystem)

	return this
}

//__TestPlayer_OnStep is intentionally names with two underlines,
// telling user never call this function in other functions, that will not work,
// even damaging the whole game logic.
func __TestPlayer_OnStep(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	isKeyHeld := false
	if input.IsKeyHeld(keys.KeyW) {
		this.tf.Translate(0, -1)
		isKeyHeld = true
	} else if input.IsKeyHeld(keys.KeyS) {
		this.tf.Translate(0, 1)
		isKeyHeld = true
	}
	if input.IsKeyHeld(keys.KeyA) {
		this.tf.Translate(-1, 0)
		isKeyHeld = true
	} else if input.IsKeyHeld(keys.KeyD) {
		this.tf.Translate(1, 0)
		isKeyHeld = true
	}
	if isKeyHeld {
		this.Sprite.EnableAnimation()
	} else {
		this.Sprite.DisableAnimation()
	}
}

func __TestPlayer_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	this.Sprite.Render(sdk.GetCamera(), linalg.Point2f64(this.tf.Pos))
	// this.Sprite.RenderWire(sdk.GetCamera(0), linalg.Point2f64(this.tf.Pos), linalg.NewRgbaF64(1, 0, 0, 1))
	// graphics.DrawRectangle(this.pc.Collider.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(1, 0, 0, 1))
	graphics.DrawRectangle(this.pc.Collider.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(1, 0, 0, 1))

}

func __TestPlayer_OnDestroy(obj base.IGameObject2D) {
	fmt.Println("SDK Call onDestroy cb")
	sdk.Exit()
}

// GetGameObject2D implements IGameObject2D.
func (t TestPlayer) GetGameObject2D() *base.GameObject2D {
	return t.GameObject2D
}
