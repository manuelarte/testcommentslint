package map_inlined

import "testing"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestMapInlined(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct { // want `Expected map-non-inlined table driven test`
		in int
		out int
	} {
		"test1": {
			in: 1,
			out: 1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := abs(test.in)
			if got != test.out {
				t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
			}
		})
	}
}

func TestMapNonInlined(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in int
		out int
	} {
		"test1": {
			in: 1,
			out: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := abs(test.in)
			if got != test.out {
				t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
			}
		})
	}
}

func TestSliceInlined(t *testing.T) {
	t.Parallel()

	for _, test := range []struct { // want `Expected map-non-inlined table driven test`
		name string
		in int
		out int
	} {
		{
			name: "test1",
			in: 1,
			out: 1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := abs(test.in)
			if got != test.out {
				t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
			}
		})
	}
}

func TestSliceNonInlined(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in int
		out int
	} {
		{
			name: "test1",
			in: 1,
			out: 1,
		},
	}
	for _, test := range tests { // want `Expected map-non-inlined table driven test`
		t.Run(test.name, func(t *testing.T) {
			got := abs(test.in)
			if got != test.out {
				t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
			}
		})
	}
}
