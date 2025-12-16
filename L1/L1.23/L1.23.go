package main

import (
	"fmt"
)

func DeleteIndex(arr []int, index int) ([]int, error) {

	if index < 0 || index >= len(arr) {
		return arr, fmt.Errorf("Нет элемента в списке")
	}
	copy(arr[index:], arr[index+1:])
	arr = arr[:len(arr)-1]
	return arr, nil

}

func main() {
	arrInt := make([]int, 0, 99)

	for i := 1; i < 100; i++ {

		arrInt = append(arrInt, i)

	}
	del, err := DeleteIndex(arrInt, 1)

	if err != nil {
		fmt.Print("Ошибка:", err)
	}
	fmt.Print(del)
}
