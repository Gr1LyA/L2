package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

func main() {
	t, err := ntp.Time("ntp1.stratum2.ru")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(127)
	}

	fmt.Println(t)
}
