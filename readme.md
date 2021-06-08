# Galaxy Engine 2D

使用 OpenGL 渲染 API 的 2D 游戏开发框架。正在开发。

- 多协程核心循环高效执行你的每一物理帧更新。
- 基于生命周期体系构建你的游戏，逻辑简单清晰。


设计文档：devlog.md

## How 2 Use

一个简单的演示：

下载项目，我内置了一个测试案例，直接运行 Main 即可：

```go
func GameEngineTest() {
	sdk.StartApplication(&core.AppConfig{
		Resolution:  &linalg.Vector2f32{X: 640, Y: 320},  // 窗口大小 
		PhysicalFps: 60,	// 渲染更新频率
		RenderFps:   60,	// 物理更新频率
		WorkerCount: 1,		// 工作协程数，如果不出现掉帧，开 1 即可。
		Title:       "Test Window",
		InitFunc: func() {  // 初始化函数
			sdk.Create(objs.TestImplementedGameObject2D_OnCreate)
		},
	})
}
```

下面是对游戏对象的定义：

```go
/*
All user defined 2D objects should be put here.
*/

package objs

import (
	"fmt"
	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input"
	keys "galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/sdk"
)

// 测试用游戏对象结构体，必须包含 GameObject2D。
type TestInputDetection struct {
	*core.GameObject2D
	counter                int
	counter2               int
	keyboardCounter        int
	keyboardNotHeldCounter int
	status                 int
	statusCounter          int
}

//构造函数，当调用sdk.Create的时候，框架自动帮我们生成一个有生命周期的对象。
func TestImplementedGameObject2D_OnCreate() core.IGameObject2D {
	fmt.Println("SDK Call onCreate")
	gameObject2D := core.NewGameObject2D().
		RegisterRender(__TestImplementedGameObject2D_OnRender).
		RegisterStep(__TestImplementedGameObject2D_OnStep).
		RegisterDestroy(__TestImplementedGameObject2D_OnDestroy)
	gameObject2D.Sprite = graphics.NewSprite(fmt.Sprintf("%s/examples/testproj/static/Mudkip.png", core.GetCwd()), false, 0, 0)
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

//更新函数，作为测试，这里检测了 W 按键是否被按下。在 360 步之后（FPS=60），销毁这个游戏对象。
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

// 渲染函数。目前版本这里它还不会被调用，TODO。
func __TestImplementedGameObject2D_OnRender(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	this.counter++
	if this.counter == 60 {
		this.counter = 0
		fmt.Println("Trigger render")
	}
}

// 析构函数。当游戏对象即将被销毁时，将会执行此函数。
func __TestImplementedGameObject2D_OnDestroy(obj core.IGameObject2D) {
	this := obj.(*TestInputDetection)
	fmt.Println("SDK Call onDestroy cb")
	fmt.Println("Counter:", this.counter2)
	fmt.Println("KbdCounter:", this.keyboardCounter)
	fmt.Println("KbdNotHeldCounter", this.keyboardNotHeldCounter)
	fmt.Println("StatusCounter", this.statusCounter)

	sdk.Exit()
}

// 必须实现此方法，该游戏对象才能被框架使用。
func (t TestInputDetection) GetGameObject2D() *core.GameObject2D {
	return t.GameObject2D
}

```