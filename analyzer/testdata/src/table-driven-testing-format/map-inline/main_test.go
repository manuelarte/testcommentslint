package map_inline

import "testing"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestMapInline(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
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
