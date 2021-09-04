package objs

import (
	"time"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra/logger"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"galaxyzeta.io/engine/sdk"
)

const __TestEnemy_Name = "obj_testEnemy"

func init() {
	core.RegisterCtor(__TestEnemy_Name, TestEnemy_OnCreate)
}

type TestEnemy struct {
	// - system requirements
	*base.GameObject2D

	BasicComponentsBundle
	logger *logger.Logger

	// - custom properties
	hp              int
	maxhp           int
	lastHitTime     time.Time
	hitPreventionCD time.Duration
}

func TestEnemy_OnCreate() base.IGameObject2D {
	this := &TestEnemy{}

	animator := graphics.NewAnimator(graphics.StateClipPair{
		State: "idle",
		Clip:  graphics.NewSpriteInstance("spr_miner"),
	})

	this.tf = component.NewTransform2D()
	this.sr = component.NewSpriteRendererWithOptions(animator, this.tf, false, graphics.RenderOptions{
		Pivot: &physics.Pivot{
			Option: physics.PivotOption_BottomCenter,
		},
	})
	this.pc = component.NewPolygonCollider(animator.Spr().GetHitbox(&this.tf.Pos, physics.Pivot{Option: physics.PivotOption_BottomCenter}), this)
	this.logger = logger.New("miner")

	this.GameObject2D = base.NewGameObject2D("enemy").
		RegisterRender(__TestEnemy_OnRender).
		RegisterStep(__TestEnemy_OnStep).
		RegisterDestroy(__TestEnemy_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.pc).
		RegisterComponentIfAbsent(this.sr).
		AppendTags("enemy")

	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)

	core.SubscribeSystem(this, system.NameCollision2Dsystem)
	core.SubscribeSystem(this, system.NameRenderer2DSystem)

	// --- custom properties

	this.maxhp = 10
	this.hp = this.maxhp
	this.hitPreventionCD = time.Second

	return this
}

func __TestEnemy_OnStep(iobj base.IGameObject2D) {
	this := iobj.(*TestEnemy)
	if this.hp <= 0 {
		sdk.Destroy(this)
	}
	if time.Since(this.lastHitTime) < this.hitPreventionCD {
		// hit prevention
		// TODO implement with a phaser effector
		this.sr.Animator.Spr().DisableAnimation()
	} else {
		this.sr.Animator.Spr().EnableAnimation()
	}
}

func __TestEnemy_OnDestroy(iobj base.IGameObject2D) {
	// Your code here ...
}

func __TestEnemy_OnRender(obj base.IGameObject2D) {
	// Your code here ...
	this := obj.(*TestEnemy)

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

}

// GetGameObject2D implements IGameObject2D.
func (t TestEnemy) Obj() *base.GameObject2D {
	return t.GameObject2D
}

// TakeDamage implements IDamageable
func (t *TestEnemy) TakeDamage(dmg int) {
	if time.Since(t.lastHitTime) < t.hitPreventionCD {
		return // will not cause damage, if the enemy is already in hit prevention status.
	}
	t.logger.Debugf("take damage = %d", dmg)
	t.hp -= dmg
	t.lastHitTime = time.Now()
}
