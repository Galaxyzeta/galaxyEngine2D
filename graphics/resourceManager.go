package graphics

import (
	"sync"

	"galaxyzeta.io/engine/linalg"
)

var spriteMetaMap map[string]SpriteMeta
var frameMap map[string]*GLFrame
var shaderMap map[string]*Shader
var screenResolution *linalg.Vector2f64 = &linalg.Vector2f64{}

var mutexList []sync.RWMutex

const (
	mutexScreenResolution = iota
)

func init() {
	shaderMap = make(map[string]*Shader)
	spriteMetaMap = make(map[string]SpriteMeta)
	frameMap = make(map[string]*GLFrame)

	mutexList = make([]sync.RWMutex, 1, 8)
	mutexList[mutexScreenResolution] = sync.RWMutex{}
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
