package graphics

import "galaxyzeta.io/engine/linalg"

type RenderDepthEnum int

const RenderDetph_Default RenderDepthEnum = 10
const RenderDepth_Top RenderDepthEnum = 100
const RenderDepth_Bottom RenderDepthEnum = 0

type IRenderable2D interface {
	Render(ox float64, oy float64)
	Depth() int
}

type RenderContext struct {
	Anchor     *linalg.Vector2f64
	Renderable IRenderable2D
}

// NewRenderContext receive
func NewRenderContext(r IRenderable2D, anchorPoint *linalg.Vector2f64) *RenderContext {
	return &RenderContext{
		Anchor:     anchorPoint,
		Renderable: r,
	}
}
