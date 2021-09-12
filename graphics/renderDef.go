package graphics

type IRenderable interface {
	Render(cam *Camera)
	PostRender()
	IsStatic() bool
	Z() int64
}
