package core

import (
	"fmt"
	"strings"

	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/level"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

func NewApplicationFromFile(filePath string) *Application {
	levelData := level.ParseGameLevelFile(filePath)
	cwd = GetCwd()

	initializer := func() {
		// load static frames
		staticPath := levelData.LevelMetas.Static
		for _, dir := range levelData.LevelMetas.FrameMetas.Dirs {
			graphics.BatchNewFrames(fmt.Sprintf("%s/%s/%s", cwd, staticPath, dir.Name), func(fileName string) string {
				return fmt.Sprintf("%s%s", dir.Prefix, strings.Split(fileName, ".")[0])
			})
		}
		// register sprites
		for _, spriteMeta := range levelData.LevelMetas.SpriteMetas.Sprites {
			framesCandidates := make([]string, 0)
			for _, frameCandidate := range spriteMeta.Frames {
				framesCandidates = append(framesCandidates, frameCandidate.Name)
			}
			graphics.NewSpriteMeta(spriteMeta.Name, framesCandidates...)
		}
		// build object name-src relation map
		objName2Ctor := map[string]string{}
		for _, objectMeta := range levelData.LevelMetas.ObjectMetas.Objects {
			objName2Ctor[objectMeta.Name] = objectMeta.Name
		}
		// create systems
		// TODO
		csys := system.NewQuadTreeCollision2DSystem(0, physics.NewRectangle(0, 0, 1024, 1024), 4, 64)
		RegisterSystem(csys)
		RegisterSystem(system.NewPhysics2DSystem(1, csys))
		RegisterGfxSystem(system.NewRenderer2DSystem(0))

		// create objects in level
		for _, obj := range levelData.LevelDetails.ObjectDetails {
			ctor, ok := objName2Ctor[obj.Name]
			if !ok {
				panic("failed to find mapping between the object being initialized and object meta.")
			}
			invoker := GetCtor(ctor)
			if !ok {
				panic("failed to find mapping between the object being initialized and constructor map.")
			}

			tf := Create(invoker).Obj().GetComponent(component.NameTransform2D).(*component.Transform2D)
			tf.Pos.X = float64(obj.X)
			tf.Pos.Y = float64(obj.Y)
		}
	}

	// start application
	appCfg := levelData.LevelMetas.ApplicationMetas
	return NewApplication(&AppConfig{
		Resolution:  linalg.NewVector2f64Ptr(appCfg.Resolution.W, appCfg.Resolution.H),
		PhysicalFps: appCfg.FPS.Physics,
		RenderFps:   appCfg.FPS.Render,
		Parallelism: appCfg.Parallelism,
		Title:       appCfg.Title,
		InitFunc:    initializer,
	})
}
