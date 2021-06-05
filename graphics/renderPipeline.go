package graphics

// TODO use a general pipeline to hide troublesome details in OpenGL

type RenderPipeline struct {
	Texture      uint32
	Shader       uint32
	DrawFunction func()
}
