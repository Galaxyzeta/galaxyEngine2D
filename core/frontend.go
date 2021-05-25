package core

import (
	"fmt"
	"galaxyzeta.io/engine/input/keys"
	"runtime"
	"sync"

	"galaxyzeta.io/engine/linalg"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version = %v\n", version)
	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}

// keyboardCb is a function that will be used in OpenGL keyboard callback.
//
// We register and unregister KeyHold status when detecting KeyPress and KeyReleased
// instead of registering the KeyHold for a KeyHold callback and auto unregister it in subLoop automatically,
// because keyboard callback frequency is different from physical update frequency (also
// keyboard status reset frequency), which will cause undetected keys held situations
// during step update.
func keyboardCb(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		SetInputBuffer(keys.Action_KeyPress, keys.Key(key))
		SetInputBuffer(keys.Action_KeyHold, keys.Key(key))
	case glfw.Release:
		SetInputBuffer(keys.Action_KeyRelease, keys.Key(key))
		UnsetInputBuffer(keys.Action_KeyHold, keys.Key(key))
	}
}

func renderLoop(resolution *linalg.Vector2i, title string, renderFunc func(), wg *sync.WaitGroup, sigKill <-chan struct{}) {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic("Glfw init failed")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(resolution.X, resolution.Y, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	glfw.SwapInterval(1)

	window.SetKeyCallback(keyboardCb)

	for !window.ShouldClose() {
		// check sigkill
		select {
		case <-sigKill:
			wg.Done()
			return
		default:
		}
		// Do OpenGL stuff.
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// ---- render ----
		renderFunc()
		// ---- render ----

		window.SwapBuffers()
		glfw.PollEvents()
	}
	wg.Done()
}
