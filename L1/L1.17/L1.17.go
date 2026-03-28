package main

import (
	"fmt"
)

func binarySearch(target int, arr []int) int {

	last := len(arr) - 1
	first := 0

	for first <= last {

		middele := (last + first) / 2

		if arr[middele] == target {
			return middele
		} else {
			if arr[middele] < target {
				first = middele + 1
			} else {
				last = middele - 1
			}
		}

	}
	return -1

}

func main() {

	var arrInt []int

	for i := 25; i < 100; i++ {
		if i%2 == 0 {
			arrInt = append(arrInt, i)
		}
	}
	target := 28

	search := binarySearch(target, arrInt)
	fmt.Print(search)
}
