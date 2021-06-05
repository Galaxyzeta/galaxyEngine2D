package graphics

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image"
)

type Sprite struct {
	OffsetX   float64
	OffsetY   float64
	img       image.Image
	glTexture uint32
	vbo       uint32
	Z         int
	lazyload  bool
}

func (spr *Sprite) GetImg() image.Image {
	return spr.img
}

func (spr *Sprite) GetLazyLoad() bool {
	return spr.lazyload
}

// NewSprite creates a new sprite.
func NewSprite(fileNamePng string, lazyLoad bool, OffsetX float64, OffsetY float64) *Sprite {
	var img image.Image = nil
	var err error
	if !lazyLoad {
		img, err = ReadPng(fileNamePng)
		if err != nil {
			panic(err)
		}
	}
	return &Sprite{
		OffsetX:  OffsetX,
		OffsetY:  OffsetY,
		img:      img,
		lazyload: lazyLoad,
		vbo:      GLNewVBO(1),
	}
}

// Render sprite. Sprite must exist.
func (spr *Sprite) Render(ox float64, oy float64) {
	if spr.img == nil {
		return
	}
	if spr.glTexture == 0 {
		fmt.Println("[System] textureActivated")
		GLSpriteRegister(spr.img, spr)
	}
	vertices := []float32{
		-0.5, 0.5, 0, 0, 0,
		-0.5, -0.5, 0, 0, 1,
		0.5, -0.5, 0, 1, 1,
		0.5, 0.5, 0, 1, 0,
	}
	GLActivateTexture(spr.glTexture)
	GLBindData(spr.vbo, vertices, len(vertices)*4, gl.DYNAMIC_DRAW)
	GLActivateShader("default")

	gl.DrawArrays(gl.QUADS, 0, 4)
}

// Depth gets sprite's z direction depth.
func (spr *Sprite) Depth() int {
	return spr.Z
}
