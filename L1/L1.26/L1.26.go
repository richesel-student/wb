package main

import (
	"fmt"
	"strings"
)

func mapSearch(s1 string) bool {
	lowerstr := strings.ToLower(s1)
	strune1 := []rune(lowerstr)
	mapSearch := make(map[rune]bool)
	for _, key := range strune1 {
		mapSearch[key] = true
	}
	if len(mapSearch) == len(strune1) {
		return true
	}
	return false
}

func main() {
	str1 := "abcd"
	str2 := "abCdefAaf"
	str3 := "aabcd"

	fmt.Println(mapSearch(str1))
	fmt.Println(mapSearch(str2))
	fmt.Println(mapSearch(str3))

}
