package main

import (
	"fmt"
	"testing"
)

// реализуйте функцию Reverse
// ...

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func TestReverse(t *testing.T) {
	if fmt.Sprint(Reverse([]int{10, -6, 34, 54})) != "[54 34 -6 10]" {
		t.Errorf(`wrong []int reverse`)
	}
	if fmt.Sprint(Reverse([]string{"foo", "buzz", "generic", "go"})) != "[go generic buzz foo]" {
		t.Errorf(`wrong []string reverse`)
	}
	if fmt.Sprint(Reverse([]float64{4.67, 5, -2.34, 7.88, 100})) != "[100 7.88 -2.34 5 4.67]" {
		t.Errorf(`wrong []float64 reverse`)
	}
}
