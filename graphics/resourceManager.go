package graphics

import (
	"sync"

	"galaxyzeta.io/engine/linalg"
)

var spriteMetaMap map[string]SpriteMeta
var frameMap map[string]*GLFrame
var shaderMap map[string]*Shader
var vboManager *vboPool
var cameraPool []*Camera

var currentCamera int

var screenResolution *linalg.Vector2f64 = &linalg.Vector2f64{}

var mutexList []*sync.RWMutex

const (
	mutexScreenResolution = iota
	mutexVboManager
	mutexCurrentCamera
)

func init() {
	shaderMap = make(map[string]*Shader)
	spriteMetaMap = make(map[string]SpriteMeta)
	frameMap = make(map[string]*GLFrame)

	mutexList = make([]*sync.RWMutex, 0, 8)
	for i := 0; i < cap(mutexList); i++ {
		mutexList = append(mutexList, &sync.RWMutex{})
	}

}

func SetScreenResolution(x float64, y float64) {
	mutexList[mutexScreenResolution].Lock()
	screenResolution.X = x
	screenResolution.Y = y
	mutexList[mutexScreenResolution].Unlock()
}

func GetScreenResolution() linalg.Vector2f64 {
	mutexList[mutexScreenResolution].RLock()
	defer mutexList[mutexScreenResolution].RUnlock()
	return *screenResolution
}

// GetSpriteMeta gets a sequence of static frames from spriteMap.
func GetSpriteMeta(name string) SpriteMeta {
	spr, ok := spriteMetaMap[name]
	if !ok {
		panic("sprite not found !")
	}
	return spr
}

// GetFrame from frame map. Will panic if wanted frame is not found.
func GetFrame(name string) *GLFrame {
	frm, ok := frameMap[name]
	if !ok {
		panic("sprite not found !")
	}
	return frm
}

func GetCurrentCamera() *Camera {
	return GetCamera(currentCamera)
}

func GetCurrentCameraIndex() (index int) {
	mu := mutexList[mutexCurrentCamera]
	mu.RLock()
	index = currentCamera
	mu.RUnlock()
	return
}

func GetCamera(index int) *Camera {
	return cameraPool[index]
}

func SetCurrentCamera(index int) {
	if index > len(cameraPool) {
		panic("invalid index, should be less than the length of cameraPool")
	}
	mu := mutexList[mutexCurrentCamera]
	mu.Lock()
	currentCamera = index
	mu.Unlock()
}

// InitCameraPool inits camera pool with given camera counts. It will be called by core. Do not use this in ypur game logic.
func InitCameraPool(camCnt int) {
	// init camera list
	cameraPool = make([]*Camera, 0, camCnt)
	for i := 0; i < camCnt; i++ {
		cameraPool = append(cameraPool, NewCamera(linalg.NewVector2f64(0, 0), linalg.NewVector2f64(640, 480)))
	}
}

func GetVboManager() (ret *vboPool) {
	mutexList[mutexVboManager].RLock()
	ret = vboManager
	mutexList[mutexVboManager].RUnlock()
	return ret
}
