package main

import (
	"errors"
	"fmt"
	"math/big"
)

type Value struct {
	A *big.Int
	B *big.Int
}

func (v *Value) Add() *big.Int {

	return new(big.Int).Add(v.A, v.B)
}

func (v *Value) Sub() *big.Int {

	return new(big.Int).Sub(v.A, v.B)
}

func (v *Value) Mul() *big.Int {

	return new(big.Int).Mul(v.A, v.B)

}

func (v *Value) Div() (*big.Int, error) {

	if len(v.B.Bits()) == 0 {
		return nil, errors.New("Деление на ноль")

	}

	return new(big.Int).Div(v.A, v.B), nil

}

func main() {
	numbers := Value{A: big.NewInt(1 << 27), B: big.NewInt(1 << 27)}
	fmt.Printf("\nСумма двух чисел:%v", numbers.Add())
	fmt.Printf("\nРазность двух чисел:%v", numbers.Sub())
	fmt.Printf("\nУмножение двух чисел:%v", numbers.Mul())

	res, err := numbers.Div()
	if err != nil {
		fmt.Print("\nОшибка:", err)
	} else {
		fmt.Println("\nДеление двух чисел:", res)

	}
}
