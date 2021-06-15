package physics

type IShape interface {
	Intersect(shape IShape) bool
}

type Rectangle struct {
	Width  float64
	Height float64
	Left   float64
	Top    float64
}

func (rect *Rectangle) Intersect(shape IShape) bool {
	switch shape := shape.(type) {
	case *Rectangle:
		return rect.IntersectWithRectangle(shape)
	}
	return false
}

func (rect *Rectangle) IntersectWithRectangle(anotherRect *Rectangle) bool {
	anotherRight := anotherRect.Left + anotherRect.Width
	anotherBottom := anotherRect.Top + anotherRect.Height
	thisRight := rect.Left + rect.Width
	thisBottom := rect.Top + rect.Height
	return thisRight >= anotherRect.Left && rect.Left <= anotherRight && rect.Top >= anotherBottom && thisBottom <= anotherRect.Top
}
