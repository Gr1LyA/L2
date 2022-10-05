package main

// реализовать -n -h -b

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var month = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,
}

type flags struct {
	k                   int
	n, r, u, m, b, c, h bool
}

// реализует sort.Interface
type sortUtility struct {
	//	Флаги
	f flags
	//	Полные строки
	strs []string
	//	Конкретная колонка, если указан -k, иначе исходная строка
	column []string
}

func main() {
	var s sortUtility

	parseFlags(&s)

	// открываю, читаю, закрываю файл
	if err := openRead(&s); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Если указан флаг -c
	if s.f.c {
		isSorted(&s)
		return
	}

	// res - результат сортировки ([]string)
	res := sorting(&s)

	for i := range res {
		fmt.Println(res[i])
	}
}

func isSorted(s *sortUtility) {
	for i, j := 0, len(s.column); i < j-1; i++ {
		if !s.Less(i, i+1) {
			fmt.Println("disorder:", s.strs[i+1])
			return
		}
	}
}

func parseFlags(s *sortUtility) {
	flag.IntVar(&s.f.k, "k", 0, " указание колонки для сортировки")
	flag.BoolVar(&s.f.n, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&s.f.r, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&s.f.u, "u", false, "не выводить повторяющиеся строки")
	flag.BoolVar(&s.f.m, "M", false, "сортировка по названию месяца")
	flag.BoolVar(&s.f.b, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&s.f.c, "c", false, "проверять отсортированы ли данные")
	flag.BoolVar(&s.f.h, "h", false, "сортировать по числовому значению с учётом суффиксов")
	flag.Parse()
}

func openRead(s *sortUtility) error {
	if len(flag.Args()) == 0 {
		return errors.New("Нет пути к файлу")
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		s.strs = append(s.strs, scanner.Text())

		// если указана сортировка по колонке то вырезаю нужный кусок
		if s.f.k != 0 {
			columns := strings.Fields(s.strs[i])
			idx := s.f.k
			length := len(columns)

			if idx > length {
				idx = length
			}

			if idx <= 0 {
				s.column = append(s.column, "")
			} else {
				s.column = append(s.column, columns[idx-1])
			}
		} else {
			s.column = append(s.column, scanner.Text())
		}
	}

	return nil
}

func sorting(s *sortUtility) []string {
	sort.Sort(s)

	res := make([]string, 0, len(s.strs))

	for i := range s.strs {
		// если указан флаг -u и повторение перескакиваю на след итерацию
		if s.f.u && i != 0 && s.strs[i] == s.strs[i-1] {
			continue
		}
		res = append(res, s.strs[i])
	}

	return res
}

func (s *sortUtility) Len() int { return len(s.strs) }

func (s *sortUtility) Swap(i, j int) {
	s.strs[i], s.strs[j] = s.strs[j], s.strs[i]
	s.column[i], s.column[j] = s.column[j], s.column[i]
}

func (s *sortUtility) Less(i, j int) (res bool) {
	if s.f.n {
		res = s.LessNumber(i, j)
	} else if s.f.m {
		res = s.LessMonth(i, j)
	} else {
		res = s.column[i] < s.column[j]
	}

	// Если указан -r
	if s.f.r {
		return !res
	}
	return res
}

// если указан флаг -n
func (s *sortUtility) LessNumber(i, j int) (res bool) {
	var idx1, idx2 int
	for _, v := range s.column[i] {
		if !unicode.IsDigit(v) {
			break
		}
		idx1++
	}
	num1, _ := strconv.Atoi(s.column[i][:idx1])

	for _, v := range s.column[j] {
		if !unicode.IsDigit(v) {
			break
		}
		idx2++
	}
	num2, _ := strconv.Atoi(s.column[j][:idx2])

	if num1 == num2 || s.column[i] == "" || s.column[j] == "" {
		if s.column[i] == s.column[j] {
			res = s.strs[i] < s.strs[j]
		} else {
			res = s.column[i] < s.column[j]
		}
	} else {
		res = num1 < num2
	}
	return
}

// если указан флаг -M
func (s *sortUtility) LessMonth(i, j int) (res bool) {
	length1, length2 := utf8.RuneCountInString(s.column[i]), utf8.RuneCountInString(s.column[j])
	k, m := 3, 3

	if length1 < 3 {
		k = length1
	}
	if length2 < 3 {
		m = length2
	}

	val1, ok1 := month[strings.ToLower(s.column[i][:k])]
	val2, ok2 := month[strings.ToLower(s.column[j][:m])]

	switch {
	case !ok1 && !ok2:
		res = s.column[i] < s.column[j]
	case !ok1:
		res = true
	case !ok2:
		res = false
	default:
		res = val1 < val2
	}
	return
}
