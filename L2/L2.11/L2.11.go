package main

import (
	"fmt"
	"sort"
	"strings"
)

// makeKey создаёт "отпечаток" слова
// сортирует буквы и возвращает новую строку
func makeKey(word string) string {
	r := []rune(word)
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return string(r)
}

// FindAnagrams находит множества анаграмм
func FindAnagrams(words []string) map[string][]string {
	groups := make(map[string][]string)   // временная группировка по ключу
	firstWord := make(map[string]string)  // хранит первое встретившееся слово

	for _, word := range words {
		word = strings.ToLower(word)

		key := makeKey(word)

		if _, exists := groups[key]; !exists {
			firstWord[key] = word
		}

		groups[key] = append(groups[key], word)
	}

	result := make(map[string][]string)

	for key, group := range groups {
		if len(group) > 1 {        // исключаем одиночные слова
			sort.Strings(group)   // сортируем слова внутри группы
			result[firstWord[key]] = group
		}
	}

	return result
}

func main() {
	words := []string{
		"пятак", "пятка", "тяпка",
		"листок", "слиток", "столик",
		"стол",
	}

	result := FindAnagrams(words)

	for key, group := range result {
		fmt.Println(key,":", group)
	}
}