package graphics

import (
	"fmt"
	"image"
	"image/draw"
	"strings"

	"galaxyzeta.io/engine/infra/file"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"github.com/go-gl/gl/v4.1-core/gl"
)

// Shader is a representation of Shader / vao descriptor.
type Shader struct {
	shader        uint32       // shader in OpenGL descriptor
	vao           uint32       // vertex array object descriptor
	AttributeFunc func(uint32) // this is used for setting shader variable descriptions
}

// GLNewShader creates a new Shader.
func GLNewShader(name string, shader uint32, vao uint32, attr func(program uint32)) *Shader {
	s := &Shader{
		shader:        shader,
		vao:           vao,
		AttributeFunc: attr,
	}
	shaderMap[name] = s
	return s
}

func GLEnableWireframe() {
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

}

func GLDisableWireFrame() {
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

// GLActivateTexture uses texture for following drawings.
func GLActivateTexture(texture uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 3*4, 0)
	gl.EnableVertexAttribArray(0)
}

// GLDeactivateTexture removes texture binding.
func GLDeactivateTexture() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// GLActivateShader uses specific Shader program. This should be put at the end of everything ahead of drawing.
// to parse shader program parameters.
func GLActivateShader(name string) {
	shader := shaderMap[name]
	gl.UseProgram(shader.shader)
	gl.BindVertexArray(shader.vao)
	shader.AttributeFunc(shader.shader)
}

// GLBindData binds data into VBO.
// Notice the calling sequence, AttributeFunc should always be called
// after data binding, otherwise the shader will not working.
// DataSize is not the length of data interface{}, it is the total length of the data interface{}.
func GLBindData(vbo uint32, data interface{}, dataSize int, bindingMode uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, dataSize, gl.Ptr(data), bindingMode)
}

// GLMustPrepareShaderProgram reads content from vert and frag file, compiles Shader,
// creates Shader program and then link them up. Returns Shader program descriptor and err.
func GLMustPrepareShaderProgram(vert string, frag string) uint32 {
	fmt.Printf("[System] file = %s %s\n", vert, frag)
	vertexShaderSource, err := file.OpenAndRead(vert)
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := file.OpenAndRead(frag)
	if err != nil {
		panic(err)
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		panic(fmt.Errorf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

// GLNewVBO allocates an VBO.
func GLNewVBO(size int32) uint32 {
	var vbo uint32 = 0
	gl.GenBuffers(size, &vbo)
	return vbo
}

// GLReleaseVBO destory an VBO.
func GLReleaseVBO(vbo uint32) {
	gl.DeleteBuffers(1, &vbo)
}

// GLNewVAO allocates an VAO.
func GLNewVAO(size int32) uint32 {
	var vao uint32 = 0
	gl.GenVertexArrays(size, &vao)
	return vao
}

func GLRegisterTexture(img image.Image, slot *uint32) {

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// opengl generate texture
	gl.GenTextures(1, slot)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *slot)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func GLRenderRectangle(vbo uint32, rect physics.Rectangle, rgba linalg.Rgba) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	vertices := []float64{
		rect.Left, rect.Top, 0,
		rect.Left + rect.Width, rect.Top, 0,
		rect.Left + rect.Width, rect.Top + rect.Height, 0,
		rect.Left, rect.Top + rect.Height, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, 4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	gl.DrawArrays(gl.QUADS, 0, 4*2)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
