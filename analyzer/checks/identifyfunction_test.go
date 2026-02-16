package checks

import (
	"testing"
)

func TestContainsFunctionNameString(t *testing.T) {
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
		"selector expr complete function, no parenthesis": {
			functionName:   "test.YourFunction",
			failureMessage: "test.YourFunction = %v, want %v",
			want:           true,
		},
		"selector expr wrong selector function, no parenthesis": {
			functionName:   "test.YourFunction",
			failureMessage: "x.YourFunction = %v, want %v",
			want:           false,
		},
		"two selector expr valid function, no parenthesis": {
			functionName:   "test.mystruct.YourFunction",
			failureMessage: "YourFunction = %v, want %v",
			want:           true,
		},
		"two selector expr complete function, no parenthesis": {
			functionName:   "test.mystruct.YourFunction",
			failureMessage: "test.mystruct.YourFunction = %v, want %v",
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

			got := containsFunctionNameString(tc.functionName, tc.failureMessage)
			if got != tc.want {
				t.Errorf("isRecommendedGotWantFailureMessage = %t, want %t", got, tc.want)
			}
		})
	}
}
