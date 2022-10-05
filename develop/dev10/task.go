package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type args struct {
	timeout time.Duration
	host    string
	port    string
}

func main() {
	f := parse()

	// Подключаюсь к указанному хосту по TCP
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(f.host, f.port), f.timeout)
	if err != nil {
		time.Sleep(f.timeout)
		fmt.Fprintln(os.Stderr, "Connection failed")
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// Читаю данные из сокета
			reader := bufio.NewReader(conn)
			b, err := reader.ReadString('\n')

			if err != nil {
				// Если EOF со стороны хоста
				if err == io.EOF {
					conn.Close()
					fmt.Println("Connection closed by foreign host.")
					os.Exit(0)
				}
				return
			}
			// вывожу в stdout данные из сокета
			fmt.Print(b)
		}
	}()

	// Пишу в сокет из  Stdin
	writer := bufio.NewWriter(conn)
	_, err = writer.ReadFrom(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	// Если EOF со стороны клиента
	fmt.Println("EOF")
	conn.Close()
	wg.Wait()
	fmt.Println("connection closed")
}

func parse() (f args) {
	flag.DurationVar(&f.timeout, "timeout", time.Second*10, "timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "expected host and port")
	}
	f.host = flag.Arg(0)
	f.port = flag.Arg(1)

	return
}
