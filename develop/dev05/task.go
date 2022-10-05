package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type flags struct {
	// флаг C интегрирован в A и B
	A, B          int
	c, i, v, F, n bool
}

func main() {
	var src []string
	var err error

	f := parseFlags()
	args := flag.Args()

	// выбираю из чего читать, stdin или файлы
	switch len(args) {
	case 0:
		fmt.Fprintln(os.Stderr, "not enough arg")
		os.Exit(1)
	case 1:
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			src = append(src, scanner.Text())
		}
	default:
		src, err = openRead(args[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	pattern := args[0]
	grep(src, f, pattern)
}

func grep(src []string, f flags, pattern string) {
	length := len(src)

	// Реализован флаг c
	if f.c {
		var count int

		for i := 0; i < length; i++ {
			if compare(src[i], pattern, f) {
				count++
			}
		}
		fmt.Println(count)
	} else {
		var l int
		for i := 0; i < length; i++ {
			// захожу в условие если есть совпадение с паттерном
			if compare(src[i], pattern, f) {

				if l < i-1 && l != 0 && (f.A != 0 || f.B != 0) {
					fmt.Println("--")
				}

				//	Реализован флаг B (двигает левую границу)
				if l < i-f.B {
					l = i - f.B
				}
				if l < 0 {
					l = 0
				}

				// Флаг A реализован в этом цикле (сдвигает правую границу)
				for k := i + f.A; l <= k && l < length; l++ {
					printString(l+1, src[l], f)
					if compare(src[l], pattern, f) {
						k = l + f.A
					}
				}
				i = l - 1
			}
		}
	}
}

// реализован флаг i, v, F
func compare(s string, pattern string, f flags) (res bool) {
	// флаг i
	if f.i {
		s = strings.ToLower(s)
		pattern = strings.ToLower(pattern)
	}
	// флаг F
	if f.F {
		if s == pattern {
			res = true
		}
	} else if strings.Contains(s, pattern) {
		res = true
	}

	// флаг v
	if f.v {
		res = !res
	}

	return
}

//	флаг n
func printString(i int, s string, f flags) {
	if f.n {
		fmt.Print(i, ":", s, "\n")
	} else {
		fmt.Println(s)
	}
}

func openRead(args []string) (res []string, err error) {
	i := 0
	for _, v := range args {
		file, err := os.Open(v)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		for ; scanner.Scan(); i++ {
			res = append(res, scanner.Text())
		}
		file.Close()
	}

	return
}

func parseFlags() (res flags) {
	var tmp int

	flag.IntVar(&tmp, "C", 0, " \"context\" (A+B) печатать ±N строк вокруг совпадения")
	flag.IntVar(&res.A, "A", 0, " \"after\" печатать +N строк после совпадения")
	flag.IntVar(&res.B, "B", 0, " \"before\" печатать +N строк до совпадения")
	flag.BoolVar(&res.c, "c", false, " \"count\" (количество строк)")
	flag.BoolVar(&res.i, "i", false, " \"ignore-case\" (игнорировать регистр)")
	flag.BoolVar(&res.v, "v", false, " \"invert\" (вместо совпадения, исключать)")
	flag.BoolVar(&res.F, "F", false, " \"fixed\", точное совпадение со строкой, не паттерн")
	flag.BoolVar(&res.n, "n", false, " \"line num\", напечатать номер строки")

	flag.Parse()

	// реализован флаг C ( A и B в приоритете )
	if res.A == 0 {
		res.A = tmp
	}
	if res.B == 0 {
		res.B = tmp
	}

	if res.A < 0 || res.B < 0 {
		fmt.Fprintln(os.Stderr, "неверный аргумент длины контекста")
		os.Exit(1)
	}

	return
}
