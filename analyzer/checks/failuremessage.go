package checks

import (
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
	"github.com/manuelarte/testcommentslint/analyzer/slicesutils"
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
			testBlock, isTestBlock := newTestFuncBlock(testFunc.ImportGroup(), testVar, stmts[i-1], ifStmt)
			if !isTestBlock {
				continue
			}

			if testBlock.isRecommendedFailureMessage() {
				continue
			}
			diag := analysis.Diagnostic{
				Pos:      testBlock.tErrorCallExpr.CallExpr().Pos(),
				End:      testBlock.tErrorCallExpr.CallExpr().End(),
				Category: c.category,
				Message:  testBlock.expectedFailureMessage(),
				URL:      "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#failure-message",
			}
			pass.Report(diag)
		}
	}
}

// Auxiliary structs to facilitate the business logic.
type (
	// testFuncBlock is a struct that holds the typical testing block like:
	// got := myFunction(in)	<- testedFunc
	// if got != want { 		<- ifComparingResult
	//   t.Errorf(...)			<- tErrorfCallExpr
	// }.
	testFuncBlock struct {
		importGroup model.ImportGroup

		// testedFunc contain the actual call to the function tested.
		testedFunc model.TestedCallExpr

		// ifComparing contains the if statement that leads to t.Errorf or t.Fatalf.
		ifComparing model.IfComparing

		// tErrorCallExpr contains the call to t.Errorf or t.Fatalf and its parameters.
		tErrorCallExpr model.TErrorfCallExpr
	}
)

func newTestFuncBlock(
	importGroup model.ImportGroup,
	testVar string,
	prev ast.Stmt,
	ifStmt *ast.IfStmt,
) (testFuncBlock, bool) {
	testedFunc, isTestedFunc := model.NewTestedCallExpr(prev)
	if !isTestedFunc {
		return testFuncBlock{}, false
	}

	ifComparing, isComparingIfStmt := model.NewIfComparingResult(importGroup, testedFunc.Params(), ifStmt)
	if !isComparingIfStmt {
		return testFuncBlock{}, false
	}

	teCallExpr, istErrorf := model.NewTErrorfCallExpr(testVar, ifStmt.Body)
	if !istErrorf {
		return testFuncBlock{}, false
	}

	return testFuncBlock{
		importGroup:    importGroup,
		testedFunc:     testedFunc,
		ifComparing:    ifComparing,
		tErrorCallExpr: teCallExpr,
	}, true
}

// isRecommendedFailureMessage returns whether the failure message honors the expected format for comparison used.
func (t testFuncBlock) isRecommendedFailureMessage() bool {
	currentFailureMessage := t.tErrorCallExpr.FailureMessage()

	unquoted, err := strconv.Unquote(currentFailureMessage)
	if err != nil {
		unquoted = currentFailureMessage
	}

	switch t.ifComparing.(type) {
	case model.ComparingParamsIfStmt:
		funName := t.testedFunc.FunctionName()
		quotedFunName := regexp.QuoteMeta(funName)
		pattern := fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, quotedFunName)

		matched, _ := regexp.MatchString(pattern, unquoted)

		return matched
	case model.DiffIfStmt:
		pattern := `(?:-want \+got|\(-want \+got\)):\n%s$`
		matched, _ := regexp.MatchString(pattern, unquoted)

		return matched
	}

	return true
}

func (t testFuncBlock) expectedFailureMessage() string {
	if _, ok := t.ifComparing.(model.DiffIfStmt); ok {
		return "Prefer \"diff -want +got:\\n%s\" format for this failure message"
	}

	in := strings.Join(slicesutils.Map(t.testedFunc.CallExpr().Args, func(in ast.Expr) string {
		return "%v"
	}), ", ")

	out := strings.Join(slicesutils.Map(t.testedFunc.Params(), func(in *ast.Ident) string {
		if in.Name == "_" {
			return "_"
		}

		return "%v"
	}), ", ")

	funcFailurePart := fmt.Sprintf("%s(%s) = %s", t.testedFunc.FunctionName(), in, out)

	return fmt.Sprintf("Prefer \"%s, want %%v\" format for this failure message", funcFailurePart)
}
