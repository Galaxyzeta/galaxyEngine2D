package graphics

import (
	"image"
)

type Sprite struct {
	OffsetX   float64
	OffsetY   float64
	img       image.Image
	lazyload  bool
	glTexture uint32
	Z         int
}

func (spr *Sprite) GetImg() image.Image {
	return spr.img
}

func (spr *Sprite) GetLazyLoad() bool {
	return spr.lazyload
}

func NewSprite(fileName string, lazyLoad bool, OffsetX float64, OffsetY float64) *Sprite {
	var img image.Image = nil
	var err error
	if !lazyLoad {
		img, err = ReadPng(fileName)
		if err != nil {
			panic(err)
		}
	}
	return &Sprite{
		OffsetX:  OffsetX,
		OffsetY:  OffsetY,
		img:      img,
		lazyload: lazyLoad,
	}
}

// ---- Implement IRenderable2D ----

// Render sprite.
func (spr *Sprite) Render(ox float64, oy float64) {
	if spr.glTexture == 0 {
		// not registered
		GlSpriteRegister(spr.img, spr)
	}
}

// Z gets sprite's z direction depth.
func (spr *Sprite) Depth() int {
	return spr.Z
}
