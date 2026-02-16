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
		t.Errorf("got %v, want %v", got, want) // want `Failure messages should include the name of the function that failed`
	}
}

func TestSumAndBool(t *testing.T) {
	t.Parallel()

	want := 2
	got, _ := sumAndBool(1, 1)
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Failure messages should include the name of the function that failed`
	}
}

func TestPrintHelloWorldValueCheck(t *testing.T) {
	t.Parallel()

	want := 10
	got, _ := printHelloWorld()
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Failure messages should include the name of the function that failed`
	}
}

func TestPrintHelloWorldErrCheck(t *testing.T) {
	t.Parallel()

	_, err := printHelloWorld()
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
}

func TestPrintHelloWorld(t *testing.T) {
	t.Parallel()

	want := 10
	got, err := printHelloWorld()
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want) // want `Failure messages should include the name of the function that failed`
	}

}
