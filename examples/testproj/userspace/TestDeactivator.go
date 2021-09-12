package objs

import (
	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/collision"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
)

const __TestDeactivatorName = "obj_testDeactivator"

func init() {
	core.RegisterCtor(__TestDeactivatorName, TestDeactivator_OnCreate)
}

type TestDeactivator struct {
	*base.GameObject2D
	tf *component.Transform2D

	csys collision.ICollisionSystem
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestDeactivator_OnCreate() base.IGameObject2D {
	this := &TestDeactivator{}

	this.tf = component.NewTransform2D()

	this.GameObject2D = base.NewGameObject2D("deactivator").
		RegisterStep(__TestDeactivator_OnStep).
		RegisterComponentIfAbsent(this.tf)

	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)

	return this
}

func __TestDeactivator_OnStep(iobj base.IGameObject2D) {
	this := iobj.(*TestDeactivator)

	cam := graphics.GetCurrentCamera()
	iobjsMap := core.GetIGameobjectsMap()
	cols := collision.CollidersAtPolygonWithAny(this.csys, cam.GetPolygon(), collision.All)
	camColliderIobjMap := make(map[base.IGameObject2D]struct{})
	for _, col := range cols {
		camColliderIobjMap[col.I()] = struct{}{}
	}

	for testIobj := range iobjsMap {
		if _, ok := camColliderIobjMap[testIobj]; !ok && testIobj != iobj {
			core.Deactivate(testIobj)
		}
	}

	for testIobj := range camColliderIobjMap {
		if _, ok := iobjsMap[testIobj]; !ok && testIobj != iobj {
			core.Activate(testIobj)
		}
	}
}

// GetGameObject2D implements IGameObject2D.
func (t TestDeactivator) Obj() *base.GameObject2D {
	return t.GameObject2D
}
