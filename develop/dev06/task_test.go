package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var in []string = []string{
	`Winter white snow 	frost`,
	`Spring:  green: grass:	warm`,
	`Summer: colorful: blossom:`,
	`Autumn: yellow: leaves: 	cool`,
}

func TestCut(t *testing.T) {
	f, err := os.OpenFile("input", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range in {
		fmt.Fprintln(f, v)
	}

	f.Close()

	testCases := []struct {
		arg []string
	}{
		{
			arg: []string{"run", "task.go", "-f", "1", "input"},
		},
		{
			arg: []string{"run", "task.go", "-f", "1-3", "-d", ":", "input"},
		},
		{
			arg: []string{"run", "task.go", "-f", "1-2, 3", "-d", ":", "-s", "input"},
		},
		{
			arg: []string{"run", "task.go", "-f", "1-2, 3", "-s", "input"},
		},
	}

	for _, v := range testCases {
		b1, _ := exec.Command("go", v.arg...).Output()
		b2, _ := exec.Command("cut", v.arg[2:]...).Output()
		if string(b1) != string(b2) {
			t.Fatal("cut output: ", string(b2), "my cut output: ", string(b1))
		}
	}

	exec.Command("rm", "input").Run()
}
