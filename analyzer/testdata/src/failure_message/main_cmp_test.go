package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type (
	User struct {
		name, surname string
		address       Address
	}

	Address struct {
		street, city, country string
	}
)

func TestCmpEqualSum(t *testing.T) {
	t.Parallel()

	want := 2
	got := double(1)
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want) // want `Prefer "double\(%v\) = %v, want %v" format for this failure message`
	}
}

func TestCmpDiff(t *testing.T) {
	t.Parallel()

	want := User{
		name:    "John",
		surname: "Doe",
		address: Address{},
	}
	got := double(1)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("diff %s", diff) // want `Prefer "diff -want \+got:\\n%s" format for this failure message`
	}
}
