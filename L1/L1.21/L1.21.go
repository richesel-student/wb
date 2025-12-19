package main

import (
	"fmt"
)

func strbyte(s string) []byte {
	strbyte := []byte(s)
	return strbyte

}
func reverse(b []byte, start int, end int) {
	for start < end {
		b[start], b[end] = b[end], b[start]
		start++
		end--
	}

}

func reverseWords(b []byte) []byte {
	start := 0 // начало текущего слова

	for i := 0; i <= len(b); i++ {

		if i == len(b) || b[i] == ' ' {
			left := start
			right := i - 1

			reverse(b, left, right)
			start = i + 1
		}
	}

	return b
}

func main() {
	strings := "snow dog sun"
	byteS := strbyte(strings)
	start := 0
	end := len(byteS) - 1
	reverse(byteS, start, end)
	result := reverseWords(byteS)
	fmt.Println(string(result))

}

