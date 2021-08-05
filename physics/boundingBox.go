package physics

import "galaxyzeta.io/engine/linalg"

type BoundingBox [4]linalg.Vector2f64

const BB_TopLeft = 1
const BB_TopRight = 0
const BB_BotLeft = 2
const BB_BotRight = 3

func (bb BoundingBox) ToRectangle() Rectangle {
	return Rectangle{
		Width:  bb[BB_TopRight].X - bb[BB_TopLeft].X,
		Height: bb[BB_BotLeft].Y - bb[BB_TopLeft].Y,
		Left:   bb[BB_TopLeft].X,
		Top:    bb[BB_TopLeft].Y,
	}
}

func (bb BoundingBox) GetBottomLeftPoint() linalg.Vector2f64 {
	return bb[BB_BotLeft]
}

func (bb BoundingBox) GetBottomRightPoint() linalg.Vector2f64 {
	return bb[BB_BotRight]
}

func (bb BoundingBox) GetTopLeftPoint() linalg.Vector2f64 {
	return bb[BB_TopLeft]
}

func (bb BoundingBox) GetTopRightPoint() linalg.Vector2f64 {
	return bb[BB_TopRight]
}

func SliceToBoundingBox(vecSlice []linalg.Vector2f64) BoundingBox {
	array := [4]linalg.Vector2f64{}
	copy(array[:], vecSlice)
	return BoundingBox(array)
}
