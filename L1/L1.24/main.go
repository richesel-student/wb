package main

import (
	"L1.24/point"
	"fmt"
)

func main() {

	point1 := point.NewPoint(1, 2)
	point2 := point.NewPoint(4, 6)
	distance := point1.Distance(point2)
	fmt.Print(distance)

}
