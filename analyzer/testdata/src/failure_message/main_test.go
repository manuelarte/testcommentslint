package main

import (
	"testing"
)

func sum(a, b int) int {
	return a + b
}

func sumAndBool(a, b int) (int, bool) {
	return a + b, true
}

func TestSum(t *testing.T) {
	t.Parallel()

	want := 2
	got := sum(1, 1)
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "sum\(%v, %v\) = %v, want %v" format for failure message`
	}
}

func TestSumAndBool(t *testing.T) {
	t.Parallel()

	want := 2
	got, _ := sumAndBool(1, 1)
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "sum\(%v, %v\) = %v, want %v" format for failure message`
	}
}
