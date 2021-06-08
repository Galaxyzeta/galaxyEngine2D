package linalg

type Rgba struct {
	X uint32
	Y uint32
	Z uint32
	W uint32
}

type Vector2f64 struct {
	X float64
	Y float64
}

type Vector2i64 struct {
	X int64
	Y int64
}

type Vector2i32 struct {
	X int32
	Y int32
}

type Vector2f32 struct {
	X float32
	Y float32
}

type Vector2i struct {
	X int
	Y int
}

type Point2f32 Vector2f32
type Point2i64 Vector2i64
type Point2i Vector2i
