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
		},
		"failure message": {
			patterns: "failure_message",
			options: map[string]string{
				EqualityComparisonCheckName: "false",
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
