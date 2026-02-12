package model

import "go/ast"

// TErrorfCallExpr contains the call to t.Errorf and its parameters.
type TErrorfCallExpr struct {
	callExpr       *ast.CallExpr
	failureMessage string
	params         []*ast.Ident
}

// NewTErrorfCallExpr creates a tErrorfCallExpr after checking that the stmt is a call to t.Errorf.
func NewTErrorfCallExpr(testVar string, blStmts *ast.BlockStmt) (TErrorfCallExpr, bool) {
	if blStmts == nil {
		return TErrorfCallExpr{}, false
	}

	if len(blStmts.List) != 1 {
		return TErrorfCallExpr{}, false
	}

	stmt := blStmts.List[0]

	exprStmt, isExprStmt := stmt.(*ast.ExprStmt)
	if !isExprStmt {
		return TErrorfCallExpr{}, false
	}

	callExpr, isCallExpr := exprStmt.X.(*ast.CallExpr)
	if !isCallExpr {
		return TErrorfCallExpr{}, false
	}

	selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	if !isSelectorExpr {
		return TErrorfCallExpr{}, false
	}

	ident, isIdent := selectorExpr.X.(*ast.Ident)
	if !isIdent || ident.Name != testVar || selectorExpr.Sel.Name != "Errorf" {
		return TErrorfCallExpr{}, false
	}

	if len(callExpr.Args) < 2 {
		return TErrorfCallExpr{}, false
	}

	basicLit, isBasicLit := callExpr.Args[0].(*ast.BasicLit)
	if !isBasicLit || basicLit.Kind.String() != "STRING" {
		return TErrorfCallExpr{}, false
	}

	params := make([]*ast.Ident, 0)

	for i := 1; len(callExpr.Args) > i; i++ {
		// TODO here we need to accept also selectorExpr that comes from the table driven tests
		paramIdent, isParamIdent := callExpr.Args[i].(*ast.Ident)

		_, isTestSelectorExpr := callExpr.Args[i].(*ast.SelectorExpr)
		if !isParamIdent && !isTestSelectorExpr {
			return TErrorfCallExpr{}, false
		}

		params = append(params, paramIdent)
	}

	return TErrorfCallExpr{
		callExpr:       callExpr,
		failureMessage: basicLit.Value,
		params:         params,
	}, true
}

func (t TErrorfCallExpr) CallExpr() *ast.CallExpr {
	return t.callExpr
}

func (t TErrorfCallExpr) FailureMessage() string {
	return t.failureMessage
}

func (t TErrorfCallExpr) Params() []*ast.Ident {
	return t.params
}
