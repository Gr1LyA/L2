package main

import (
	"testing"
)

func TestShell(t *testing.T) {

	t.Log("launch: yes | head")
	errs := executePipeline([]string{"yes", "head"})
	for _, v := range errs {
		t.Fatal(v)
	}

	t.Log("launch: yes | yes | head")
	errs = executePipeline([]string{"yes", "yes", "head"})
	for _, v := range errs {
		t.Fatal(v)
	}

	t.Log("launch: echo -n alo")
	err := execute([]string{"echo", "-n", "alo"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("launch: echo alo")
	err = execute([]string{"echo", "alo"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("launch: ls")
	err = execute([]string{"ls"})
	if err != nil {
		t.Fatal(err)
	}

}
