package main

import (
	"fmt"
	"testing"
)

func double(a int) int {
	return 2 * a
}

func sumAndBool(a, b int) (int, bool) {
	return a + b, true
}

func printHelloWorld() (int, error) {
	return fmt.Println("Hello World")
}

func TestSum(t *testing.T) {
	t.Parallel()

	want := 2
	got := double(1)
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "double\(%v\) = %v, want %v" format for this failure message`
	}
}

func TestSumAndBool(t *testing.T) {
	t.Parallel()

	want := 2
	got, _ := sumAndBool(1, 1)
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "sumAndBool\(%v, %v\) = %v, _, want %v" format for this failure message`
	}
}

func TestPrintHelloWorld(t *testing.T) {
	t.Parallel()

	want := 10
	got, _ := printHelloWorld()
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "printHelloWorld\(\) = %v, _, want %v" format for this failure message`
	}
}
