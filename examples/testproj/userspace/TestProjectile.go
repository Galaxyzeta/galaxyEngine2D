package objs

import (
	"math"
	"time"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/essentials/rpg/rpgbase"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra/logger"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/sdk"
)

const __TestProjectile_Name = "obj_testProjectile"

func init() {
	core.RegisterCtor(__TestProjectile_Name, TestProjectile_OnCreate)
}

type TestProjectile struct {
	// ---- requirements -----

	*base.GameObject2D

	// ---- components ----
	BasicComponentsBundle
	csys collision.ICollisionSystem

	// ---- custom properties -----
	dmg                  int
	selfDestructDuration time.Duration
	createdAt            time.Time
	owner                base.IGameObject2D
	speed                float64
	directionRad         float64
}

// GetGameObject2D implements IGameObject2D.
func (t TestProjectile) Obj() *base.GameObject2D {
	return t.GameObject2D
}

func (t *TestProjectile) SetOwner(owner base.IGameObject2D) {
	t.owner = owner
}

func TestProjectile_OnCreate() base.IGameObject2D {
	this := &TestProjectile{}

	animator := graphics.NewAnimator(graphics.StateClipPair{
		State: "idle",
		Clip:  graphics.NewSpriteInstance("spr_bullet"),
	})

	this.tf = component.NewTransform2D()
	this.rb = component.NewRigidBody2D()
	this.sr = component.NewSpriteRenderer(animator, this.tf, false)
	this.pc = component.NewPolygonCollider(this.sr.GetHitbox(), this)
	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)

	this.GameObject2D = base.NewGameObject2D("projectile").
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.rb).
		RegisterComponentIfAbsent(this.pc).
		RegisterComponentIfAbsent(this.sr).
		RegisterStep(__TestProjectile_OnStep).
		RegisterDestroy(__TestProjectile_OnDestroy).
		RegisterRender(__TestProjectile_OnRender)

	core.SubscribeSystem(this, system.NameRenderer2DSystem)
	core.SubscribeSystem(this, system.NameCollision2Dsystem)

	// - custom properties
	this.dmg = 5

	return this
}

func __TestProjectile_OnStep(obj base.IGameObject2D) {
	// Your code here ...
	this := obj.(*TestProjectile)

	if time.Since(this.createdAt) >= this.selfDestructDuration {
		sdk.Destroy(this)
	}

	val := collision.ColliderAtPolygonWithAny(this.csys, this.pc.Collider)
	if val != nil {

		for tag := range val.I().Obj().Tags {
			if tag == "solid" {
				// destroy it
				sdk.Destroy(obj)
				return
			}
		}

		// if it is a IDamageable, do damage to it and self-destruct
		if _, ok := val.I().(rpgbase.IDamageable); ok {
			if this.owner != val.I() {
				idmg := val.I().(rpgbase.IDamageable)
				idmg.TakeDamage(this.dmg)
				sdk.Destroy(obj)
				return
			}
		}
	}

	// move bullet
	this.tf.Pos.X += this.speed * math.Cos(this.directionRad)
	this.tf.Pos.Y += this.speed * math.Sin(this.directionRad)

}

func __TestProjectile_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestProjectile)
	graphics.DrawRectangle(this.pc.Collider.GetBoundingBox().ToRectangle(), linalg.NewRgbaF64(1, 0, 0, 1))
}

func __TestProjectile_OnDestroy(obj base.IGameObject2D) {
	// Your code here ...
	a := obj.(*TestProjectile)
	logger.GlobalLogger.Infof("%v", a)
}
