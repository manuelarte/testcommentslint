package main

import "testing"

func TestSum(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		x, y, want int
	}{
		"one operand is 0": {
			x:    0,
			y:    1,
			want: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := add(tc.x, tc.y)
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}
