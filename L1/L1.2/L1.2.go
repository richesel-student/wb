package main

import (
	"fmt"
)

func square(ch chan int, num int) {
	ch <- num * num

}

func runSquares(ch chan int, array []int) {
	for _, value := range array {
		go square(ch, value)

	}

}

func resultSquares(ch chan int, array []int) {
	for i := 0; i < len(array); i++ {
		fmt.Println(<-ch)

	}

}

func main() {
	array := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(array))
	runSquares(ch, array)
	resultSquares(ch, array)

}
