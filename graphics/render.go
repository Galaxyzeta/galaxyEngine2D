package graphics

import "galaxyzeta.io/engine/linalg"

type RenderDepthEnum int

const RenderDetph_Default RenderDepthEnum = 10
const RenderDepth_Top RenderDepthEnum = 100
const RenderDepth_Bottom RenderDepthEnum = 0

type IRenderable2D interface {
	Render(camera *Camera, pos linalg.Point2f32)
}
