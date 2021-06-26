package graphics

import (
	"image"
	"time"

	"galaxyzeta.io/engine/linalg"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type SpriteInstance struct {
	Animator
	frames  SpriteMeta
	OffsetX float64
	OffsetY float64
	Z       int
	vbo     uint32
}

// SpriteMeta is a sequence of frames that consists of an playable animation.
type SpriteMeta []*GLFrame

// GLFrame is a single img included in Sprite object.
type GLFrame struct {
	img       image.Image
	glTexture uint32
}

type Animator struct {
	currentFrame   int
	updateInterval time.Duration
	lastUpdateTime time.Time
	isPlaying      bool
}

func (spr *SpriteInstance) GetImg() image.Image {
	return spr.frames[spr.currentFrame].img
}

// NewFrame returns a new single frame.
func NewFrame(name string, frameFile string) *GLFrame {
	var img image.Image = nil
	var err error
	// load img
	img, err = ReadPng(frameFile)
	if err != nil {
		panic(err)
	}
	ret := &GLFrame{
		img: img,
	}
	GLRegisterTexture(img, &ret.glTexture)
	// save to graphic hashmap
	frameMap[name] = ret
	return ret
}

// NewSpriteMeta creates a new sprite meta from given sprite names.
func NewSpriteMeta(name string, frameNames ...string) {
	ret := make(SpriteMeta, len(frameNames))
	for idx := range ret {
		ret[idx] = GetFrame(frameNames[idx])
	}
	spriteMetaMap[name] = ret
}

// NewSpriteInstance creates a new sprite.
func NewSpriteInstance(sprMetaName string) (spr *SpriteInstance) {
	ret := &SpriteInstance{
		vbo:    GLNewVBO(1),
		frames: GetSpriteMeta(sprMetaName),
		Animator: Animator{
			currentFrame:   0,
			updateInterval: time.Millisecond * 200,
			lastUpdateTime: time.Now(),
			isPlaying:      true,
		},
	}
	return ret
}

// Render sprite. Sprite must exist.
func (spr *SpriteInstance) Render(camera *Camera, pos linalg.Point2f64) {
	currentGLImg := spr.frames[spr.currentFrame]
	dx := float64(currentGLImg.img.Bounds().Dx())
	dy := float64(currentGLImg.img.Bounds().Dy())
	vertices := []float64{
		pos.X, pos.Y, 0, 0, 0,
		pos.X, pos.Y + dy, 0, 0, 1,
		pos.X + dx, pos.Y + dy, 0, 1, 1,
		pos.X + dx, pos.Y, 0, 1, 0,
	}
	linalg.WorldVertice2OpenGL(&vertices, 0, 5, camera.Pos, camera.Resolution, GetScreenResolution())
	// vertices := []float64{
	// 	-0.5, 0.5, 0, 0, 0,
	// 	-0.5, -0.5, 0, 0, 1,
	// 	0.5, -0.5, 0, 1, 1,
	// 	0.5, 0.5, 0, 1, 0,
	// }
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	GLActivateTexture(currentGLImg.glTexture)
	GLBindData(spr.vbo, vertices, len(vertices)*8, gl.DYNAMIC_DRAW)
	GLActivateShader("default")
	gl.DrawArrays(gl.QUADS, 0, 4)
	gl.Disable(gl.BLEND)

}

func (spr *SpriteInstance) DoFrameStep() {
	if spr.isPlaying {
		if time.Since(spr.lastUpdateTime) >= spr.updateInterval {
			spr.currentFrame += 1
			if spr.currentFrame >= len(spr.frames) {
				spr.currentFrame = 0
			}
			spr.lastUpdateTime = time.Now()
		}
	}
}

func (spr *SpriteInstance) SetUpdateInterval(dur time.Duration) *SpriteInstance {
	spr.updateInterval = dur
	return spr
}

func (spr *SpriteInstance) EnableAnimation() *SpriteInstance {
	spr.isPlaying = true
	return spr
}

func (spr *SpriteInstance) DisableAnimation() *SpriteInstance {
	spr.isPlaying = false
	return spr
}

func (spr *SpriteInstance) IsPlaying() bool {
	return spr.isPlaying
}
