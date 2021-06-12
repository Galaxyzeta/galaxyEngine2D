package graphics

import (
	"image"
	"time"

	"galaxyzeta.io/engine/linalg"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Sprite struct {
	Animator
	OffsetX float64
	OffsetY float64
	Z       int
}

type GLImg struct {
	img       image.Image
	vbo       uint32
	glTexture uint32
}

type Animator struct {
	frames         []*GLImg
	currentFrame   int
	frameLen       int
	updateInterval time.Duration
	lastUpdateTime time.Time
	isPlaying      bool
}

func (spr *Sprite) GetImg() image.Image {
	return spr.frames[spr.currentFrame].img
}

// NewSprite creates a new sprite.
func NewSprite(OffsetX float64, OffsetY float64, updateInterval time.Duration, frameFiles ...string) (spr *Sprite) {
	var img image.Image = nil
	var err error
	// load img
	glImgFrames := make([]*GLImg, len(frameFiles))
	for idx, frame := range frameFiles {
		img, err = ReadPng(frame)
		if err != nil {
			panic(err)
		}
		glImgFrames[idx] = &GLImg{
			img: img,
			vbo: GLNewVBO(1),
		}
		GLRegisterTexture(img, &glImgFrames[idx].glTexture)
	}

	return &Sprite{
		OffsetX: OffsetX,
		OffsetY: OffsetY,
		Animator: Animator{
			frames:         glImgFrames,
			currentFrame:   0,
			frameLen:       len(glImgFrames),
			updateInterval: updateInterval,
			lastUpdateTime: time.Now(),
			isPlaying:      true,
		},
	}
}

// Render sprite. Sprite must exist.
func (spr *Sprite) Render(camera *Camera, pos linalg.Point2f32) {
	currentGLImg := spr.frames[spr.currentFrame]
	dx := float32(currentGLImg.img.Bounds().Dx())
	dy := float32(currentGLImg.img.Bounds().Dy())
	vertices := []float32{
		pos.X, pos.Y, 0, 0, 0,
		pos.X, pos.Y + dy, 0, 0, 1,
		pos.X + dx, pos.Y + dy, 0, 1, 1,
		pos.X + dx, pos.Y, 0, 1, 0,
	}
	linalg.WorldVertice2OpenGL(&vertices, 0, 5, camera.Pos, camera.Resolution, GetScreenResolution())
	// vertices := []float32{
	// 	-0.5, 0.5, 0, 0, 0,
	// 	-0.5, -0.5, 0, 0, 1,
	// 	0.5, -0.5, 0, 1, 1,
	// 	0.5, 0.5, 0, 1, 0,
	// }
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	GLActivateTexture(currentGLImg.glTexture)
	GLBindData(currentGLImg.vbo, vertices, len(vertices)*4, gl.DYNAMIC_DRAW)
	GLActivateShader("default")

	gl.DrawArrays(gl.QUADS, 0, 4)
	gl.Disable(gl.BLEND)

}

func (spr *Sprite) DoFrameStep() {
	if spr.isPlaying {
		if time.Since(spr.lastUpdateTime) >= spr.updateInterval {
			spr.currentFrame += 1
			if spr.currentFrame >= spr.frameLen {
				spr.currentFrame = 0
			}
			spr.lastUpdateTime = time.Now()
		}
	}
}

func (spr *Sprite) SetUpdateInterval(dur time.Duration) {
	spr.updateInterval = dur
}

func (spr *Sprite) EnableAnimation() {
	spr.isPlaying = true
}

func (spr *Sprite) DisableAnimation() {
	spr.isPlaying = false
}

func (spr *Sprite) IsPlaying() bool {
	return spr.isPlaying
}

// Depth gets sprite's z direction depth.
func (spr *Sprite) Depth() int {
	return spr.Z
}
