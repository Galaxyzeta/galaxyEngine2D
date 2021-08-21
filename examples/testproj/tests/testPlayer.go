package tests

import (
	"fmt"
	"time"

	"galaxyzeta.io/engine/core"
	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	objs "galaxyzeta.io/engine/examples/testproj/userspace"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"galaxyzeta.io/engine/sdk"
)

func init() {
	core.GlobalInitializer()
}

func GameEngineTest() {
	sdk.StartApplication(&core.AppConfig{
		Resolution:  &linalg.Vector2f64{X: 640, Y: 480},
		PhysicalFps: 60,
		RenderFps:   60,
		Parallelism: 4,
		Title:       "Test Window",
		InitFunc: func() {
			loadResource()
			csys := system.NewQuadTreeCollision2DSystem(0, physics.NewRectangle(0, 0, 1024, 1024), 4, 128)
			core.RegisterSystem(csys)
			core.RegisterSystem(system.NewPhysics2DSystem(0, csys))
			sdk.Create(objs.TestPlayer_OnCreate)
			var i float64
			var j float64
			for j = 0; j < 5; j++ {
				for i = 0; i < 24; i++ {
					b := sdk.Create(objs.TestBlock_OnCreate)
					tf := b.Obj().GetComponent(component.NameTransform2D).(*component.Transform2D)
					tf.Pos.X = i*16 + (j * 96)
					tf.Pos.Y = 96 * (j + 1)
					this := b.(*objs.TestBlock)
					this.SelfDestructTime = time.Now().Add(time.Millisecond * 200 * time.Duration(int(j*24+i+5)))
				}
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
	graphics.NewFrame("frm_bullet", fmt.Sprintf("%s/examples/testproj/static/megaman/bullet.png", cwd))

	graphics.NewSpriteMeta("spr_megaman", "frm_megaman_01", "frm_megaman_02", "frm_megaman_03")
	graphics.NewSpriteMeta("spr_block", "frm_block")
	graphics.NewSpriteMeta("spr_bullet", "frm_bullet")

}
