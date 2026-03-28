package main

import (
	"errors"
	"strings"
	"unicode"
)

func unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	flagDigit := false
	flagLetter := false

	for _, r := range s {
		if unicode.IsDigit(r) {
			flagDigit = true
		}
		if unicode.IsLetter(r) {
			flagLetter = true

		}
		if unicode.IsSpace(r) {
			return "", errors.New("(пустая строка -> пустая строка)")
		}

	}
	if flagDigit && flagLetter {
		result := ""

		for i := 1; i < len(s); i++ {
			// буква + цифра
			if unicode.IsLetter(rune(s[i-1])) && unicode.IsDigit(rune(s[i])) {
				part := strings.Repeat(string(s[i-1]), int(s[i]-'0'))
				result += part

				// буква + буква
			} else if i+1 < len(s) &&
				unicode.IsLetter(rune(s[i])) &&
				unicode.IsLetter(rune(s[i+1])) {

				result += string(s[i])
			}
		}

		// последняя буква
		if unicode.IsLetter(rune(s[len(s)-1])) {
			result += string(s[len(s)-1])
		}

		return result, nil
	}

	if !flagDigit && flagLetter {
		return s, nil
	}

	if flagDigit && !flagLetter {
		return "", errors.New("(некорректная строка, т.к. в строке только цифры)")
	}

	return "", nil
}
