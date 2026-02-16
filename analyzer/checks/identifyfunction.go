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
		if isRecommendedFailureMessage(testBlock) {
			continue
		}

		diag := analysis.Diagnostic{
			Pos:      testBlock.TErrorCallExpr().CallExpr().Pos(),
			End:      testBlock.TErrorCallExpr().CallExpr().End(),
			Category: c.category,
			Message:  expectedFailureMessage(testBlock),
			URL:      "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#identify-the-function",
		}
		pass.Report(diag)
	}
}

// isRecommendedFailureMessage returns whether the failure message honors the expected format for comparison used.
// TODO: check that the function name is present, only that. Not that it contains got and want
func isRecommendedFailureMessage(t model.TestPartBlock) bool {
	currentFailureMessage := t.TErrorCallExpr().FailureMessage()

	unquoted, err := strconv.Unquote(currentFailureMessage)
	if err != nil {
		unquoted = currentFailureMessage
	}

	switch t.IfComparing().(type) {
	case model.ComparingParamsIfStmt:
		funName := t.TestedFunc().FunctionName()
		//isRecommendedGotWantFailureMessage(funName, unquoted)
		quotedFunName := regexp.QuoteMeta(funName)
		pattern := fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, quotedFunName)

		matched, _ := regexp.MatchString(pattern, unquoted)
		if matched {
			return true
		}

		if selExpr, ok := t.TestedFunc().CallExpr().Fun.(*ast.SelectorExpr); ok {
			funName = selExpr.Sel.Name
			quotedFunName = regexp.QuoteMeta(funName)
			pattern = fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, quotedFunName)
			matched, _ = regexp.MatchString(pattern, unquoted)

			return matched
		}

		return false
	case model.DiffIfStmt:
		pattern := `(?:-want \+got|\(-want \+got\)):\n%s$`

		matched, _ := regexp.MatchString(pattern, unquoted)
		if matched {
			return true
		}

		funName := t.TestedFunc().FunctionName()
		quotedFunName := regexp.QuoteMeta(funName)
		pattern = fmt.Sprintf(`^%s mismatch (?:-want \+got|\(-want \+got\)):\n%%s$`, quotedFunName)

		matched, _ = regexp.MatchString(pattern, unquoted)
		if matched {
			return true
		}

		if selExpr, ok := t.TestedFunc().CallExpr().Fun.(*ast.SelectorExpr); ok {
			funName = selExpr.Sel.Name
			quotedFunName = regexp.QuoteMeta(funName)
			pattern = fmt.Sprintf(`^%s mismatch (?:-want \+got|\(-want \+got\)):\n%%s$`, quotedFunName)
			matched, _ = regexp.MatchString(pattern, unquoted)

			return matched
		}

		return false
	}

	return true
}

func expectedFailureMessage(t model.TestPartBlock) string {
	in := strings.Join(slicesutils.Map(t.TestedFunc().CallExpr().Args, func(in ast.Expr) string {
		return "%v"
	}), ", ")

	out := strings.Join(slicesutils.Map(t.TestedFunc().Params(), func(in *ast.Ident) string {
		if in.Name == "_" {
			return "_"
		}

		return "%v"
	}), ", ")

	funcFailurePart := fmt.Sprintf("%s(%s) = %s", t.TestedFunc().FunctionName(), in, out)

	switch t.IfComparing().(type) {
	case model.ComparingParamsIfStmt:
		return fmt.Sprintf("Prefer \"%s, want %%v\" format for this failure message", funcFailurePart)
	case model.DiffIfStmt:
		return fmt.Sprintf("Prefer \"%s mismatch (-want +got):\\n%%s\" format for this failure message", funcFailurePart)
	}

	return ""
}

func isRecommendedGotWantFailureMessage(functionName, failureMessage string) bool {
	pattern := fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, functionName)

	matched, _ := regexp.MatchString(pattern, failureMessage)
	if matched {
		return true
	}

	if strings.Contains(functionName, ".") {
		splitted := strings.Split(functionName, ".")
		funName := splitted[len(splitted)-1]
		pattern = fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, funName)
		matched, _ = regexp.MatchString(pattern, functionName)

		return matched
	}

	return false
}

func isRecommendedDiffFailureMessage(functionName, failureMessage string) bool {
	pattern := `(?:-want \+got|\(-want \+got\)):\n%s$`

	matched, _ := regexp.MatchString(pattern, failureMessage)
	if matched {
		return true
	}

	pattern = fmt.Sprintf(`^%s mismatch (?:-want \+got|\(-want \+got\)):\n%%s$`, functionName)

	matched, _ = regexp.MatchString(pattern, functionName)
	if matched {
		return true
	}

	if strings.Contains(functionName, ".") {
		splitted := strings.Split(functionName, ".")
		funName := splitted[len(splitted)-1]
		pattern = fmt.Sprintf(`^%s mismatch (?:-want \+got|\(-want \+got\)):\n%%s$`, funName)
		matched, _ = regexp.MatchString(pattern, functionName)

		return matched
	}

	return false
}
