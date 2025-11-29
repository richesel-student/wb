package main

import "fmt"

func quickSort(num []int) []int {
	if len(num) <= 2 {
		return num
	}

	pivot := num[len(num)/2]
	left := []int{}
	middle := []int{}
	right := []int{}

	for _, v := range num {
		if v == pivot {
			middle = append(middle, v)

		}
		if v < pivot {
			left = append(left, v)

		}
		if v > pivot {
			right = append(right, v)

		}
	}

	result := append(append(quickSort(left), middle...), quickSort(right)...)
	return result

}

func main() {
	var arrInt = []int{42, 7, 89, 13, 56, 4, 77, 21, 98, 30}

	fmt.Println(quickSort(arrInt))

}
