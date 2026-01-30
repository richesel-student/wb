// Что выведет программа?
// Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil // type *os.PathError value  nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)        // <nil>, так как значение внутри интерфейса равно nil
	fmt.Println(err == nil) // Вывод false
}

// В Go непустой интерфейс представлен структурой iface,
// содержащей указатель на таблицу методов (*itab)
// и указатель на данные (unsafe.Pointer).
//
// Пустой интерфейс (interface{}) представлен структурой eface,
// не содержащей таблицы методов и способной хранить значение любого типа.
