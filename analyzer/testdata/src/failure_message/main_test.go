package main

import (
	"testing"
)

func sum(a, b int) int {
	return a + b
}

func TestNewMyStruct(t *testing.T) {
	t.Parallel()

	want := 2
	got := sum(1, 1)
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "sum(%v, %v) = %v, want %v" format for failure message`
	}
}
