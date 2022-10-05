package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var in []string = []string{
	`-rwxrwxrwx 1 ilya ilya    2231 Aug 21 23:40  Discord.lnk`,
	`drwxrwxrwx 1 ilya ilya     bbb512444 Sep 14 23:13  Games`,
	`drwxrwxrwx 1 ilya ilya     512 Aug 22 00:30  инст`,
	`drwxrwxrwx 1 ilya ilya     512 S 14 23:13  Games`,
	`drwxrwxrwx 1 ilya ilya     512 S 14 23:13  Games`,
	`drwxrwxrwx 1 ilya ilya     5122 15 19:12  L2_golang`,
	`drwxrwxrwx 1 ilya ilya     5122 15 19:12  L2_golang`,
	`-rwxrwxrwx 1 ilya ilya    1212 May 11 19:30 'Visual Studio 2022.lnk'`,
	`drwxrwxrwx 1 ilya ilya     512 buu 26 13:46  доки`,
	`drwxrwxrwx 1 ilya ilya     512 S 14 23:13  Games`,
	`drwxrwxrwx 1 ilya ilya     512 vuu 10 08:04 'задания реп'`,
}

func TestGrep(t *testing.T) {
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
			arg: []string{"run", "task.go", "-n", "Games", "input"},
		},
		{
			arg: []string{"run", "task.go", "-A", "1", "Games", "input"},
		},
		{
			arg: []string{"run", "task.go", "-A", "4", "Games", "input"},
		},
		{
			arg: []string{"run", "task.go", "-B", "3", "Games", "input"},
		},
		{
			arg: []string{"run", "task.go", "-C", "5", "Games", "input"},
		},
		{
			arg: []string{"run", "task.go", "-C", "5", "-v", "Games", "input"},
		},
		{
			arg: []string{"run", "task.go", "-C", "5", "-i", "gAmes", "input"},
		},
		{
			arg: []string{"run", "task.go", "-C", "5", "-i", "-v", "gAmes", "input"},
		},
		{
			arg: []string{"run", "task.go", "-C", "5", "-i", "-v", "-c", "gAmes", "input"},
		},
		{
			arg: []string{"run", "task.go", "-C", "5", "-i", "-v", "-c", "1", "input"},
		},
	}

	for _, v := range testCases {
		b1, _ := exec.Command("go", v.arg...).Output()
		b2, _ := exec.Command("grep", v.arg[2:]...).Output()
		if string(b1) != string(b2) {
			t.Fatal("grep output: ", string(b2), "my grep output: ", string(b1))
		}
	}

	exec.Command("rm", "input").Run()
}
