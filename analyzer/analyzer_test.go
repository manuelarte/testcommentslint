package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testCases := map[string]struct {
		patterns string
	}{
		"default": {
			patterns: "simple",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			a := New()

			analysistest.Run(t, analysistest.TestData(), a, test.patterns)
		})
	}
}
