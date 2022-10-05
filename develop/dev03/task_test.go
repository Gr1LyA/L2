package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var in []string = []string{
	`bssd   1 ilya ilya 2231      Aug 21 23:40  Discord.lnk`,
	`sadasd 1 ilya ilya bbb512444 Sep 14 23:13  Games`,
	`bb     1 ilya ilya 5         Aug 22 00:30  инст`,
	`ala    1 ilya ilya 4         jul 14 23:13  Games`,
	`aa     1 ilya ilya 3         jul 14 23:13  Games`,
	`dd     1 ilya ilya 5122      feb 19:12  L2_golang`,
	`asd    1 ilya ilya 663b      May 11 19:30 'Visual Studio 2022.lnk'`,
	`c      1 ilya ilya 512       jun 19:12  L2_golang`,
	`asdd   1 ilya ilya 1         aug 26 13:46  доки`,
	`vv     1 ilya ilya -1        oct 14 23:13  Games`,
	`asd    1 ilya ilya 0         nov 10 08:04 'задания реп'`,
	`asd    1 ilya ilya 0         nov 10 08:04 'задания реп'`,
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
			arg: []string{"run", "task.go", "input"},
		},
		{
			arg: []string{"run", "task.go", "-k", "5", "-n", "input"},
		},
		{
			arg: []string{"run", "task.go", "-u", "input"},
		},
		{
			arg: []string{"run", "task.go", "-u", "-r", "input"},
		},
	}

	for _, v := range testCases {
		b1, _ := exec.Command("go", v.arg...).Output()
		b2, _ := exec.Command("sort", v.arg[2:]...).Output()
		if string(b1) != string(b2) {
			t.Fatal("\ncut output: ", string(b2), "\nmy cut output: ", string(b1))
		}
	}

	exec.Command("rm", "input").Run()
}
