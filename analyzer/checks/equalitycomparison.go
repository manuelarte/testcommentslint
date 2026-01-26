package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// EqualityComparisonCheck checks that reflect.DeepEqual can be replaced by newer cmp.Equal.
type EqualityComparisonCheck struct {
	// testVar is the identifier given for the testing.T parameter in the test function.
	testVar string
	// reflectImportName is the import name of the reflect package.
	reflectImportName string
}

// NewEqualityComparisonCheck creates a new EqualityComparisonCheck with the given reflect import name.
func NewEqualityComparisonCheck(testVar, reflectImportName string) *EqualityComparisonCheck {
	if testVar == "" {
		testVar = "t"
	}

	if reflectImportName == "" {
		reflectImportName = "reflect"
	}

	return &EqualityComparisonCheck{
		testVar:           testVar,
		reflectImportName: reflectImportName,
	}
}

func (c EqualityComparisonCheck) Check(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// TODO
	var stmts []ast.Stmt
	if funcDecl.Body != nil {
		stmts = funcDecl.Body.List
	}
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.CallExpr:
				return true
			}
			return false
		})
	}
}
