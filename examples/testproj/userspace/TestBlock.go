/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"

	"galaxyzeta.io/engine/base"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/infra/constdef"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/sdk"
)

// TestInputDetection is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.
type TestBlock struct {
	*base.GameObject2D
	tf *component.Transform2D
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestBlock_OnCreate() base.IGameObject2D {
	fmt.Println("SDK Call onCreate")
	gameObject2D := base.NewGameObject2D().
		RegisterRender(__TestBlock_OnRender).
		RegisterStep(constdef.DefaultGameFunction).
		RegisterDestroy(constdef.DefaultGameFunction).
		RegisterComponentIfAbsent(component.NewTransform2D())
	gameObject2D.Sprite = graphics.NewSpriteInstance("spr_block")
	gameObject2D.Sprite.DisableAnimation()
	gameObject2D.Sprite.Z = 10
	return &TestBlock{
		GameObject2D: gameObject2D,
		tf:           gameObject2D.GetComponent(component.NameTransform2D).(*component.Transform2D),
	}
}

func __TestBlock_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestBlock)
	this.Sprite.Render(sdk.GetCamera(0), linalg.Point2f64{X: this.tf.X, Y: this.tf.Y})
}

// GetGameObject2D implements IGameObject2D.
func (t TestBlock) GetGameObject2D() *base.GameObject2D {
	return t.GameObject2D
}
