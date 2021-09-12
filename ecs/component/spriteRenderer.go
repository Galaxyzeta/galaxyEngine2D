package component

import (
	"galaxyzeta.io/engine/graphics"
	"galaxyzeta.io/engine/linalg"
	"galaxyzeta.io/engine/physics"
)

const NameSpriteRenderer string = "spriteRenderer"

type SpriteRenderer struct {
	*graphics.Animator
	Name     string
	tf       *Transform2D
	z        int64             // render depth. bigger Z means deeper behind current view.
	Scale    linalg.Vector2f64 // scale value on X and Y
	Pivot    *physics.Pivot
	Offset   linalg.Vector2f64 // the offset from sprite to player, negative value means drawing at left of an object.
	isStatic bool              // read only, if marked true, will not involve in Z-depth sorting.
	Enabled  bool              // is visible or not
}

// GetName returns sprite renderer's name.
func (sr *SpriteRenderer) GetName() string {
	return sr.Name
}

// NewSpriteRenderer returns a new renderer that render sprites.
func NewSpriteRenderer(animator *graphics.Animator, tf *Transform2D, isStatic bool) *SpriteRenderer {
	return &SpriteRenderer{
		Animator: animator,
		tf:       tf,
		Enabled:  true,
		Scale:    linalg.NewVector2f64(1, 1),
		Name:     NameSpriteRenderer,
		Pivot: &physics.Pivot{
			Option: physics.PivotOption_TopLeft,
		},
		isStatic: isStatic,
	}
}

func NewSpriteRendererWithOptions(animator *graphics.Animator, tf *Transform2D, isStatic bool, options graphics.RenderOptions) (sr *SpriteRenderer) {
	sr = NewSpriteRenderer(animator, tf, isStatic)
	if options.Scale != nil {
		sr.Scale = *options.Scale
	}
	if options.Pivot != nil {
		sr.Pivot = options.Pivot
	}
	return sr
}

// Spr returns the sprite instance.
func (sr *SpriteRenderer) Spr() *graphics.SpriteInstance {
	return sr.Animator.Spr()
}

func (sr *SpriteRenderer) Render(cam *graphics.Camera) {
	sr.Spr().Render(cam, sr.tf.Pos, graphics.RenderOptions{
		Scale: &sr.Scale,
		Pivot: sr.Pivot,
	})
}

func (sr *SpriteRenderer) IsStatic() bool {
	return sr.isStatic
}

func (sr *SpriteRenderer) Z() int64 {
	return sr.z
}

func (sr *SpriteRenderer) PostRender() {
	sr.Spr().DoFrameStep()
}

func (sr *SpriteRenderer) SetZ(z int64) {
	sr.z = z
}

func (sr *SpriteRenderer) GetHitbox() physics.Polygon {
	return sr.Animator.Spr().GetHitbox(&sr.tf.Pos, physics.Pivot{Option: sr.Pivot.Option})
}
