package main

import (
	"reflect"
	"testing"
)

func TestNewMyStruct(t *testing.T) {
	t.Parallel()

	want := MyStruct{
		id: 0,
		name: "John",
	}
	got := NewMyStruct(want.id, want.name)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTableDrivenTestNewMyStruct(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		id int
		name string
		want MyStruct
	}{
		"normal values": {
			id: 0,
			name: "John",
			want: MyStruct{
				id: 0,
				name: "John",
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewMyStruct(tc.id, tc.name)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
