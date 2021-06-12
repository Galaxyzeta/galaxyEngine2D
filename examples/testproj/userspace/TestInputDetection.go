/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"
	"time"

	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/sdk"
)

// TestInputDetection is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.
type TestInputDetection struct {
	*core.GameObject2D
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestImplementedGameObject2D_OnCreate() core.IGameObject2D {
	fmt.Println("SDK Call onCreate")
	gameObject2D := core.NewGameObject2D().
		RegisterRender(__TestImplementedGameObject2D_OnRender).
		RegisterStep(__TestImplementedGameObject2D_OnStep).
		RegisterDestroy(__TestImplementedGameObject2D_OnDestroy)
	gameObject2D.Sprite = graphics.NewSprite(0, 0, time.Millisecond*250,
		fmt.Sprintf("%s/examples/testproj/static/megaman/megaman-running-01.png", core.GetCwd()),
		fmt.Sprintf("%s/examples/testproj/static/megaman/megaman-running-02.png", core.GetCwd()),
		fmt.Sprintf("%s/examples/testproj/static/megaman/megaman-running-03.png", core.GetCwd()))
	return &TestInputDetection{
		GameObject2D: gameObject2D,
	}
}

//__TestImplementedGameObject2D_OnStep is intentionally names with two underlines,
// telling user never call this function in other functions, that will not work,
// even damaging the whole game logic.
func __TestImplementedGameObject2D_OnStep(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	isKeyHeld := false
	if input.IsKeyHeld(keys.KeyW) {
		this.CurrentStats.Position.Y -= 1
		isKeyHeld = true
	} else if input.IsKeyHeld(keys.KeyS) {
		this.CurrentStats.Position.Y += 1
		isKeyHeld = true
	}
	if input.IsKeyHeld(keys.KeyA) {
		this.CurrentStats.Position.X -= 1
		isKeyHeld = true
	} else if input.IsKeyHeld(keys.KeyD) {
		this.CurrentStats.Position.X += 1
		isKeyHeld = true
	}
	if isKeyHeld {
		this.Sprite.EnableAnimation()
	} else {
		this.Sprite.DisableAnimation()
	}
}

func __TestImplementedGameObject2D_OnRender(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	this.Sprite.Render(sdk.GetCamera(0), this.CurrentStats.Position)
}

func __TestImplementedGameObject2D_OnDestroy(obj core.IGameObject2D) {
	// this := obj.(*TestInputDetection)
	fmt.Println("SDK Call onDestroy cb")

	sdk.Exit()
}

// GetGameObject2D implements IGameObject2D.
func (t TestInputDetection) GetGameObject2D() *core.GameObject2D {
	return t.GameObject2D
}
