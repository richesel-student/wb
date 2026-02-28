package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func match(line, pattern string, ignoreCase bool, fixed bool) bool {

	if ignoreCase {
		line = strings.ToLower(line)
		pattern = strings.ToLower(pattern)
	}

	if fixed {
		return strings.Contains(line, pattern)
	}

	re := regexp.MustCompile(pattern)
	return re.MatchString(line)
}

func main() {

	A := flag.Int("A", 0, "lines after")
	B := flag.Int("B", 0, "lines before")
	C := flag.Int("C", 0, "context")
	c := flag.Bool("c", false, "count only")
	i := flag.Bool("i", false, "ignore case")
	v := flag.Bool("v", false, "invert match")
	F := flag.Bool("F", false, "fixed string")
	n := flag.Bool("n", false, "show line numbers")

	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("usage: grep [flags] pattern file")
		return
	}

	pattern := args[0]
	filename := args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	before := *B
	after := *A

	if *C > 0 {
		before = *C
		after = *C
	}

	matchCount := 0

	for iLine, line := range lines {

		ok := match(line, pattern, *i, *F)

		if *v {
			ok = !ok
		}

		if ok {

			matchCount++

			if *c {
				continue
			}

			start := iLine - before
			if start < 0 {
				start = 0
			}

			end := iLine + after
			if end >= len(lines) {
				end = len(lines) - 1
			}

			for j := start; j <= end; j++ {

				if *n {
					fmt.Printf("%d:%s\n", j+1, lines[j])
				} else {
					fmt.Println(lines[j])
				}
			}

		}
	}

	if *c {
		fmt.Println(matchCount)
	}
}
