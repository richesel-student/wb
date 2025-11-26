package main

import "fmt"

func binarySearch(target int, arrInt []int) bool {

	last := len(arrInt) - 1
	first := 0

	for first <= last {
		middele := (last + first) / 2

		if arrInt[middele] == target {
			fmt.Println(arrInt[middele])
			return true
		} else {
			if arrInt[middele] < target {
				first = middele + 1
			} else {
				last = middele - 1
			}
		}

	}
	return false

}

func main() {
	// left:=[]int{}
	var arrInt []int
	// arrNum:=[]int{}

	for i := 25; i < 100; i++ {
		if i%2 == 0 {
			arrInt = append(arrInt, i)
		}
	}

	// fmt.Println(arrNum)

	target := 26

	search := binarySearch(target, arrInt)
	fmt.Print(search)
}
