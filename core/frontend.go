package core

import (
	"fmt"
	"runtime"
	"time"

	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/input/keys"
	"galaxyzeta.io/engine/linalg"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
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

func mouseButtonCb(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		SetInputBuffer(keys.Action_KeyPress, keys.Key(button))
		SetInputBuffer(keys.Action_KeyHold, keys.Key(button))
	case glfw.Release:
		SetInputBuffer(keys.Action_KeyRelease, keys.Key(button))
		UnsetInputBuffer(keys.Action_KeyHold, keys.Key(button))
	}
}

func cursorCb(w *glfw.Window, xpos float64, ypos float64) {
	mutexList[Mutex_CursorPos].Lock()
	cursorX = xpos
	cursorY = ypos
	mutexList[Mutex_CursorPos].Unlock()
}

// InitOpenGL will be called at the very beginning of the whole program.
func InitOpenGL(resolution linalg.Vector2f64, title string) *glfw.Window {
	// glfw init
	err := glfw.Init()
	if err != nil {
		panic("Glfw init failed")
	}

	// window creation
	window, err := glfw.CreateWindow(int(resolution.X), int(resolution.Y), title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(keyboardCb)
	window.SetMouseButtonCallback(mouseButtonCb)
	window.SetCursorPosCallback(cursorCb)
	glfw.SwapInterval(1)

	// opengl init
	if err = gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Printf("OpenGL version = %v\n", version)

	// install shaders
	installShaders()
	// init vbo pool
	graphics.InitVboPool(128)

	return window
}

func installShaders() {
	// handle default shader
	// TODO shader loading refactor, using configuration
	graphics.GLNewShader(
		"default",
		graphics.GLMustPrepareShaderProgram(fmt.Sprintf("%s/graphics/shaders/simpleVertex.glsl", GetCwd()), fmt.Sprintf("%s/graphics/shaders/simpleFragment.glsl", GetCwd())),
		graphics.GLNewVAO(1),
		func(program uint32) {
			// process uniform
			textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
			gl.Uniform1i(textureUniform, 0)
			// process input
			// -- position
			aPos := uint32(gl.GetAttribLocation(program, gl.Str("aPos\x00")))
			gl.VertexAttribPointerWithOffset(aPos, 3, gl.DOUBLE, false, 5*8, 0)
			gl.EnableVertexAttribArray(aPos)
			// -- uv
			texcoord := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
			gl.EnableVertexAttribArray(texcoord)
			gl.VertexAttribPointerWithOffset(texcoord, 2, gl.DOUBLE, false, 5*8, 3*8)
		})
	graphics.GLNewShader(
		"color",
		graphics.GLMustPrepareShaderProgram(fmt.Sprintf("%s/graphics/shaders/colorVertex.glsl", GetCwd()), fmt.Sprintf("%s/graphics/shaders/colorFragment.glsl", GetCwd())),
		graphics.GLNewVAO(1),
		func(program uint32) {
			// process input
			// -- position
			aPos := uint32(gl.GetAttribLocation(program, gl.Str("aPos\x00")))
			gl.VertexAttribPointerWithOffset(aPos, 3, gl.DOUBLE, false, 7*8, 0)
			gl.EnableVertexAttribArray(aPos)
			// -- color
			color := uint32(gl.GetAttribLocation(program, gl.Str("inputColor\x00")))
			gl.EnableVertexAttribArray(color)
			gl.VertexAttribPointerWithOffset(color, 4, gl.DOUBLE, false, 7*8, 3*8)
		})
	graphics.GLNewShader("noshader", 0, graphics.GLNewVAO(1), nil)
}

func RenderLoop(window *glfw.Window, renderFunc func(), sigKill <-chan struct{}) {

	fmt.Println("[System] renderLoop entered")

	defer glfw.Terminate()
	gl.ClearColor(0.5, 0.5, 1, 1)

	vboBufferCheckTimestamp := time.Now()

	for !window.ShouldClose() {

		// check sigkill
		select {
		case <-sigKill:
			return
		default:
		}
		// Do OpenGL stuff.
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// ---- exec pipeline cmd ----
		renderFunc()
		for _, gfxsys := range gfxSystemPriorityList {
			gfxsys.Execute(app.executor)
		}
		executeAdditionalDrawCalls()

		// ---- check buffer status ----
		tryAddVboStock(&vboBufferCheckTimestamp)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	fmt.Println("[System] OpenGL routine killed")
}

func executeAdditionalDrawCalls() {
	drawcalls := len(RenderCmdChan)
	for i := 0; i < drawcalls; i++ {
		// systemLogger.Debugf("exec external drawcall")
		(<-RenderCmdChan)()
	}
}

// tryAddVboStock checks whether there's enough vbo for allocation.
// if not, add more vbos to the stock.
// buffer allocation must be done on rendering thread,
// else, will cause panic
func tryAddVboStock(vboBufferCheckTimestamp *time.Time) {
	if time.Since(*vboBufferCheckTimestamp) > time.Second {
		*vboBufferCheckTimestamp = time.Now()
		vboManager := graphics.GetVboManager()
		if vboManager.Len() < 16 {
			vboManager.Enlarge(32)
		}
	}
}
