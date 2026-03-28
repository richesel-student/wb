package main

import (
	"fmt"
)

func createArr(nums []int) []int {
	for i := 1; i < 101; i++ {
		nums = append(nums, i)
	}
	return nums

}

func readerArr(nums []int) <-chan int {
	reader := make(chan int)
	go func() {
		for i := 0; i < len(nums); i++ {
			reader <- nums[i]
		}
		close(reader)
	}()

	return reader

}

func squareChan(reader <-chan int) <-chan int {

	sqare := make(chan int)
	go func() {
		for num := range reader {
			sqare <- num * 2
		}
		close(sqare)
	}()
	return sqare
}

func out(sqare <-chan int) {
	for res := range sqare {
		fmt.Printf("%d ", res)

	}

}

func main() {

	var nums []int
	nums = createArr(nums)

	chanReaderArr := readerArr(nums)

	sqare := squareChan(chanReaderArr)

	out(sqare)

}
