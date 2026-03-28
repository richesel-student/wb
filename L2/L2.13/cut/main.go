package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {

	fieldsFlag := flag.String("f", "", "fields")
	delimiter := flag.String("d", "\t", "delimiter")
	separated := flag.Bool("s", false, "only separated")

	flag.Parse()

	if *fieldsFlag == "" {
		fmt.Println("flag -f required")
		os.Exit(1)
	}

	fields := ParseFields(*fieldsFlag)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		res, ok := ProcessLine(
			scanner.Text(),
			*delimiter,
			fields,
			*separated,
		)

		if ok {
			fmt.Println(res)
		}
	}
}
