package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		patterns string
		options  map[string]string
	}{
		"equality comparison": {
			patterns: "equality_comparison",
			options: map[string]string{
				FailureMessageCheckName: "false",
			},
		},
		"failure message": {
			patterns: "failure_message",
			options: map[string]string{
				EqualityComparisonCheckName: "false",
			},
		},
		"table-driven test format map-inline": {
			patterns: "table-driven-testing-format/map-inline",
			options: map[string]string{
				TableDrivenFormatCheckTypeName: "map",
				TableDrivenFormatCheckInlineName: "true",
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			a := New()

			for k, v := range test.options {
				err := a.Flags.Set(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, analysistest.TestData(), a, test.patterns)
		})
	}
}
