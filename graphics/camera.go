package graphics

import "galaxyzeta.io/engine/linalg"

type Camera struct {
	Pos        linalg.Point2f64
	Resolution linalg.Vector2f64
}
