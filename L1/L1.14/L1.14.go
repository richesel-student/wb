package main

import "fmt"

func funcInterface(i interface{}) {

	switch v := i.(type) {
	case int:
		fmt.Printf("Тип данных int: %d\n", v)
	case string:
		fmt.Printf("Тип данных string: %s\n", v)

	case chan int:
		fmt.Printf("Тип данных канл int: %d\n", <-v)

	case bool:
		fmt.Printf("Тип данных bool: %t\n", v)

	}

}

func main() {

	in := make(chan int, 1)
	in <- 10

	funcInterface(21)
	funcInterface("hello")
	funcInterface(true)
	funcInterface(in)

}
