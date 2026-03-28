package main

import (
	"strconv"
	"strings"
)

// ParseFields разбирает строку с указанием полей 
// и возвращает map с номерами полей, которые нужно вывести.
func ParseFields(spec string) map[int]bool {

	fields := make(map[int]bool)

	parts := strings.Split(spec, ",")

	for _, p := range parts {

		if strings.Contains(p, "-") {

			r := strings.Split(p, "-")

			start, _ := strconv.Atoi(r[0])
			end, _ := strconv.Atoi(r[1])

			for i := start; i <= end; i++ {
				fields[i] = true
			}

		} else {

			n, _ := strconv.Atoi(p)
			fields[n] = true
		}
	}

	return fields
}

// ProcessLine обрабатывает одну строку входных данных,
// разделяет её по delimiter и возвращает выбранные поля.
// Если включён флаг separated и строка не содержит разделителя,
// строка пропускается.
func ProcessLine(line string, delimiter string, fields map[int]bool, separated bool) (string, bool) {

	if separated && !strings.Contains(line, delimiter) {
		return "", false
	}

	cols := strings.Split(line, delimiter)

	out := []string{}

	for i, v := range cols {
		if fields[i+1] {
			out = append(out, v)
		}
	}

	if len(out) == 0 {
		return "", false
	}

	return strings.Join(out, delimiter), true
}