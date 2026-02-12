package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// FailureMessage check that the failure messages in t.Errorf follow the format expected.
// The format expected can be as the following:
//   - When the condition is `reflect.DeepEqual`, `cmp.Equal` or `got != want`: "YourFunction(%v) = %v, want %v"
//   - When the condition is `cmp.Diff`: YourFunction(%v) mismatch (-want +got):\n%s
//
// This checks blocks like the following:
// got := MyFunction(in)
//
//	if got != want {
//	  t.Errorf(...)
//	}
type FailureMessage struct {
	category string
}

// NewFailureMessage creates a new FailureMessage.
func NewFailureMessage() FailureMessage {
	return FailureMessage{
		category: "Failure Message",
	}
}

// Check checks that the failure messages in t.Errorf follow the format expected.
func (c FailureMessage) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	blStmt := testFunc.GetActualTestBlockStmt()
	testVar := testFunc.GetTestVar()

	var stmts []ast.Stmt
	if blStmt != nil {
		stmts = blStmt.List
	}

	for i, stmt := range stmts {
		if ifStmt, ok := stmt.(*ast.IfStmt); ok {
			if i == 0 {
				continue
			}

			// create an auxiliary testBlock struct that holds:
			// - if statement
			// - the t.Errorf call
			// - the tested function previous to the if statement
			testBlock, isTestBlock := model.NewTestPartBlock(testFunc.ImportGroup(), testVar, stmts[i-1], ifStmt)
			if !isTestBlock {
				continue
			}

			if testBlock.IsRecommendedFailureMessage() {
				continue
			}

			diag := analysis.Diagnostic{
				Pos:      testBlock.TErrorCallExpr().CallExpr().Pos(),
				End:      testBlock.TErrorCallExpr().CallExpr().End(),
				Category: c.category,
				Message:  testBlock.ExpectedFailureMessage(),
				URL:      "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#failure-message",
			}
			pass.Report(diag)
		}
	}
}
