package core

import (
	"fmt"
	"strings"

	"galaxyzeta.io/engine/ecs/component"
	"galaxyzeta.io/engine/ecs/system"
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/parser"
	"galaxyzeta.io/engine/physics"
)

var worldMeta *parser.LevelConfig
var objName2Ctor map[string]string = make(map[string]string)

// NewApplicationFromFile creates a new application from given level definition XML file.
// Not concurrently safe, no need to create multiple applications at same time.
func NewApplicationFromFile(filePath string) *Application {
	worldMeta = parser.ParseGameLevelFile(filePath)
	cwd = GetCwd()

	initializer := func() {
		// load static frames
		staticPath := worldMeta.LevelMetas.Static
		for _, dir := range worldMeta.LevelMetas.FrameMetas.Dirs {
			graphics.BatchNewFrames(fmt.Sprintf("%s/%s/%s", cwd, staticPath, dir.Name), func(fileName string) string {
				return fmt.Sprintf("%s%s", dir.Prefix, strings.Split(fileName, ".")[0])
			})
		}
		// register sprites
		for _, spriteMeta := range worldMeta.LevelMetas.SpriteMetas.Sprites {
			framesCandidates := make([]string, 0)
			for _, frameCandidate := range spriteMeta.Frames {
				framesCandidates = append(framesCandidates, frameCandidate.Name)
			}
			graphics.NewSpriteMeta(spriteMeta.Name, framesCandidates...)
		}
		// build object name-src relation map
		for _, objectMeta := range worldMeta.LevelMetas.ObjectMetas.Objects {
			objName2Ctor[objectMeta.Name] = objectMeta.Name
		}
		// create systems
		// TODO only load user defined systems
		csys := system.NewQuadTreeCollision2DSystem(0, physics.NewRectangle(0, 0, 1024, 1024), 4, 64)
		RegisterSystem(csys)
		RegisterSystem(system.NewPhysics2DSystem(1, csys))
		RegisterGfxSystem(system.NewRenderer2DSystem(0))
		// create cameras
		graphics.InitCameraPool(worldMeta.LevelMetas.CameraCount)
		// register scenes
		for _, scene := range worldMeta.LevelDetails.Scene {
			registerScene(scene.SceneName, &scene)
		}
		// load default scene
		doSceneLoad(&worldMeta.LevelDetails.Scene[0])
	}

	// start application
	appCfg := worldMeta.LevelMetas.ApplicationMetas
	return NewApplication(&AppConfig{
		Resolution:  linalg.NewVector2f64Ptr(appCfg.Resolution.W, appCfg.Resolution.H),
		PhysicalFps: appCfg.FPS.Physics,
		RenderFps:   appCfg.FPS.Render,
		Parallelism: appCfg.Parallelism,
		Title:       appCfg.Title,
		InitFunc:    initializer,
	})
}

func doSceneLoad(scene *parser.Scene) {
	// create objects in level
	for _, obj := range scene.ObjectDetails.Objects {
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
	// set camera pos
	for _, cam := range scene.SceneMetas.Cameras.Cameras {
		graphics.GetCamera(cam.Index).SetPos(cam.X, cam.Y)
	}
}

func doSceneUnload() {
	activePoolMu := mutexList[Mutex_ActivePool]
	activePoolMu.Lock()
	for _, pool := range activePool {
		for iobj := range pool {
			Destroy(iobj)
		}
	}
	activePoolMu.Unlock()
	inactivePoolMu := mutexList[Mutex_InactivePool]
	inactivePoolMu.Lock()
	for _, pool := range inactivePool {
		for iobj := range pool {
			Destroy(iobj)
		}
	}
	inactivePoolMu.Unlock()
}

// doChangeScene totally destroys current scene and switch to next one.
func doChangeScene(scene *parser.Scene) {
	// destory all old objects
	doSceneUnload()
	// load new objects
	doSceneLoad(scene)
}

// register scene adds a new scene to the world.
func registerScene(name string, sceneMeta *parser.Scene) {
	sceneCfgMap[name] = sceneMeta
}

// ChangeScene with provided name. Will panic if th provided name was not registered yet.
// Thread safe.
func ChangeScene(name string) {
	mu := mutexList[Mutex_SceneCfgMap]
	mu.Lock()
	if scene, ok := sceneCfgMap[name]; ok {
		doChangeScene(scene)
		mu.Unlock()
	} else {
		mu.Unlock()
		panic("scene not found")
	}
}
