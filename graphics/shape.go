package graphics

type IShape interface {
	Intersect(shape IShape) bool
}

type Rectangle struct {
	width  float64
	height float64
	left   float64
	top    float64
}

func (rect *Rectangle) Intersect(shape IShape) bool {
	switch shape := shape.(type) {
	case *Rectangle:
		return rect.IntersectWithRectangle(shape)
	}
	return false
}

func (rect *Rectangle) IntersectWithRectangle(anotherRect *Rectangle) bool {
	anotherRight := anotherRect.left + anotherRect.width
	anotherBottom := anotherRect.top + anotherRect.height
	thisRight := rect.left + rect.width
	thisBottom := rect.top + rect.height
	return thisRight >= anotherRect.left && rect.left <= anotherRight && rect.top >= anotherBottom && thisBottom <= anotherRect.top
}
