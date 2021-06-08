package graphics

import "galaxyzeta.io/engine/linalg"

type Camera struct {
	Pos        linalg.Point2f32
	Resolution linalg.Vector2f32
}
