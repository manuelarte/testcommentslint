package checks

import (
	"testing"
)

func TestIsRecommendedGotWantFailureMessage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		functionName   string
		failureMessage string
		want           bool
	}{
		"valid function with one parameter": {
			functionName:   "YourFunction",
			failureMessage: "YourFunction(%v) = %v, want %v",
			want:           true,
		},
		"valid function with zero parameter": {
			functionName:   "YourFunction",
			failureMessage: "YourFunction() = %v, want %v",
			want:           true,
		},
		"valid function, no parenthesis": {
			functionName:   "YourFunction",
			failureMessage: "YourFunction = %v, want %v",
			want:           true,
		},
		"selector expr valid function, no parenthesis": {
			functionName:   "test.YourFunction",
			failureMessage: "YourFunction = %v, want %v",
			want:           true,
		},
		"selector expr wrong selector function, no parenthesis": {
			functionName:   "test.YourFunction",
			failureMessage: "x.YourFunction = %v, want %v",
			want:           false,
		},
		"two selector expr valid function, no parenthesis": {
			functionName:   "test.mystruct.YourFunction",
			failureMessage: "MyFunction = %v, want %v",
			want:           true,
		},
		"different function name with zero parameter": {
			functionName:   "MyFunction",
			failureMessage: "YourFunction() = %v, want %v",
			want:           false,
		},
		"got want, no function name": {
			functionName:   "MyFunction",
			failureMessage: "got %v, want %v",
			want:           false,
		},
		"expected, actual, no function name": {
			functionName:   "YourFunction",
			failureMessage: "actual: %v, expected %v",
			want:           false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := isRecommendedGotWantFailureMessage(tc.functionName, tc.failureMessage)
			if got != tc.want {
				t.Errorf("isRecommendedGotWantFailureMessage = %t, want %t", got, tc.want)
			}
		})
	}
}

func TestIsRecommendedDiffFailureMessage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		functionName   string
		failureMessage string
		want           bool
	}{
		"diff, expecting -want +got:\\n%s": {
			functionName:   "YourFunction",
			failureMessage: "diff: %s",
			want:           false,
		},
		"diff with function name, with (-want +got):\\n%s": {
			functionName:   "YourFunction",
			failureMessage: "YourFunction mismatch (-want +got):\n%s",
			want:           true,
		},
		"diff with function name, with -want +got:\\n%s": {
			functionName:   "YourFunction",
			failureMessage: "YourFunction mismatch -want +got:\n%s",
			want:           true,
		},
		"diff without function name, with -want +got:\\n%s": {
			functionName:   "YourFunction",
			failureMessage: "diff -want +got:\n%s",
			want:           true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := isRecommendedDiffFailureMessage(tc.functionName, tc.failureMessage)
			if got != tc.want {
				t.Errorf("isRecommendedDiffFailureMessage = %t, want %t", got, tc.want)
			}
		})
	}
}
