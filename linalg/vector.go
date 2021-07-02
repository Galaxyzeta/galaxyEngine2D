package linalg

import "math"

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

type Segmentf64 struct {
	Point1 Vector2f64
	Point2 Vector2f64
}

type Point2f32 Vector2f32
type Point2f64 Vector2f64
type Point2i64 Vector2i64
type Point2i Vector2i

func (vec1 Vector2f64) Add(vec2 Vector2f64) Vector2f64 {
	return Vector2f64{vec1.X + vec2.X, vec1.Y + vec2.Y}
}

func (vec1 Vector2f64) Sub(vec2 Vector2f64) Vector2f64 {
	return Vector2f64{vec1.X - vec2.X, vec1.Y - vec2.Y}
}

func (vec1 Vector2f64) Dot(vec2 Vector2f64) float64 {
	return vec1.X*vec2.X + vec1.Y*vec2.Y
}

func (vec1 Vector2f64) Magnitude() float64 {
	return math.Sqrt(vec1.X*vec1.X + vec1.Y*vec1.Y)
}

func (vec1 Vector2f64) Normalize() Vector2f64 {
	magnitude := vec1.Magnitude()
	return Vector2f64{X: vec1.X / magnitude, Y: vec1.Y / magnitude}
}

func (vec1 Vector2f64) NormalVec() Vector2f64 {
	return Vector2f64{X: -vec1.Y, Y: vec1.X}
}

func (vec1 Vector2f64) ProjectOn(vec2 Vector2f64) Vector2f64 {
	mag := vec2.Magnitude()
	scale := vec1.Dot(vec2) / (mag * mag)
	return Vector2f64{X: vec2.X * scale, Y: vec2.Y * scale}
}
