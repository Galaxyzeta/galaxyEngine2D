package graphics

import (
	"image"

	"galaxyzeta.io/engine/linalg"
	"github.com/go-gl/gl/v3.2-core/gl"
)

func GlSpriteRegister(img image.Image, spr *Sprite) {
	rect := img.Bounds()
	w := rect.Dx()
	h := rect.Dy()
	// init rgba 2d array
	rgba := make([][]linalg.Rgba, h)
	for idx := range rgba {
		rgba[idx] = make([]linalg.Rgba, w)
	}
	// populate rgba 2d array
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			r, g, b, a := img.At(dx, dy).RGBA()
			rgba[dx][dy] = linalg.Rgba{X: r, Y: g, Z: b, W: a}
		}
	}
	// opengl generate texture
	gl.GenTextures(1, &spr.glTexture)
}

func GLQuad() {
	// gl.DrawArrays(gl.QUADS)
}
