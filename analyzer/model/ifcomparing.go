package model

import (
	"go/ast"
	"go/token"
)

var (
	_ IfComparing = new(ComparingParamsIfStmt)
	_ IfComparing = new(DiffIfStmt)
)

type (
	// IfComparing interface that contains the if statement that leads to t.Errorf or t.Fatalf.
	IfComparing interface {
		IfStmt() *ast.IfStmt
	}

	// ComparingParamsIfStmt if statement that leads to t.Errorf or t.Fatalf and that indicates:
	// 1. if param1 != param2
	// 2. if !reflect.DeepEqual(param1, param2)
	// 3. if !cmp.Equal(param1, param2).
	ComparingParamsIfStmt struct {
		ifStmt *ast.IfStmt

		// got param that identifies the result of the tested function
		got *ast.Ident
		// want expr that identifies the expected result, it can be *ast.Ident or *ast.SelectExpr
		want ast.Expr
	}

	// DiffIfStmt if statement that leads to t.Errorf or t.Fatalf and that indicates:
	// if diff := cmp.Diff(param1, param2); diff != "".
	DiffIfStmt struct {
		ifStmt *ast.IfStmt
	}
)

// NewIfComparingResult creates a new IfComparingResult based on the if condition.
func NewIfComparingResult(
	importGroup ImportGroup,
	testedFunctionParams []*ast.Ident,
	ifStmt *ast.IfStmt,
) (IfComparing, bool) {
	if ifStmt == nil || ifStmt.Body == nil {
		return nil, false
	}

	if len(ifStmt.Body.List) != 1 {
		return nil, false
	}

	if ifStmt.Init == nil {
		// case got != equal and !reflect.DeepEqual or !cmp.Equal
		got, want, ok := getGotWantParams(importGroup, testedFunctionParams, ifStmt.Cond)
		if !ok {
			return nil, false
		}

		return ComparingParamsIfStmt{
			ifStmt: ifStmt,
			got:    got,
			want:   want,
		}, true
	}

	// case cmp.Diff
	ok := isDiffParamIfStmt(importGroup, ifStmt)
	if !ok {
		return nil, false
	}

	return DiffIfStmt{
		ifStmt: ifStmt,
	}, true
}

func (c ComparingParamsIfStmt) IfStmt() *ast.IfStmt {
	return c.ifStmt
}

func (c ComparingParamsIfStmt) Got() *ast.Ident {
	return c.got
}

func (c ComparingParamsIfStmt) Want() ast.Expr {
	return c.want
}

func (d DiffIfStmt) IfStmt() *ast.IfStmt {
	return d.ifStmt
}

//nolint:gocognit // refactor later
func isDiffParamIfStmt(importGroup ImportGroup, ifStmt *ast.IfStmt) bool {
	var diffParam *ast.Ident

	switch node := ifStmt.Init.(type) {
	case *ast.AssignStmt:
		if len(node.Lhs) != 1 {
			return false
		}

		ident, ok := node.Lhs[0].(*ast.Ident)
		if !ok {
			return false
		}

		if len(node.Rhs) != 1 {
			return false
		}

		callExpr, ok := node.Rhs[0].(*ast.CallExpr)
		if !ok {
			return false
		}

		if len(callExpr.Args) != 2 {
			return false
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return false
		}

		goCmpImportAlias, _ := importGroup.GoCmpImportName()

		if !isGoCmpDiff(goCmpImportAlias, selectorExpr) {
			return false
		}

		diffParam = ident
	default:
		return false
	}

	switch node := ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		// check "ident1 != ident2" and both are used in the failure message `t.Errorf`.
		if node.Op != token.NEQ {
			return false
		}

		xIdent, isXIdent := isNotBlankIdent(node.X)
		if !isXIdent {
			return false
		}

		if basicLit, isBasicLit := node.Y.(*ast.BasicLit); !isBasicLit || basicLit.Value != "\"\"" {
			return false
		}

		if xIdent.Name != diffParam.Name {
			return false
		}
	default:
		return false
	}

	return true
}

//nolint:gocognit // refactor later
func getGotWantParams(
	importGroup ImportGroup,
	testedFunctionParams []*ast.Ident,
	cond ast.Expr,
) (*ast.Ident, ast.Expr, bool) {
	switch node := cond.(type) {
	case *ast.BinaryExpr:
		// check "ident1 != ident2".
		if node.Op != token.NEQ {
			return nil, nil, false
		}

		xIdent, isXIdent := isNotBlankIdent(node.X)

		yIdent, isYIdent := isNotBlankIdent(node.Y)
		if isYIdent && yIdent.Name == "nil" {
			return nil, nil, false
		}

		for _, p := range testedFunctionParams {
			if isXIdent && p.Name == xIdent.Name {
				got := xIdent
				want := node.Y

				return got, want, true
			}

			if isYIdent && p.Name == yIdent.Name {
				got := yIdent
				want := node.X

				return got, want, true
			}
		}
	case *ast.UnaryExpr:
		// check either `!reflect.DeepEqual` or `!cmp.Equal`.
		if node.Op != token.NOT {
			return nil, nil, false
		}

		callExpr, ok := node.X.(*ast.CallExpr)
		if !ok {
			return nil, nil, false
		}

		if len(callExpr.Args) != 2 {
			return nil, nil, false
		}

		xIdent, isXIdent := isNotBlankIdent(callExpr.Args[0])
		yIdent, isYIdent := isNotBlankIdent(callExpr.Args[1])

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil, nil, false
		}

		goCmpImportAlias, _ := importGroup.GoCmpImportName()
		reflectImportAlias, _ := importGroup.ReflectImportName()

		if !isGoCmpEqual(goCmpImportAlias, selectorExpr) && !isReflectEqual(reflectImportAlias, selectorExpr) {
			return nil, nil, false
		}

		for _, p := range testedFunctionParams {
			if isXIdent && p.Name == xIdent.Name {
				got := xIdent
				want := callExpr.Args[1]

				return got, want, true
			}

			if isYIdent && p.Name == yIdent.Name {
				got := yIdent
				want := callExpr.Args[0]

				return got, want, true
			}
		}
	default:
		return nil, nil, false
	}

	return nil, nil, false
}
