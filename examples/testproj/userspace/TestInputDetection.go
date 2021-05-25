/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/input"
	keys "galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/sdk"
)

// TestInputDetection is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.
type TestInputDetection struct {
	*core.GameObject2D
	counter                int
	counter2               int
	keyboardCounter        int
	keyboardNotHeldCounter int
	status                 int
	statusCounter          int
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestImplementedGameObject2D_OnCreate() core.IGameObject2D {
	fmt.Println("SDK Call onCreate")
	gameObject2D := core.NewGameObject2D().
		RegisterRender(__TestImplementedGameObject2D_OnRender).
		RegisterStep(__TestImplementedGameObject2D_OnStep).
		RegisterDestroy(__TestImplementedGameObject2D_OnDestroy)
	return &TestInputDetection{
		GameObject2D:           gameObject2D,
		counter:                0,
		counter2:               0,
		keyboardCounter:        0,
		status:                 0,
		statusCounter:          0,
		keyboardNotHeldCounter: 0,
	}
}

//__TestImplementedGameObject2D_OnStep is intentionally names with two underlines,
// telling user never call this function in other functions, that will not work,
// even damaging the whole game logic.
func __TestImplementedGameObject2D_OnStep(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	this.counter2++
	if input.IsKeyPressed(keys.KeyW) {
		fmt.Println("Key W pressed")
		this.status = 1
	}
	if input.IsKeyReleased(keys.KeyW) {
		fmt.Println("Key W released")
		this.status = 0
	}
	if input.IsKeyHeld(keys.KeyW) {
		this.keyboardCounter++
	} else {
		this.keyboardNotHeldCounter++
	}
	if this.status == 1 {
		this.statusCounter++
	}
	if this.counter2 == 360 {
		sdk.Destroy(obj)
	}
}

func __TestImplementedGameObject2D_OnRender(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	this.counter++
	if this.counter == 60 {
		this.counter = 0
		fmt.Println("Trigger render")
	}
}

func __TestImplementedGameObject2D_OnDestroy(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	fmt.Println("SDK Call onDestroy cb")
	fmt.Println("Counter:", this.counter2)
	fmt.Println("KbdCounter:", this.keyboardCounter)
	fmt.Println("KbdNotHeldCounter", this.keyboardNotHeldCounter)
	fmt.Println("StatusCounter", this.statusCounter)

	sdk.Exit()
}

// GetGameObject2D implements IGameObject2D.
func (t TestInputDetection) GetGameObject2D() *core.GameObject2D {
	return t.GameObject2D
}
