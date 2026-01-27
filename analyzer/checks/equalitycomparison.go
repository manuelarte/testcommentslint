package checks

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// EqualityComparisonCheck checks that reflect.DeepEqual can be replaced by newer cmp.Equal.
type EqualityComparisonCheck struct{}

// NewEqualityComparisonCheck creates a new EqualityComparisonCheck.
func NewEqualityComparisonCheck() EqualityComparisonCheck {
	return EqualityComparisonCheck{}
}

func (c EqualityComparisonCheck) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	// TODO
	reflectImportName, ok := testFunc.ReflectImportName()
	if !ok {
		return
	}

	blStmt := testFunc.GetActualTestBlockStmt()

	var stmts []ast.Stmt
	if blStmt != nil {
		stmts = blStmt.List
	}

	for _, stmt := range stmts {
		ast.Inspect(stmt, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.CallExpr:
				fmt.Println(reflectImportName)

				return true
			}

			return false
		})
	}
}
