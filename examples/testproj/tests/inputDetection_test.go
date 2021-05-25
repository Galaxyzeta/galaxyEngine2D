package tests

import (
	"galaxyzeta.io/engine/sdk"
	"testing"

	"galaxyzeta.io/engine/core"
	objs "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/linalg"
)

func init() {
	core.GlobalInitializer()
}

func TestGameEngine(t *testing.T) {
	sdk.StartApplication(&core.AppConfig{
		Resolution:  &linalg.Vector2i{X: 640, Y: 320},
		PhysicalFps: 60,
		RenderFps:   60,
		WorkerCount: 4,
		Title:       "Test Window",
		InitFunc: func() {
			sdk.Create(objs.TestImplementedGameObject2D_OnCreate)
		},
	})
}
