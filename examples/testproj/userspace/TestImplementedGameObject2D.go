/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"

	"galaxyzeta.io/engine/core"
)

// TestImplementedGameObject2D is a golang GameObject2D testing template,
// It illustrates how to use Galaxy2DEngine.
type TestImplementedGameObject2D struct {
	*core.GameObject2D
	counter  int
	counter2 int
}

//TestImplementedGameObject2D_OnCreate is a public constructor.
func TestImplementedGameObject2D_OnCreate() core.IGameObject2D {
	fmt.Println("SDK Call onCreate")
	gameObject2D := core.NewGameObject2D().
		RegisterRender(__TestImplementedGameObject2D_OnRender).
		RegisterStep(__TestImplementedGameObject2D_OnStep).
		RegisterDestroy(__TestImplementedGameObject2D_OnDestroy)
	return &TestImplementedGameObject2D{
		GameObject2D: gameObject2D,
		counter:      0,
		counter2:     0,
	}
}

//__TestImplementedGameObject2D_OnStep is intentionally names with two underlines,
// telling user never call this function in other functions, that will not work,
// even damaging the whole game logic.
func __TestImplementedGameObject2D_OnStep(obj core.IGameObject2D) {
	fmt.Println("onstep cb")
	this := obj.(*TestImplementedGameObject2D)
	this.counter2++
	if this.counter2 == 360 {
		core.Destroy(obj)
	}
}

func __TestImplementedGameObject2D_OnRender(obj core.IGameObject2D) {
	this := obj.(*TestImplementedGameObject2D)
	this.counter++
	if this.counter == 60 {
		this.counter = 0
		fmt.Println("Trigger render")
	}
}

func __TestImplementedGameObject2D_OnDestroy(obj core.IGameObject2D) {
	this := obj.(*TestImplementedGameObject2D)
	fmt.Println("SDK Call onDestory cb")
	fmt.Println(this.counter2)
	core.Exit()
}

// GetGameObject2D implements IGameObject2D.
func (t TestImplementedGameObject2D) GetGameObject2D() *core.GameObject2D {
	return t.GameObject2D
}
