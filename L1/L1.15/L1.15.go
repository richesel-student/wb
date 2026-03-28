// Переменная justString ссылается на срез из первых 100 символов строки v,
// но при этом удерживает в памяти весь буфер, созданный функцией createHugeString.
// При многократных вызовах такой код может приводить к избыточному расходу памяти
// и выглядеть как утечка.
package main

import (
	"strconv"
	"strings"
)

var justString string

func createHugeString(num int) string {

	var myString string
	for i := 1; i <= num; i++ {
		myString += strconv.Itoa(i)
	}
	return myString

}

func someFunc() {
	v := createHugeString(1 << 10)
	justString = strings.Clone(v[:100])

}

func main() {

	createHugeString(1 << 10)
	someFunc()

}
