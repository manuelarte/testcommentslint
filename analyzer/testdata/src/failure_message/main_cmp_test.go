package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCmpSum(t *testing.T) {
	t.Parallel()

	want := 2
	got := double(1)
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "double\(%v\) = %v, want %v" format for a failure message`
	}
}
