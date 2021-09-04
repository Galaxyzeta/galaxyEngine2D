# Galaxy Engine 2D

-  一个从 0 开始的 2D 游戏引擎，基于 OpenGL 渲染。

## 目标和非目标

- 打造基于 ECS 架构的，易于使用的 Golang 2D 游戏引擎。
- 暂不覆盖 3D 游戏开发。

## 代码Preview

启动游戏引擎：
```go
func GameEngineTest() {
	sdk.StartApplication(&core.AppConfig{
		Resolution:  &linalg.Vector2f64{X: 640, Y: 480},	// 屏幕大小
		PhysicalFps: 60,	// 物理帧率
		RenderFps:   60,	// 渲染帧率
		Parallelism: 4,		// 工作线程个数
		Title:       "Test Window",		// 屏幕标题
		InitFunc: func() {
			// 这里一般是创建场景前的初始化，例如载入资源，创建对象等。
		},
	})
}
```

游戏对象的通用模板：
```go
package objs

import (
	... ...
)

type TestPlayer struct {
	// -- 必须具备的一个组件
	*base.GameObject2D

	// -- 可选，但一般会带上的组件
	tf   *component.Transform2D			// 控制物体的移动
	rb   *component.RigidBody2D			// 物理相关
	pc   *component.PolygonCollider		// 碰撞组件
	sr   *component.SpriteRenderer		// 渲染组件
	csys collision.ICollisionSystem		// 碰撞系统
	logger *logger.Logger				// 异步日志

	// -- 用户自定义属性
	... ...
}

// Create 构造函数将在调用对象创建函数 sdk.Create() 后立即执行，对象将在下一物理帧开始前被加入框架管理
func TestPlayer_OnCreate() base.IGameObject2D {
	this := &TestPlayer{}	// this
	animator := graphics.NewAnimator(graphics.StateClipPair{
		State: "run",
		Clip:  graphics.NewSpriteInstance("spr_megaman"),
	})	// 动画播片机

	this.tf = component.NewTransform2D()
	this.rb = component.NewRigidBody2D()
	this.sr = component.NewSpriteRendererWithOptions(animator, this.tf, false, graphics.RenderOptions{
		Pivot: &physics.Pivot{
			Option: physics.PivotOption_BottomCenter,
		},
	})
	this.pc = component.NewPolygonCollider(animator.Spr().GetHitbox(&this.tf.Pos, physics.Pivot{Option: physics.PivotOption_BottomCenter}), this)
	
	// 1. 注册生命周期事件
	// 2. 将组件注册到游戏对象上，完成E-C绑定
	this.GameObject2D = base.NewGameObject2D("player").
		RegisterRender(__TestPlayer_OnRender).
		RegisterStep(__TestPlayer_OnStep).
		RegisterDestroy(__TestPlayer_OnDestroy).
		RegisterComponentIfAbsent(this.tf).
		RegisterComponentIfAbsent(this.rb).
		RegisterComponentIfAbsent(this.pc).
		RegisterComponentIfAbsent(this.sr)

	// 开启重力
	this.rb.UseGravity = true
	this.rb.SetGravity(270, 0.15)
	
	// 异步日志
	this.logger = logger.New("player")

	// 获取碰撞处理系统
	this.csys = core.GetSystem(system.NameCollision2Dsystem).(collision.ICollisionSystem)

	// 将组件注册到系统中，完成C-S绑定
	core.SubscribeSystem(this, system.NamePhysics2DSystem)
	core.SubscribeSystem(this, system.NameCollision2Dsystem)
	core.SubscribeSystem(this, system.NameRenderer2DSystem)

	// 为2D跳台初始化一些属性...
	this.jumpPreventionTime = time.Millisecond * 50
	this.lastJumpTime = time.Now()
	this.speed = 2

	return this
}

// Step 函数将在每一物理帧执行
// 主要逻辑都在这里，例如控制玩家移动，跳跃，攻击等。
func __TestPlayer_OnStep(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	// Your code here...
}

// Render 函数将在每一渲染帧执行
// 你不需要手动绘制玩家，因为绘制玩家的操作已经交给 SpriteRenderer 组件了
// 这里一般绘制一些特效，绘制碰撞检测框用于Debug
func __TestPlayer_OnRender(obj base.IGameObject2D) {
	this := obj.(*TestPlayer)
	// Your code here...
}

// Destroy 函数将在物体即将被摧毁时( sdk.Destroy() )调用，物体将会在当前物理帧结束后被移除
func __TestPlayer_OnDestroy(obj base.IGameObject2D) {
	// Your code here...
}

```

## RoadMap

- 低阶渲染器
  - Sprite 绘制
  - Shader 插件
  - 动画状态机
- 基础设施
  - 线程池、回环屏障
  - ID 生成器
  - 异步日志
  - 测试断言
- 物理
  - 基本形状的碰撞检测
  - 力计算
  - 碰撞检测优化
- 资源管理
  - LRU 动态加载、卸载资源
  - VBO 管理
- 前端编辑器
  - 关卡 Xml 文件解析
  - 基于 Electron 的编辑界面
- ECS架构
  - Transform
  - RigidBody
  - PolygonCollider
  - SpriteRenderer
- 碰撞检测及其优化
- 用户插件
  - RPG基础组件
  - 窗体控件
  - 伤害计算框架
- RPC网络通信
  - 状态同步组件
  - IO多路复用网络框架
  - 自定义RPC协议

## 