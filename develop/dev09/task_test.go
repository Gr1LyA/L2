package main

import (
	"testing"
)

func TestDefinePathAndName(t *testing.T) {
	testCases := []struct {
		testCase string
		expected string
	}{
		{
			testCase: "https://test",
			expected: ".//index.html",
		},
		{
			testCase: "https://test/two",
			expected: "./two/index.html",
		},
		{
			testCase: "https://test/zone/three",
			expected: "./zone/three/index.html",
		},
	}

	for _, v := range testCases {
		path, name, _ := definePathAndName(v.testCase)
		if path+name != v.expected {
			t.Errorf("expected '%s', got '%s'\n", v.expected, path+name)
		}
	}
}
