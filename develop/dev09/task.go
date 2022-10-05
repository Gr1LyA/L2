package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Wget - Основная структура
type Wget struct {
	// Хранит в себе ссылку, ключ - ссылка, значение - была ли она скачана
	pages      sync.Map
	httpClient http.Client
	domain     string
	saveDir    string
	wg         sync.WaitGroup
}

// NewWget вернет Wget с заполнеными полями
func NewWget(domain string) *Wget {
	u, err := url.Parse(domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dir := u.Host
	if dir == "" {
		dir = domain
	}

	return &Wget{
		domain:     domain,
		saveDir:    dir,
		httpClient: http.Client{Transport: &http.Transport{MaxIdleConns: 100, IdleConnTimeout: 2 * time.Second}}}
}

// Run Запускает процесс загрузки сайта
func (w *Wget) Run() {
	// Создаю директорию в которую будет осуществялться загрузка файла
	os.MkdirAll(w.saveDir, os.ModePerm)

	// Перехожу в эту директроию
	if err := os.Chdir(w.saveDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Заношу эту дирректорию в мап,
	// которая хранит сайты которые были обработаны или нужно обработать
	w.pages.Store(w.domain, false)

	for {
		// Итерриуруюсь по мап, если встречаю не обработанную ссылку то загружаю ее
		w.pages.Range(func(k, v any) bool {
			if !v.(bool) {
				w.pages.Store(k, true)
				w.wg.Add(1)
				go w.download(k.(string))
			}
			return true
		})
		w.wg.Wait()
		fl := true
		// Проверяю есть ли еще не обработанные ссылки
		w.pages.Range(func(k, v any) bool {
			if !v.(bool) {
				fl = false
			}
			return true
		})
		if fl {
			return
		}

	}
}

func (w *Wget) download(pageURL string) {
	defer w.wg.Done()

	fmt.Println("start download:", pageURL)

	if length := len(pageURL); pageURL[length-1] == '/' {
		pageURL = pageURL[:length-1]
	}

	// Получаю старничку Get запросом
	dataPage, err := w.get(pageURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Вычисляю имя файла и путь, создаю нужные директории
	path, name, err := definePathAndName(pageURL)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Сохраняю страничку по найденному пути
	f, err := os.OpenFile(path+name, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	defer f.Close()

	// Записываю информацию полученную Get запросом в файл
	if _, err = f.Write(dataPage); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// Парсю страничку и заношу ссылки которые еще не были обработаны
	w.addLinksInMap(string(dataPage))
}

func definePathAndName(pageURL string) (path, name string, err error) {
	u, _ := url.Parse(pageURL)

	pathAndName := u.RequestURI()

	if pathAndName == "" {
		path = ""
		name = "index.html"
	} else {
		idxDel := strings.LastIndex(pathAndName, "/")

		if !strings.Contains(pathAndName, ".") {
			name = "index.html"
			path = "." + pathAndName + "/"
			if err := os.MkdirAll("."+pathAndName, os.ModePerm); err != nil {
				return "", "", err
			}
		} else {
			path = "." + pathAndName[:idxDel]

			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return "", "", err
			}
			path += "/"
			name = pathAndName[idxDel+1:]
		}
	}
	return
}

func (w *Wget) addLinksInMap(dataPage string) {
	r := regexp.MustCompile(`href="(.*?)"`)
	links := r.FindAllString(dataPage, -1)

	for i, length := 0, len(links); i < length; i++ {
		if strings.Index(links[i][6:], w.domain) == 0 {
			w.pages.LoadOrStore(links[i][6:len(links[i])-1], false)
		}
	}
}

func (w *Wget) get(pageURL string) ([]byte, error) {
	resp, err := w.httpClient.Get(pageURL)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp.Body.Close()
	return b, nil
}

func main() {
	var u string
	flag.StringVar(&u, "u", "", "url")
	flag.Parse()

	if u == "" {
		fmt.Fprintln(os.Stderr, "expected url")
	} else {
		NewWget(u).Run()
	}
}
