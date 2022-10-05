package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func unpackingString(src string) (string, error) {
	str := []rune(src)
	length := len(str)

	if length == 0 {
		return "", nil
	}

	if unicode.IsDigit(str[0]) {
		return "", errors.New("invalid string")
	}

	var res strings.Builder

	res.Grow(length)

	for i := 0; i < length; i++ {
		//	если встретили цифру
		if unicode.IsDigit(str[i]) {
			j := i + 1
			//	ищу правую граница числа
			for j < length && unicode.IsDigit(str[j]) {
				j++
			}

			//	конвертирую в int
			count, err := strconv.Atoi(string(str[i:j]))
			if err != nil {
				return "", err
			}

			//	повторяю символ после числа нужое кол-во раз
			for ; count > 1; count-- {
				res.WriteString(string(str[i-1]))
			}

			i = j - 1

		} else {
			//	если экранирование то перепрыгиваю через символ
			if str[i] == '\\' {
				i++
			}
			res.WriteString(string(str[i]))
		}
	}

	return res.String(), nil
}
