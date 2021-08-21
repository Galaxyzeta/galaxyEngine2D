/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"container/list"
	"fmt"
	"time"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
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
	// -- must require
	*base.GameObject2D

	// -- system requirement
	tf     *component.Transform2D
	rb     *component.RigidBody2D
	pc     *component.PolygonCollider
	csys   collision.ICollisionSystem
	logger *logger.Logger

	// -- user defined
	canJump            bool          // whether the user can jump or not
	lastJumpTime       time.Time     // last player jump time
	jumpPreventionTime time.Duration // stop the user from operating a jump in this duration
	speed              float64
	jumpForceElem      *list.Element
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
	this.GameObject2D = base.NewGameObject2D("player").
		RegisterRender(__TestPlayer_OnRender).
		RegisterStep(__TestPlayer_OnStep).
		RegisterDestroy(__TestPlayer_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.rb).
		RegisterComponentIfAbsent(this.pc)
	this.GameObject2D.Sprite = spr

	// Enable gravity
	this.rb.UseGravity = true
	this.rb.SetGravity(270, 0.1)

	this.logger = logger.New("player")
	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)

	core.SubscribeSystem(this, system.NamePhysics2DSystem)
	core.SubscribeSystem(this, system.NameCollision2Dsystem)

	this.jumpPreventionTime = time.Millisecond * 50
	this.lastJumpTime = time.Now()
	this.speed = 2

	return this
}

//__TestPlayer_OnStep is intentionally names with two underlines,
// telling user never call this function in other functions, that will not work,
// even damaging the whole game logic.
func __TestPlayer_OnStep(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	isKeyHeld := false
	var dx float64 = 0
	var dy float64 = 0

	// movement
	if input.IsKeyHeld(keys.KeyA) && !collision.HasColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(-2, 0), "solid") {
		dx = -this.speed
		isKeyHeld = true
	} else if input.IsKeyHeld(keys.KeyD) && !collision.HasColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(2, 0), "solid") {
		dx = this.speed
		isKeyHeld = true
	}

	// functional debug
	if input.IsKeyPressed(keys.KeyE) {
		this.speed += 1
	}
	if input.IsKeyPressed(keys.KeyQ) {
		this.speed -= 1
	}

	// jump
	if input.IsKeyPressed(keys.KeyW) && this.canJump && time.Since(this.lastJumpTime) > this.jumpPreventionTime {
		this.canJump = false
		this.lastJumpTime = time.Now()
		this.jumpForceElem = this.rb.AddForce(component.SpeedVector{
			Acceleration: 0.05,
			Direction:    90,
			Speed:        5,
		})
	}

	// ceil detecting
	if collision.HasColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(0, this.rb.GetVspeed()), "solid") {
		if this.jumpForceElem != nil {
			this.rb.RemoveForce(this.jumpForceElem)
			this.jumpForceElem = nil
		}
	}

	// animation
	if isKeyHeld {
		this.Sprite.EnableAnimation()
	} else {
		this.Sprite.DisableAnimation()
	}

	this.tf.Translate(dx, dy)

	if val := collision.ColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(0, 1), "solid"); val != nil {
		if time.Since(this.lastJumpTime) > this.jumpPreventionTime {
			this.canJump = true
		}
	} else {
		this.canJump = false
	}
}

func __TestPlayer_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	this.Sprite.Render(sdk.GetCamera(), linalg.Point2f64(this.tf.Pos))
	this.csys.(*system.QuadTreeCollision2DSystem).Traverse(false, func(pc *component.PolygonCollider, qn *collision.QTreeNode, at collision.AreaType, idx int) bool {
		graphics.DrawRectangle(qn.GetArea(), linalg.NewRgbaF64(0, 1, 0, 1))
		if pc.GetIGameObject2D().GetGameObject2D().Name == "player" {
			pcRect := pc.Collider.GetBoundingBox().ToRectangle()
			graphics.DrawRectangle(pc.Collider.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(0, 0, 1, 1))
			graphics.DrawRectangle(pcRect.CropOutside(pcRect.Height/2, pcRect.Width/2), linalg.NewRgbaF64(0, 0, 1, 1))
			graphics.DrawRectangle(qn.GetArea().CropOutside(-1, -1), linalg.NewRgbaF64(1, 0, 0, 1))
		}
		ceilTest := this.pc.Collider.Shift(0, -this.rb.GravityVector.Speed*2)
		graphics.DrawRectangle(ceilTest.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(1, 0, 0, 1))

		return false
	})
	// graphics.DrawSegment(linalg.NewSegmentf64(this.tf.X(), this.tf.Y(), this.tf.X()+64, this.tf.Y()), linalg.NewRgbaF64(1, 0, 0, 0))
}

func __TestPlayer_OnDestroy(obj base.IGameObject2D) {
	fmt.Println("SDK Call onDestroy cb")
	sdk.Exit()
}

// GetGameObject2D implements IGameObject2D.
func (t TestPlayer) GetGameObject2D() *base.GameObject2D {
	return t.GameObject2D
}
