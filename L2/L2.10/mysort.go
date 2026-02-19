package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// options хранит значения флагов программы.
type options struct {
	column  int  // номер колонки для сортировки (-k)
	numeric bool // числовая сортировка (-n)
	reverse bool // обратный порядок (-r)
	unique  bool // вывод только уникальных строк (-u)
}

func main() {
	opts := parseFlags()

	lines, err := readInput(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	sortLines(lines, opts)

	if opts.unique {
		lines = unique(lines)
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}

// parseFlags читает флаги командной строки.
func parseFlags() options {
	var opts options

	flag.IntVar(&opts.column, "k", 0, "sort by column N (tab separated)")
	flag.BoolVar(&opts.numeric, "n", false, "numeric sort")
	flag.BoolVar(&opts.reverse, "r", false, "reverse sort")
	flag.BoolVar(&opts.unique, "u", false, "unique lines only")

	flag.Parse()
	return opts
}

// readInput читает строки из файла или STDIN.
func readInput(filename string) ([]string, error) {
	var reader io.Reader

	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	var lines []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// sortLines сортирует строки в соответствии с флагами.
func sortLines(lines []string, opts options) {
	sort.SliceStable(lines, func(i, j int) bool {
		a := extractKey(lines[i], opts)
		b := extractKey(lines[j], opts)

		var less bool

		if opts.numeric {
			less = toFloat(a) < toFloat(b)
		} else {
			less = a < b
		}

		if opts.reverse {
			return !less
		}
		return less
	})
}

// extractKey возвращает ключ сортировки.
// Если указан -k, используется соответствующая колонка (разделитель — табуляция).
func extractKey(line string, opts options) string {
	if opts.column <= 0 {
		return line
	}

	fields := strings.Split(line, "\t")
	if opts.column-1 < len(fields) {
		return fields[opts.column-1]
	}
	return ""
}

// unique удаляет подряд идущие одинаковые строки.
func unique(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result
}

// toFloat безопасно преобразует строку в число.
// Если преобразование не удалось — возвращается 0.
func toFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}