package graphics

import (
	"galaxyzeta.io/engine/linalg"
	"github.com/go-gl/gl/v3.2-core/gl"
	"image"
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

func GLRenderSprite(x float64, y float64, sprite *Sprite) {
	var buffer uint32
	bounds := sprite.GetImg().Bounds()
	dx := float64(bounds.Dx())
	dy := float64(bounds.Dy())
	vertices := [4][2]float64{
		{x, y},
		{x+dx ,y},
		{x+dx, y+dy},
		{x, y+dy},
	}

	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.DrawArrays(gl.QUADS, 0, 4)
}
