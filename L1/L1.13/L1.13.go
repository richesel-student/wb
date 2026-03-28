package main

import (
	"errors"
	"fmt"
)

func replaceNum(A, B float64) (float64, float64) {

	fmt.Println("Числа до перестановки:", A, B)
	A = A + B
	B = A - B
	A = A - B
	return A, B
}

func main() {

	var num1, num2 float64
	_, ok := fmt.Scan(&num1)
	_, ok1 := fmt.Scan(&num2)

	if ok != nil || ok1 != nil {
		fmt.Println(errors.New("oшибка чтения"))

	} else {
		val1, val2 := replaceNum(num1, num2)
		fmt.Print("Числа после перестановки: ", val1, val2)

	}

}
