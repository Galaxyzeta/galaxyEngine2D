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
	}
}

// ---- Implement IRenderable2D ----

// Render sprite. Sprite must exist.
func (spr *Sprite) Render(ox float64, oy float64) {
	if spr.img == nil {
		return;
	}
	if spr.glTexture == 0 {
		// not registered
		GlSpriteRegister(spr.img, spr)
	}
	GLRenderSprite(ox, oy, spr)
}

// Depth gets sprite's z direction depth.
func (spr *Sprite) Depth() int {
	return spr.Z
}
