package main

import (
	"sort"
	"testing"
)

func TestDictionary(t *testing.T) {
	testCase := []struct {
		in     []string
		length int
	}{
		{
			in:     []string{"alo", "aa", "aa", "пятка", "фыв", "bb", "выф", "пятак", "слиток", "листок", "столик", "тяпка"},
			length: 4,
		},
		{
			in:     []string{"абвг", "абвгд", "абввг", "абвг", "абввг"},
			length: 2,
		},
		{
			in:     []string{},
			length: 0,
		},
	}

	for _, v := range testCase {
		res := dictionary(v.in)
		t.Log("check", v.in, "\n\toutput dictionary: ", res)
		if len(res) != v.length {
			t.Errorf("expected length: %d, got : %d", v.length, len(res))
		}
		for _, val := range res {
			if len(val) == 0 {
				t.Error(val, " len == 0")
			}
			if !sort.StringsAreSorted(val) {
				t.Error(val, " is not sorted")
			}
		}
	}
}
