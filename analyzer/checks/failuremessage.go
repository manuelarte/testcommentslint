package checks

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// FailureMessage check that the failure messages in t.Errorf follow the format expected.
// The format expected can be as the following:
//   - when the condition is reflect.DeepEqual, cmp.Equal or got != want: "YourFunction(%v) = %v, want %v"
//   - when the condition is cmp.Diff: YourFunction() mismatch (-want +got):\n%s
type FailureMessage struct {
	category string
}

// NewFailureMessage creates a new FailureMessage.
func NewFailureMessage() FailureMessage {
	return FailureMessage{
		category: "Failure Message",
	}
}

func (c FailureMessage) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	blStmt := testFunc.GetActualTestBlockStmt()
	testVar := testFunc.GetTestVar()

	var stmts []ast.Stmt
	if blStmt != nil {
		stmts = blStmt.List
	}

	for i, stmt := range stmts {
		if ifStmt, ok := stmt.(*ast.IfStmt); ok {
			// Check the condition to see if it's a comparison like got != want
			if i == 0 {
				continue
			}

			testBlock, isTestBlock := newTestFuncBlock(testVar, stmts[i-1], ifStmt)
			if !isTestBlock {
				continue
			}

			diag := analysis.Diagnostic{
				Pos:      testBlock.ifStmt.errorCallExpr.callExpr.Pos(),
				End:      testBlock.ifStmt.errorCallExpr.callExpr.End(),
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
	// got := myFunction(in)
	// if got != want {
	//   t.Errorf(...)
	// }.
	testFuncBlock struct {
		// testedFunc contain the actual call to the function tested.
		testedFunc testFuncStmt
		// contain the got != want condition and the t.Errorf call.
		ifStmt gotWantIfStmt
	}

	// testFuncStmt contains the actual call to the function tested.
	testFuncStmt struct {
		callExpr *ast.CallExpr

		params []*ast.Ident
	}

	// gotWantIfStatement struct holding an if statement that contains a comparison of got and want.
	// Conditions allowed:
	// - got != want
	// - reflect.DeepEqual
	// - cmp.Equal
	// - cmp.Diff
	// and inside the if statement there is a call to t.Errorf.
	gotWantIfStmt struct {
		// original ifStmt
		ifStmt *ast.IfStmt

		params        []*ast.Ident
		errorCallExpr tErrorfCallExpr
	}

	// tErrorfCallExpr contains the call to t.Errorf and its parameters.
	tErrorfCallExpr struct {
		callExpr *ast.CallExpr

		params []*ast.Ident
	}
)

func newTestFuncBlock(testVar string, prev ast.Stmt, ifStmt *ast.IfStmt) (testFuncBlock, bool) {
	gwStmt, isGotWant := newGotWantIfStmt(testVar, ifStmt)
	if !isGotWant {
		return testFuncBlock{}, false
	}

	testedFunc, isTestedFunc := newTestedFuncExpr(prev)
	if !isTestedFunc {
		return testFuncBlock{}, false
	}

	return testFuncBlock{
		testedFunc: testedFunc,
		ifStmt:     gwStmt,
	}, true
}

func (t testFuncBlock) getFunctionName() string {
	return t.testedFunc.getFunctionName()
}

//nolint:unused // to be used later
func (t testFuncBlock) getGotName() string {
	for _, param := range t.testedFunc.params {
		if param.Name == "_" {
			continue
		}

		for _, ifParam := range t.ifStmt.params {
			if param.Name == ifParam.Name {
				return param.Name
			}
		}
	}

	// impossible case
	return "got"
}

func (t testFuncBlock) expectedFailureMessage() string {
	in := make([]string, len(t.testedFunc.callExpr.Args))
	for i := range in {
		in[i] = "%v"
	}

	out := make([]string, len(t.testedFunc.params))
	for i := range out {
		if t.testedFunc.params[i].Name == "_" {
			out[i] = "_"
		} else {
			out[i] = "%v"
		}
	}

	funcFailurePart := fmt.Sprintf("%s(%s) = %s", t.getFunctionName(), strings.Join(in, ", "), strings.Join(out, ", "))

	return fmt.Sprintf("Prefer \"%s, want %%v\" format for failure message", funcFailurePart)
}

// newTestedFuncExpr creates a testedFuncStmt after checking that the stmt is a typical function call.
func newTestedFuncExpr(stmt ast.Stmt) (testFuncStmt, bool) {
	var callExpr *ast.CallExpr

	params := make([]*ast.Ident, 0)

	assignStmt, isAssignStmt := stmt.(*ast.AssignStmt)
	if !isAssignStmt {
		return testFuncStmt{}, false
	}

	if len(assignStmt.Rhs) != 1 {
		return testFuncStmt{}, false
	}

	for _, expr := range assignStmt.Lhs {
		ident, ok := expr.(*ast.Ident)
		if !ok {
			return testFuncStmt{}, false
		}

		params = append(params, ident)
	}

	ce, ok := assignStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return testFuncStmt{}, false
	}

	callExpr = ce

	return testFuncStmt{
		callExpr: callExpr,
		params:   params,
	}, true
}

func (t testFuncStmt) getFunctionName() string {
	switch fn := t.callExpr.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		if ident, ok := fn.X.(*ast.Ident); ok {
			return ident.Name + "." + fn.Sel.Name
		}
	}

	return ""
}

// newGotWantIfStmt creates a new gotWantIfStmt.
// only if the condition applies.
func newGotWantIfStmt(testVar string, ifStmt *ast.IfStmt) (gotWantIfStmt, bool) {
	if ifStmt == nil || ifStmt.Body == nil {
		return gotWantIfStmt{}, false
	}

	if len(ifStmt.Body.List) != 1 {
		return gotWantIfStmt{}, false
	}

	teCallExpr, istErrorf := newTErrorfCallExpr(testVar, ifStmt.Body.List[0])
	if !istErrorf {
		return gotWantIfStmt{}, false
	}

	params := make([]*ast.Ident, 2)

	switch node := ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		// check ident1 != ident2 and both are used in the failure message t.Errorf
		if node.Op.String() != "!=" {
			return gotWantIfStmt{}, false
		}

		xIdent, isXIdent := isNotBlankIdent(node.X)

		yIdent, isYIdent := isNotBlankIdent(node.Y)
		if !isXIdent || !isYIdent {
			return gotWantIfStmt{}, false
		}

		params[0] = xIdent
		params[1] = yIdent
	default:
		return gotWantIfStmt{}, false
	}

	// check params match
	if params[0].Name != teCallExpr.params[0].Name && params[0].Name != teCallExpr.params[1].Name {
		return gotWantIfStmt{}, false
	}

	if params[1].Name != teCallExpr.params[0].Name && params[1].Name != teCallExpr.params[1].Name {
		return gotWantIfStmt{}, false
	}

	return gotWantIfStmt{
		ifStmt:        ifStmt,
		params:        params,
		errorCallExpr: teCallExpr,
	}, true
}

// newTErrorfCallExpr creates a tErrorfCallExpr after checking that the stmt is a call to t.Errorf.
func newTErrorfCallExpr(testVar string, stmt ast.Stmt) (tErrorfCallExpr, bool) {
	exprStmt, isExprStmt := stmt.(*ast.ExprStmt)
	if !isExprStmt {
		return tErrorfCallExpr{}, false
	}

	callExpr, isCallExpr := exprStmt.X.(*ast.CallExpr)
	if !isCallExpr {
		return tErrorfCallExpr{}, false
	}

	selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	if !isSelectorExpr {
		return tErrorfCallExpr{}, false
	}

	ident, isIdent := selectorExpr.X.(*ast.Ident)
	if !isIdent || ident.Name != testVar || selectorExpr.Sel.Name != "Errorf" {
		return tErrorfCallExpr{}, false
	}

	if len(callExpr.Args) != 3 {
		return tErrorfCallExpr{}, false
	}

	basicLit, isBasicLit := callExpr.Args[0].(*ast.BasicLit)
	if !isBasicLit || basicLit.Kind.String() != "STRING" {
		return tErrorfCallExpr{}, false
	}

	firstIdent, isFirstIdent := callExpr.Args[1].(*ast.Ident)

	secondIdent, isSecondIdent := callExpr.Args[2].(*ast.Ident)
	if !isFirstIdent || !isSecondIdent {
		return tErrorfCallExpr{}, false
	}

	return tErrorfCallExpr{
		callExpr: callExpr,
		params:   []*ast.Ident{firstIdent, secondIdent},
	}, true
}

func isNotBlankIdent(expr ast.Expr) (*ast.Ident, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return nil, false
	}

	if ident.Name == "_" {
		return nil, false
	}

	return ident, true
}
