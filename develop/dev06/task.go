package main

// Флаг d просто меняет разделитель, по умолчанию \t
// Флаг f указывает столбцы которые нужно вывести (можно указывать диапазон, например -f 1-5, и через запятую, колонки сортируются по возрастатнию позициии)
// Флаг s
// Чтение либо из stdin либо из файла

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type flags struct {
	s bool
	d string
	// массив номеров колонок
	f []int
}

func main() {
	var in []string
	var fd *os.File = os.Stdin

	fl, numColumns := parseFlags()

	// если в аргументе есть файл меня fd стдина на fd файла
	if len(flag.Args()) != 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fd = file
	}

	// читаю из дескриптора и заношу построчно в массив
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		in = append(in, scanner.Text())
	}

	// делю строки по указанному разделителю
	src, maxNumColumns := spltStr(in, fl)

	// заполняю массив номеров колонок
	fl.f = initNumColumns(maxNumColumns, numColumns)

	// итерируюсь по множеству
	for _, v := range src {
		// итерируюсь по подмножеству
		for i, length, lenSrc := 0, len(fl.f), len(v); i < length && i < lenSrc; i++ {
			fmt.Print(v[fl.f[i]-1])
			// если надо печатаю разделитель
			if length > 1 && i != length-1 {
				fmt.Print(fl.d)
			}
		}
		fmt.Print("\n")
	}
}

func spltStr(src []string, fl flags) (res [][]string, maxNumColumns int) {
	for i := range src {
		tmp := strings.Split(src[i], fl.d)
		// если нет delimiter то не заношу в массив
		if length := len(tmp); length == 1 && fl.s {
			continue
		} else if length > maxNumColumns {
			// определяю максимальное кол-во колонок в строке
			maxNumColumns = length
		}
		// добавляю подмножество в массив
		res = append(res, tmp)
	}
	return
}

func parseFlags() (res flags, numColumns string) {
	flag.BoolVar(&res.s, "s", false, "только строки с разделителем")
	flag.StringVar(&res.d, "d", "\t", "использовать другой разделитель")
	flag.StringVar(&numColumns, "f", "", "выбрать поля (колонки)")
	flag.Parse()

	if numColumns == "" {
		fmt.Fprintln(os.Stderr, "you must specify fields")
		os.Exit(1)
	}

	if utf8.RuneCountInString(res.d) != 1 {
		fmt.Fprintln(os.Stderr, "the delimiter must be a single character")
		os.Exit(1)
	}

	return
}

// функция для того чтоб корректно распарсить , - в флаге f
func initNumColumns(maxNumColumns int, numColumns string) (result []int) {
	var err error
	var res []int

	// сначала делю по запятым
	strs := strings.Split(numColumns, ",")

	for i := range strs {
		// если встретил '-', то заношу числа в массив номеров колонок
		if strings.Contains(strs[i], "-") {
			var l, r int

			ranges := strings.Split(strs[i], "-")
			if len(ranges) > 2 {
				fmt.Fprintln(os.Stderr, "invalid field range")
				os.Exit(1)
			}
			// если левая граница не указана то она = 0
			if ranges[0] == "" {
				l = 0
			} else {
				l, err = strconv.Atoi(ranges[0])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
			// если правая граница не указана то она = maxNumColumns
			if ranges[1] == "" {
				r = maxNumColumns
			} else {
				r, err = strconv.Atoi(ranges[1])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if r > maxNumColumns {
					r = maxNumColumns
				}
			}
			if l > r {
				fmt.Fprintln(os.Stderr, "invalid decreasing range")
				os.Exit(1)
			}
			if l == 0 {
				fmt.Fprintln(os.Stderr, "fields are numbered from 1")
				os.Exit(1)
			}
			// добавлю номера колонок в массив из полученного диапазона
			for ; l <= r; l++ {
				res = append(res, l)
			}

		} else {
			// иначе просто заношу число в массив
			num, err := strconv.Atoi(strs[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			res = append(res, num)
		}
	}

	// сортирую массив номеров колонок
	sort.Ints(res)

	// оставляю только уникальные элементы
	for i := range res {
		if (i != 0 && res[i-1] != res[i]) || i == 0 {
			result = append(result, res[i])
		}
	}

	return
}
