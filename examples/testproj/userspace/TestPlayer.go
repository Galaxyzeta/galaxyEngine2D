/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"container/list"
	"fmt"
	"math"
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

const __TestPlayer_Name = "obj_testPlayer"

func init() {
	core.RegisterCtor(__TestPlayer_Name, TestPlayer_OnCreate)
}

// TestPlayer is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.
type TestPlayer struct {
	// -- must require
	*base.GameObject2D

	// -- system requirement
	BasicComponentsBundle
	logger *logger.Logger

	// -- user defined
	canJump            bool          // whether the user can jump or not
	isAir              bool          // whether the user is in air
	lastJumpTime       time.Time     // last player jump time
	jumpPreventionTime time.Duration // stop the user from operating a jump in this duration
	speed              float64
	jumpForceElem      *list.Element

	hp int
}

//TestPlayer_OnCreate is a public constructor.
func TestPlayer_OnCreate() base.IGameObject2D {
	fmt.Println("SDK Call onCreate")

	this := &TestPlayer{}

	animator := graphics.NewAnimator(graphics.StateClipPair{
		State: "run",
		Clip:  graphics.NewSpriteInstance("spr_megaman"),
	})

	this.tf = component.NewTransform2D()
	this.rb = component.NewRigidBody2D()
	this.sr = component.NewSpriteRendererWithOptions(animator, this.tf, false, graphics.RenderOptions{
		Pivot: &physics.Pivot{
			Option: physics.PivotOption_BottomCenter,
		},
	})
	this.pc = component.NewPolygonCollider(animator.Spr().GetHitbox(&this.tf.Pos, physics.Pivot{Option: physics.PivotOption_BottomCenter}), this)
	this.GameObject2D = base.NewGameObject2D("player").
		RegisterRender(__TestPlayer_OnRender).
		RegisterStep(__TestPlayer_OnStep).
		RegisterDestroy(__TestPlayer_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.rb).
		RegisterComponentIfAbsent(this.pc).
		RegisterComponentIfAbsent(this.sr)

	// Enable gravity
	this.rb.UseGravity = true
	this.rb.SetGravity(270, 0.15)

	this.logger = logger.New("player")
	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)

	core.SubscribeSystem(this, system.NamePhysics2DSystem)
	core.SubscribeSystem(this, system.NameCollision2Dsystem)
	core.SubscribeSystem(this, system.NameRenderer2DSystem)

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
	if input.IsKeyHeld(keys.KeyA) && !collision.HasColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(-this.speed*2, 0), "solid") {
		dx = -this.speed
		this.sr.Scale.X = -1
		isKeyHeld = true
	} else if input.IsKeyHeld(keys.KeyD) && !collision.HasColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(this.speed*2, 0), "solid") {
		dx = this.speed
		this.sr.Scale.X = 1
		isKeyHeld = true
	}

	// // mouse
	// if input.IsKeyPressed(keys.MouseButton1) {
	// 	x, y := core.GetCursorPos()
	// 	this.logger.Debugf("cx, cy = %f, %f", x, y)
	// }

	// shoot
	if input.IsKeyPressed(keys.MouseButton1) {
		projectile := sdk.Create(TestProjectile_OnCreate).(*TestProjectile)
		projectile.selfDestruct = time.Now().Add(time.Second * 5)
		projectile.owner = this
		cx, cy := core.GetCursorPos()
		ox := this.tf.X()
		oy := this.tf.Y() - 16
		projectile.directionRad = math.Atan2(cy-oy, cx-ox)
		projectile.speed = 5
		projectile.tf.Pos = linalg.NewVector2f64(ox, oy)

		this.logger.Debugf("create bullet, mouse cursor is at %f, %f", cx, cy)
	}

	// change speed
	if input.IsKeyPressed(keys.KeyE) {
		this.speed += 1
	}
	if input.IsKeyPressed(keys.KeyQ) {
		this.speed -= 1
	}

	// jump
	if input.IsKeyPressed(keys.KeyW) && this.canJump {
		this.canJump = false
		this.lastJumpTime = time.Now()
		this.jumpForceElem = this.rb.AddForce(component.SpeedVector{
			Acceleration: 1,
			Direction:    90,
			Speed:        12,
		})
	}

	// animation
	if isKeyHeld {
		this.sr.Spr().EnableAnimation()
	} else {
		this.sr.Spr().DisableAnimation()
	}

	vspeed := this.rb.GetVspeed()
	if val := collision.ColliderAtPolygonWithTag(this.csys, this.pc.Collider.Shift(0, vspeed), "solid"); val != nil {
		// something is above the player / beneath the player
		if this.isAir {
			if vspeed < 0 {
				// player ascending, detect ceil
				this.logger.Debug("ceil")
				this.rb.RemoveForce(this.jumpForceElem)
				this.jumpForceElem = nil

			} else if vspeed > 0 {
				// player descending
				if time.Since(this.lastJumpTime) > this.jumpPreventionTime {
					// already jumped into the air, and in next frame he will go into ground
					// in order to prevent this, snap him to the ground
					this.logger.Debug("ground")

					this.rb.RemoveForce(this.jumpForceElem)
					this.jumpForceElem = nil

					this.tf.Pos.Y += (val.Collider.GetAnchor().Y - this.pc.Collider.GetAnchor().Y) // snap to the ground
					this.isAir = false
					this.canJump = true
				}
				// player has just jumped, but we don't want to pull him back immediately.
			}
		}
		// something above/beneath but the player is grounded, in this case do nothing
	} else {
		// if vspeed has variation, then:
		// 1. player is jumping;
		// 2. floor collapsed, triggers isAir = true, thus increased gravity speed, causing vspeed to change
		if vspeed != 0 {
			this.isAir = true
			this.canJump = false
		}
	}

	this.tf.Translate(dx, dy)
}

func __TestPlayer_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	this.csys.(*system.QuadTreeCollision2DSystem).Traverse(false, func(pc *component.PolygonCollider, qn *collision.QTreeNode, at collision.AreaType, idx int) bool {
		graphics.DrawRectangle(qn.GetArea(), linalg.NewRgbaF64(0, 1, 0, 1))
		if pc.I().Obj().Name == "player" {
			graphics.DrawRectangle(pc.Collider.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(0, 0, 1, 1))
			graphics.DrawRectangle(qn.GetArea().CropOutside(-1, -1), linalg.NewRgbaF64(1, 0, 0, 1))
		}
		return false
	})

	// mark sprite anchor center
	srx := this.sr.GetHitbox().GetAnchor().X
	sry := this.sr.GetHitbox().GetAnchor().Y
	graphics.DrawSegment(linalg.NewSegmentf64(srx-4, sry, srx+4, sry), linalg.NewRgbaF64(0, 1, 0, 1))
	graphics.DrawSegment(linalg.NewSegmentf64(srx, sry-4, srx, sry+4), linalg.NewRgbaF64(0, 1, 0, 1))

	graphics.DrawRectangle(this.pc.Collider.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(0, 1, 0, 1))

	// mark tf anchor center
	tfx := this.tf.Pos.X
	tfy := this.tf.Pos.Y
	graphics.DrawSegment(linalg.NewSegmentf64(tfx-4, tfy, tfx+4, tfy), linalg.NewRgbaF64(1, 0, 0, 1))
	graphics.DrawSegment(linalg.NewSegmentf64(tfx, tfy-4, tfx, tfy+4), linalg.NewRgbaF64(1, 0, 0, 1))

	// if this.canJump {
	// 	graphics.DrawSegment(linalg.NewSegmentf64(this.tf.X(), this.tf.Y(), this.tf.X()+32, this.tf.Y()), linalg.NewRgbaF64(0, 1, 0, 0))
	// } else {
	// 	graphics.DrawSegment(linalg.NewSegmentf64(this.tf.X(), this.tf.Y(), this.tf.X()+32, this.tf.Y()), linalg.NewRgbaF64(1, 0, 0, 0))
	// }

	// if this.isAir {
	// 	graphics.DrawSegment(linalg.NewSegmentf64(this.tf.X()+32, this.tf.Y(), this.tf.X()+64, this.tf.Y()), linalg.NewRgbaF64(0, 1, 0, 0))
	// } else {
	// 	graphics.DrawSegment(linalg.NewSegmentf64(this.tf.X()+32, this.tf.Y(), this.tf.X()+64, this.tf.Y()), linalg.NewRgbaF64(1, 0, 0, 0))
	// }
}

func __TestPlayer_OnDestroy(obj base.IGameObject2D) {
	fmt.Println("SDK Call onDestroy cb")
	sdk.Exit()
}

// GetGameObject2D implements IGameObject2D.
func (t TestPlayer) Obj() *base.GameObject2D {
	return t.GameObject2D
}

// TakeDamage implements rpgbase.TakeDamage
func (t *TestPlayer) TakeDamage(dmg int) {
	t.logger.Debugf("take damage = %v", dmg)
	t.hp -= dmg
}
