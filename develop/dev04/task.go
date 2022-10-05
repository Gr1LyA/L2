package main

import (
	"sort"
	"strings"
)

func dictionary(src []string) map[string][]string {
	var found bool

	m := make(map[string][]string)

	arr := make([]string, len(src))

	//	Конвертирую в нижний регистр
	for i := range src {
		arr[i] = strings.ToLower(src[i])
	}

	// Итерируюсь по исходному массиву слов
	for i := range arr {
		found = false

		// пробегаюсь по мап, если слово анаграмма одного из ключей мап,
		// то добавляю в подмножество
		for k := range m {
			if anagram(k, arr[i]) {
				m[k] = append(m[k], arr[i])
				found = true
			}
		}
		// если не нашел совпадения то создаю новое подмножество и добавляю в него это слово
		if !found {
			m[arr[i]] = append(m[arr[i]], arr[i])
		}
	}

	// если длина подмножества = 1, то удаляю, иначе - сортирую
	for k := range m {
		if len(m[k]) == 1 {
			delete(m, k)
		} else {
			sort.Strings(m[k])
		}
	}

	return m
}

func anagram(arr1, arr2 string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for _, v := range arr1 {
		if strings.Count(arr1, string(v)) != strings.Count(arr2, string(v)) {
			return false
		}
	}

	return true
}
