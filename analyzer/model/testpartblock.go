package model

import (
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	"github.com/manuelarte/testcommentslint/analyzer/slicesutils"
)

// TestPartBlock is a struct that holds the typical testing block like:
// got := myFunction(in)	<- testedFunc
// if got != want { 		<- ifComparing
//
//	  t.Errorf(...)			<- tErrorfCallExpr
//	}.
type TestPartBlock struct {
	importGroup ImportGroup

	// testedFunc contain the actual call to the function tested.
	testedFunc TestedCallExpr

	// ifComparing contains the if statement that leads to t.Errorf or t.Fatalf.
	ifComparing IfComparing

	// tErrorCallExpr contains the call to t.Errorf or t.Fatalf and its parameters.
	tErrorCallExpr TErrorfCallExpr
}

func NewTestPartBlock(
	importGroup ImportGroup,
	testVar string,
	prev ast.Stmt,
	ifStmt *ast.IfStmt,
) (TestPartBlock, bool) {
	testedFunc, isTestedFunc := NewTestedCallExpr(prev)
	if !isTestedFunc {
		return TestPartBlock{}, false
	}

	ifComparing, isComparingIfStmt := NewIfComparingResult(importGroup, testedFunc.Params(), ifStmt)
	if !isComparingIfStmt {
		return TestPartBlock{}, false
	}

	teCallExpr, istErrorf := NewTErrorfCallExpr(testVar, ifStmt.Body)
	if !istErrorf {
		return TestPartBlock{}, false
	}

	return TestPartBlock{
		importGroup:    importGroup,
		testedFunc:     testedFunc,
		ifComparing:    ifComparing,
		tErrorCallExpr: teCallExpr,
	}, true
}

func (t TestPartBlock) TErrorCallExpr() TErrorfCallExpr {
	return t.tErrorCallExpr
}

func (t TestPartBlock) ExpectedFailureMessage() string {
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

	switch t.ifComparing.(type) {
	case ComparingParamsIfStmt:
		return fmt.Sprintf("Prefer \"%s, want %%v\" format for this failure message", funcFailurePart)
	case DiffIfStmt:
		return fmt.Sprintf("Prefer \"%s mismatch (-want +got):\\n%%s\" format for this failure message", funcFailurePart)
	}

	return ""
}

// IsRecommendedFailureMessage returns whether the failure message honors the expected format for comparison used.
func (t TestPartBlock) IsRecommendedFailureMessage() bool {
	currentFailureMessage := t.tErrorCallExpr.FailureMessage()

	unquoted, err := strconv.Unquote(currentFailureMessage)
	if err != nil {
		unquoted = currentFailureMessage
	}

	switch t.ifComparing.(type) {
	case ComparingParamsIfStmt:
		funName := t.testedFunc.FunctionName()
		quotedFunName := regexp.QuoteMeta(funName)
		pattern := fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, quotedFunName)

		matched, _ := regexp.MatchString(pattern, unquoted)
		if matched {
			return true
		}

		if selExpr, ok := t.testedFunc.CallExpr().Fun.(*ast.SelectorExpr); ok {
			funName = selExpr.Sel.Name
			quotedFunName = regexp.QuoteMeta(funName)
			pattern = fmt.Sprintf(`^%s(?:|\(.*\)) = %%[^,]+, want %%[^,]+$`, quotedFunName)
			matched, _ = regexp.MatchString(pattern, unquoted)
			return matched
		}

		return false
	case DiffIfStmt:
		pattern := `(?:-want \+got|\(-want \+got\)):\n%s$`
		matched, _ := regexp.MatchString(pattern, unquoted)
		if matched {
			return true
		}

		funName := t.testedFunc.FunctionName()
		quotedFunName := regexp.QuoteMeta(funName)
		pattern = fmt.Sprintf(`^%s mismatch (?:-want \+got|\(-want \+got\)):\n%%s$`, quotedFunName)
		matched, _ = regexp.MatchString(pattern, unquoted)
		if matched {
			return true
		}

		if selExpr, ok := t.testedFunc.CallExpr().Fun.(*ast.SelectorExpr); ok {
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
