# Galaxy Engine 2D

使用 OpenGL 渲染 API 的 2D 游戏开发框架。正在开发。

- 多协程核心循环高效执行你的每一物理帧更新。
- 基于生命周期体系构建你的游戏，逻辑简单清晰。

设计文档：devlog.md

## How 2 Use

一个简单的演示：

```go
package core

import (
	"fmt"
	"testing"
	"time"
)

//一个测试用的游戏对象。
type TestImplementedGameObject2D struct {
	*GameObject2D	// 必须包含一个通用游戏对象父物体。
	counter  int
	counter2 int
}

//构造函数，调用Create后立即执行，物体将会在下一帧更新开始前被加入游戏。
//注意：回调函数不能直接调用，否则无法实现资源管理，必须调写好的 SDK 才能发挥效果。
func OnCreate() IGameObject2D {
	gameObject2D := NewGameObject2D().
		RegisterRender(OnStep).
		RegisterStep(OnRender).
		RegisterDestroy(OnDestroy)
	return &TestImplementedGameObject2D{
		GameObject2D: gameObject2D,
		counter:      0,
		counter2:     0,
	}
}

//更新函数，引擎以物理FPS为频率执行这个函数。
func OnStep(obj IGameObject2D) {
	this := obj.(*TestImplementedGameObject2D)
	this.counter2++
	if this.counter2 == 60 {
		Destroy(obj)
	}
}

//渲染函数，引擎以渲染FPS为频率执行这个函数。
func OnRender(obj IGameObject2D) {
	this := obj.(*TestImplementedGameObject2D)
	this.counter++
	fmt.Println("Trigger render")
}

//析构函数，在OnStep中调用Destroy后立即执行，游戏物体将在本次物理更新结束前被移出游戏。
func OnDestroy(obj IGameObject2D) {
	this := obj.(*TestImplementedGameObject2D)
	fmt.Println("onDestory")
	fmt.Println(this.counter2)
}

// 必须实现一个获取该物体 GameObject2D 的方法，以便引擎调用。
func (t TestImplementedGameObject2D) GetGameObject2D() *GameObject2D {
	return t.GameObject2D
}

// 以下内容可以写在 main 函数中。
func TestGameEngine(t *testing.T) {
	ctrl := NewMasterLoop(60, 60, 4)	// 创建核心循环，物理FPS和渲染FPS=60，4个子协程执行更新操作。
	Create(OnCreate)	// 发出物体创建请求。
	ctrl.RunNoBlocking()	// 开始执行核心循环。
	go func() {
		time.Sleep(time.Second * 3)
		ctrl.Kill()		// 3 秒后杀死核心循环。
	}()
	ctrl.Wait()	// 等待核心循环结束。
	fmt.Println("Abort")
}

```