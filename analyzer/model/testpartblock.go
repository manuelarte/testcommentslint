package model

import (
	"go/ast"
)

// TestPartBlock is a struct that holds the typical testing block like:
//
//	 	got := MyFunction(in)
//		if got != want {
//		  t.Errorf(...)
//		}
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

func (t TestPartBlock) TestedFunc() TestedCallExpr {
	return t.testedFunc
}

func (t TestPartBlock) IfComparing() IfComparing {
	return t.ifComparing
}

func (t TestPartBlock) TErrorCallExpr() TErrorfCallExpr {
	return t.tErrorCallExpr
}
