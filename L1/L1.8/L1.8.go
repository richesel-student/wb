package main

import (
	"fmt"
	"strings"
)

// устанавливает указанный бит в 1
func sum(bit int64, variable_number int64) {
	fmt.Print("\n________________________________")
	fmt.Printf("\n%b — исходное число в двоичной форме", variable_number)
	fmt.Printf("\n%d — исходное число в десятичной форме", variable_number)
	fmt.Print("\n________________________________")

	var mask int64 = 1

	mask_copy := mask << bit                   // количество бит, на которое выполняется сдвиг
	oper_result := mask_copy | variable_number // установка бита в 1 (побитовое ИЛИ)

	fmt.Print("\n________________________________")
	fmt.Printf("\n%b — полученное число в двоичной форме", oper_result)
	fmt.Printf("\n%d — полученное число в десятичной форме", oper_result)
	fmt.Print("\n________________________________\n")
}

// сбрасывает указанный бит в 0
func andnot(bit int64, variable_number int64) {
	fmt.Print("\n________________________________")
	fmt.Printf("\n%b — исходное число в двоичной форме", variable_number)
	fmt.Printf("\n%d — исходное число в десятичной форме", variable_number)
	fmt.Print("\n________________________________")

	var mask int64 = 1

	mask_copy := mask << bit                    // количество бит, на которое выполняется сдвиг
	oper_result := variable_number &^ mask_copy // установка бита в 0 (побитовое И НЕ)

	fmt.Print("\n________________________________")
	fmt.Printf("\n%b — полученное число в двоичной форме", oper_result)
	fmt.Printf("\n%d — полученное число в десятичной форме", oper_result)
	fmt.Print("\n________________________________\n")
}

func input() {
	var number int64
	var bit int64
	var strvalue string

	fmt.Print("Введите число: ")
	fmt.Scanf("%d", &number)

	fmt.Print("Введите номер бита для изменения: ")
	fmt.Scanf("%d", &bit)
	if bit < 0 || bit > 63 {
		fmt.Println("Ошибка: номер бита должен быть в диапазоне 0..63")
		return
	}

	fmt.Print("Установить бит в 1? (да/нет): ")
	fmt.Scanf("%s", &strvalue)

	selection(number, bit, strvalue)
}

func selection(number int64, bit int64, strvalue string) {
	if strings.ToLower(strvalue) == "да" {
		sum(bit, number)
		if strings.ToLower(strvalue) == "нет" {
			andnot(bit, number)
		}

	} else {
		fmt.Println("Ошибка: нужно ввести 'да' или 'нет'.")
	}

}
func main() {
	input()

}
