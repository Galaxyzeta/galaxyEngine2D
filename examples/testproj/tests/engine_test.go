package test

import (
	"testing"

	"galaxyzeta.io/engine/core"
	objs "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/linalg"
)

func init() {
	core.GlobalInitializer()
}

func TestGameEngine(t *testing.T) {
	core.StartApplication(&core.AppConfig{
		Resolution:  &linalg.Vector2i{X: 640, Y: 320},
		PhysicalFps: 60,
		RenderFps:   60,
		WorkerCount: 4,
		Title:       "Test Window",
		InitFunc: func() {
			core.Create(objs.TestImplementedGameObject2D_OnCreate)
		},
	})
}
