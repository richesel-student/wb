package main

import "fmt"

// В данной функции выполнится сначала x = 1, затем в defer func() инкрементируется на 1.
// В итоге x = 2, вернется return x = 2, так как есть
// возвращаемое значение (x int). Вывод: 2.
func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

// В данной функции:
// 1. Создается локальная переменная x = 0
// 2. x = 1 (локальной переменной присваивается 1)
// 3. return x - возвращается ТЕКУЩЕЕ значение x (копируется число 1)
// 4. Выполняется defer - x увеличивается до 2, но возвращаемое значение уже зафиксировано (1)
// Вывод: 1
func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
