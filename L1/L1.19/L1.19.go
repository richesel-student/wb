package main

import "fmt"

func reverseString(s string) {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		fmt.Printf("%s", string(runes[i]))
	}

}

func main() {
	var inputStr string
	fmt.Scanln(&inputStr)
	reverseString(inputStr)
}
