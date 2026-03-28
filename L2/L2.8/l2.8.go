package main

import (
	"fmt"
	"mylib/mylib"
	"os"
)

func main() {
	t, err := mylib.TimeNTP()
	if err != nil {
		fmt.Fprintf(os.Stderr, "NTP error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(t)
}