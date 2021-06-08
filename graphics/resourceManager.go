package graphics

import (
	"sync"

	"galaxyzeta.io/engine/linalg"
)

var spriteMap map[string]*Sprite
var shaderMap map[string]*Shader
var screenResolution *linalg.Vector2f32 = &linalg.Vector2f32{}

var mutexList []sync.RWMutex

const (
	mutexScreenResolution = iota
)

func init() {
	shaderMap = make(map[string]*Shader, 0)
	mutexList = make([]sync.RWMutex, 1, 8)
	mutexList[mutexScreenResolution] = sync.RWMutex{}
}

func SetScreenResolution(x float32, y float32) {
	mutexList[mutexScreenResolution].Lock()
	screenResolution.X = x
	screenResolution.Y = y
	mutexList[mutexScreenResolution].Unlock()
}

func GetScreenResolution() linalg.Vector2f32 {
	mutexList[mutexScreenResolution].RLock()
	defer mutexList[mutexScreenResolution].RUnlock()
	return *screenResolution
}
