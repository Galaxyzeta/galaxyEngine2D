package core

import "github.com/go-gl/glfw/v3.3/glfw"

// Key is a wrapper type for glfw.key
type Key glfw.Key

// keyboardCb is a function that will be used in OpenGL keyboard callback.
// TODO support modifierKey
func keyboardCb(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		registerKey(key, KeyPressed)
	case glfw.Repeat:
		registerKey(key, KeyHold)
	case glfw.Release:
		registerKey(key, KeyRelease)
	}
}

// registerKey writes key and actionType into inputBuffer.
func registerKey(key glfw.Key, actionType int) {
	inputBuffer[actionType][int(key)] = struct{}{}
}

func isKeyRegistered(key Key, actionType int) bool {
	_, ok := inputBuffer[actionType][int(key)]
	return ok
}

// flushInputBuffer resets input buffer to zero status.
func flushInputBuffer() {
	for _, actionKeyBinding := range inputBuffer {
		for k, _ := range actionKeyBinding {
			delete(actionKeyBinding, k)
		}
	}
}
