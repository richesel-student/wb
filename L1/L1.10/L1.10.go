package main

import (
	"fmt"
)

func uniqueKeys(nums []int) []int {
	counts := make(map[int]int)
	keys := []int{}
	for _, n := range nums {
		counts[n]++
	}
	for k := range counts {
		keys = append(keys, k)
	}
	return keys
}

func keyNum(arr []float64) []int {
	var list []int
	for _, j := range arr {
		keyNums := (int(j) / 10) * 10
		list = append(list, keyNums)
	}
	return list

}

func searchArr(arr []float64) ([]float64, []float64, []float64, []float64) {

	var list1, list2, list3, list4 []float64
	for _, i := range arr {
		if -29.9 <= i && i <= -20 {
			list1 = append(list1, i)
		}

		if 10 <= i && i <= 19.9 {
			list2 = append(list2, i)
		}

		if 20 <= i && i <= 29.9 {

			list3 = append(list3, i)
		}
		if 30 <= i && i <= 39.9 {
			list4 = append(list4, i)
		}

	}
	return list1, list2, list3, list4

}

func printer(m map[int][]float64) {
	first := true
	for key, values := range m {
		if !first {
			fmt.Print(", ")
		}
		first = false
		fmt.Printf("%d:{", key)
		for i, v := range values {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%.1f", v)
		}
		fmt.Print("}")
	}

}

func createMap(keys []int, mintwenty []float64, ten []float64, twenty []float64, thirty []float64) map[int][]float64 {
	m := make(map[int][]float64)
	for _, key := range keys {
		switch key {
		case -20:
			m[key] = mintwenty
		case 10:
			m[key] = ten
		case 20:
			m[key] = twenty
		case 30:
			m[key] = thirty
		}
	}
	return m

}

func main() {
	arrMap := make(map[int][]float64)
	var arr = []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	listKeyNums := keyNum(arr)
	keys := uniqueKeys(listKeyNums)
	minustwenty, ten, twenty, thirty := searchArr(arr)
	arrMap = createMap(keys, minustwenty, ten, twenty, thirty)
	printer(arrMap)

}
