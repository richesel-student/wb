//Что выведет программа?

//Объяснить вывод программы.

package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

// Результат программы будет "error", т.к. err имеет интерфейсный тип error
// и содержит динамический тип *customError и значение nil.
// Поэтому условие err != nil будет истинным.
