package graphics

import (
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func DrawRectangle(rect physics.Rectangle, color linalg.RgbaF64) {
	cam := GetCurrentCamera()
	vertices := []float64{
		rect.Left, rect.Top, 0, color.X, color.Y, color.Z, color.W,
		rect.Left, rect.Top + rect.Height, 0, color.X, color.Y, color.Z, color.W,
		rect.Left + rect.Width, rect.Top + rect.Height, 0, color.X, color.Y, color.Z, color.W,
		rect.Left + rect.Width, rect.Top, 0, color.X, color.Y, color.Z, color.W,
	}
	linalg.WorldVertice2OpenGL(&vertices, 0, 7, cam.Pos, cam.Resolution, GetScreenResolution())

	gl.Enable(gl.BLEND)
	GLEnableWireframe()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	GLDeactivateTexture()

	vbo := vboManager.Borrow()

	GLBindData(vbo, vertices, len(vertices)*8, gl.DYNAMIC_DRAW)
	GLActivateShader("color")
	gl.DrawArrays(gl.QUADS, 0, 4)
	gl.Disable(gl.BLEND)
	GLDisableWireFrame()

	vboManager.Release(vbo)
}

func DrawSegment(segment linalg.Segmentf64, color linalg.RgbaF64) {
	cam := GetCurrentCamera()
	vertices := []float64{
		segment.Point1.X, segment.Point1.Y, 0, color.X, color.Y, color.Z, color.W,
		segment.Point2.X, segment.Point2.Y, 0, color.X, color.Y, color.Z, color.W,
	}

	vbo := vboManager.Borrow()

	linalg.WorldVertice2OpenGL(&vertices, 0, 7, cam.Pos, cam.Resolution, GetScreenResolution())
	GLBindData(vbo, vertices, len(vertices)*8, gl.DYNAMIC_DRAW)
	GLActivateShader("color")
	gl.DrawArrays(gl.LINES, 0, 2)

	vboManager.Release(vbo)
}
