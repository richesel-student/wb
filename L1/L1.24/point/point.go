package point

import (
	"math"
)

type Point struct {
	x float64
	y float64
}

func NewPoint(a, b float64) *Point {
	return &Point{x: a, y: b}
}

func (point Point) Distance(other Point) float64 {
	Xexp := math.Pow((other.x - point.x), 2.0)
	Yexp := math.Pow((other.y - point.y), 2.0)
	return math.Sqrt(Xexp + Yexp)
}
