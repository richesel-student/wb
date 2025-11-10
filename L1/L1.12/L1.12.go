package main

import (
	"fmt"
)

func createSet(str []string) map[string]bool {
	set := make(map[string]bool)
	for _, s := range str {
		set[s] = true
	}
	return set

}

func printerSet(set map[string]bool) {
	var count int
	fmt.Printf("{")
	for s, _ := range set {
		count += 1
		fmt.Printf("\"%s\"", s)
		if count < len(set) {
			fmt.Printf(", ")

		}
	}
	fmt.Printf("}")

}

func main() {
	var strings = []string{"cat", "cat", "dog", "cat", "tree"}
	set := createSet(strings)
	printerSet(set)
}
