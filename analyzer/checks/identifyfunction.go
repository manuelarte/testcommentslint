package checks

import (
	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// IdentifyFunction check that the failure messages in t.Errorf follow the format expected.
// The format expected can be as the following:
//   - When the condition is `reflect.DeepEqual`, `cmp.Equal` or `got != want`: "YourFunction(%v) = %v, want %v"
//   - When the condition is `cmp.Diff`: YourFunction(%v) mismatch (-want +got):\n%s
//
// This checks blocks like the following:
//
//	 got := MyFunction(in)
//		if got != want {
//		  t.Errorf(...)
//		}
type IdentifyFunction struct {
	category string
}

// NewIdentifyFunction creates a new IdentifyFunction.
func NewIdentifyFunction() IdentifyFunction {
	return IdentifyFunction{
		category: "Identify The Function",
	}
}

// Check checks that the failure messages in t.Errorf follow the format expected.
func (c IdentifyFunction) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	for _, testBlock := range testFunc.TestPartBlocks() {
		if testBlock.IsRecommendedFailureMessage() {
			continue
		}

		diag := analysis.Diagnostic{
			Pos:      testBlock.TErrorCallExpr().CallExpr().Pos(),
			End:      testBlock.TErrorCallExpr().CallExpr().End(),
			Category: c.category,
			Message:  testBlock.ExpectedFailureMessage(),
			URL:      "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#identify-the-function",
		}
		pass.Report(diag)
	}
}
