package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for {
		// ожидаем подключения
		conn, err := listener.Accept()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			conn.Close()
			continue
		}

		fmt.Println("connected")

		// Reader для чтения из буффера
		bufReader := bufio.NewReader(conn)
		fmt.Println("Start reading")

		go func(conn net.Conn) {
			defer conn.Close()
			for {
				// Читаем
				b, err := bufReader.ReadString('\n')

				if err != nil {
					fmt.Fprintln(os.Stderr, "can`t read", err)
					break
				}

				fmt.Fprintln(conn, "message:", string(b[:len(b)-1]), "received")

				fmt.Print(string(b))
			}
		}(conn)
	}
}
