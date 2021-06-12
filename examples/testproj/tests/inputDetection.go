package tests

import (
	"galaxyzeta.io/engine/core"
	objs "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/sdk"
)

func init() {
	core.GlobalInitializer()
}

func GameEngineTest() {
	sdk.StartApplication(&core.AppConfig{
		Resolution:  &linalg.Vector2f32{X: 640, Y: 480},
		PhysicalFps: 60,
		RenderFps:   60,
		WorkerCount: 4,
		Title:       "Test Window",
		InitFunc: func() {
			sdk.Create(objs.TestImplementedGameObject2D_OnCreate)
		},
	})
}
