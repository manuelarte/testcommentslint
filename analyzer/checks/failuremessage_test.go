package checks

import "testing"

func TestIsRecommendedFailureMessage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		failureMessage string
		ifType         ifConditionType
		functionName   string
		want           bool
	}{
		"valid function with one parameter": {
			failureMessage: "YourFunction(%v) = %v, want %v",
			ifType:         equal,
			functionName:   "YourFunction",
			want:           true,
		},
		"valid function with zero parameter": {
			failureMessage: "YourFunction() = %v, want %v",
			ifType:         equal,
			functionName:   "YourFunction",
			want:           true,
		},
		"valid function, no parenthesis": {
			failureMessage: "YourFunction = %v, want %v",
			ifType:         equal,
			functionName:   "YourFunction",
			want:           true,
		},
		"different function name with zero parameter": {
			failureMessage: "YourFunction() = %v, want %v",
			ifType:         equal,
			functionName:   "MyFunction",
			want:           false,
		},
		"got want, no function name": {
			failureMessage: "got %v, want %v",
			ifType:         equal,
			functionName:   "MyFunction",
			want:           false,
		},
		"expected, actual, no function name": {
			failureMessage: "actual: %v, expected %v",
			ifType:         equal,
			functionName:   "MyFunction",
			want:           false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tfb := testFuncBlock{
				testedFunc: testFuncStmt{
					functionName: tc.functionName,
				},
				ifStmt: gotWantIfStmt{
					ifType: tc.ifType,
					errorCallExpr: tErrorfCallExpr{
						failureMessage: tc.failureMessage,
					},
				},
			}

			got := tfb.isRecommendedFailureMessage()
			if got != tc.want {
				t.Errorf("isRecommendedFailureMessage() = %v, want %v", got, tc.want)
			}
		})
	}
}
