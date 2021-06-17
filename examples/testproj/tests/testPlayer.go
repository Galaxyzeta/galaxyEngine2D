package tests

import (
	"fmt"

	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	objs "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/graphics"
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
		Parallelism: 4,
		Title:       "Test Window",
		InitFunc: func() {
			loadResource()
			core.RegisterSystem(system.NewPhysics2DSystem(0))
			sdk.Create(objs.TestPlayer_OnCreate)
			var i float32
			for i = 0; i < 480/16; i++ {
				b := sdk.Create(objs.TestBlock_OnCreate)
				tf := b.GetGameObject2D().GetComponent(component.NameTransform2D).(*component.Transform2D)
				tf.X = i * 16
				tf.Y = 128
			}
		},
	})
}

func loadResource() {
	cwd := core.GetCwd()
	graphics.NewFrame("frm_megaman_01", fmt.Sprintf("%s/examples/testproj/static/megaman/megaman-running-01.png", cwd))
	graphics.NewFrame("frm_megaman_02", fmt.Sprintf("%s/examples/testproj/static/megaman/megaman-running-02.png", cwd))
	graphics.NewFrame("frm_megaman_03", fmt.Sprintf("%s/examples/testproj/static/megaman/megaman-running-03.png", cwd))
	graphics.NewFrame("frm_block", fmt.Sprintf("%s/examples/testproj/static/megaman/block.png", cwd))

	graphics.NewSpriteMeta("spr_megaman", "frm_megaman_01", "frm_megaman_02", "frm_megaman_03")
	graphics.NewSpriteMeta("spr_block", "frm_block")

}
