package graphics

import (
	"time"

	"galaxyzeta.io/engine/infra/logger"
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

	timing := time.Now()

	gl.Enable(gl.BLEND)
	GLEnableWireframe()
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	GLDeactivateTexture()

	logger.GlobalLogger.Debugf("p1 - cost: %v", time.Since(timing))

	vbo := vboManager.Borrow()

	GLBindData(vbo, vertices, len(vertices)*8, gl.DYNAMIC_DRAW)
	GLActivateShader("color")
	gl.DrawArrays(gl.QUADS, 0, 4)
	gl.Disable(gl.BLEND)
	GLDisableWireFrame()
	logger.GlobalLogger.Debugf("p2 - cost: %v", time.Since(timing))

	vboManager.Release(vbo)
	logger.GlobalLogger.Debugf("p3 - cost: %v", time.Since(timing))

}
