package main

import (
	"testing"
)

// go test -v task.go task_test.go

func TestUnpackingString(t *testing.T) {
	testCase := []struct {
		in     string
		result string
	}{
		{in: `a4bc2d5e`,
			result: `aaaabccddddde`,
		},
		{
			in:     `abcd`,
			result: `abcd`,
		},
		{
			in:     `45`,
			result: ``,
		},
		{
			in:     ``,
			result: ``,
		},
		{
			in:     `qwe\4\5`,
			result: `qwe45`,
		},
		{
			in:     `qwe\45`,
			result: `qwe44444`,
		},
		{
			in:     `qwe\\5`,
			result: `qwe\\\\\`,
		},
	}

	for _, v := range testCase {
		res, _ := unpackingString(v.in)

		t.Logf("string: '%s', result: '%s'", v.in, res)

		if res != v.result {
			t.Errorf("Error: expected: '%s', got: '%s'", v.result, res)
		}
	}
}
