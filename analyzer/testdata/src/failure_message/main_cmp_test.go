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

func TestTableDrivenCmpEqualSum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{
			name: "simple case",
			want: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := double(1)
			if !cmp.Equal(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want) // want `Prefer "double\(%v\) = %v, want %v" format for this failure message`
			}
		})
	}
}

func TestCmpDiffWrongFormat(t *testing.T) {
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

func TestTableDrivenCmpDiffWrongFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want User
	}{
		{
			name: "simple example",
			want: User{
				name:    "John",
				surname: "Doe",
				address: Address{},
			},
		},
	}
	for _, test := range tests{
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := double(1)
			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("diff %s", diff) // want `Prefer "diff -want \+got:\\n%s" format for this failure message`
			}
		})
	}
}

func TestCmpDiffValidFormat(t *testing.T) {
	t.Parallel()

	want := User{
		name:    "John",
		surname: "Doe",
		address: Address{},
	}
	got := double(1)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("diff -want +got:\n%s", diff)
	}
}
