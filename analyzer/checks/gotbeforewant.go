package checks

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// GotBeforeWant struct that test outputs should output the actual value that the function returned
// before printing the value that was expected.
type GotBeforeWant struct {
	category string
}

// NewGotBeforeWant creates a new GotBeforeWant.
func NewGotBeforeWant() GotBeforeWant {
	return GotBeforeWant{
		category: "Got Before Want",
	}
}

// Check test outputs should output the actual value that the function returned before printing
// the value that was expected.
func (c GotBeforeWant) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	for _, testBlock := range testFunc.TestPartBlocks() {
		ifComparing, ok := testBlock.IfComparing().(model.ComparingParamsIfStmt)
		if !ok {
			continue
		}

		var gotIndex, wantIndex int

		got := ifComparing.Got()
		want := ifComparing.Want()

		for i, arg := range testBlock.TErrorCallExpr().GetArgs() {
			if ident, isIdent := arg.(*ast.Ident); isIdent && ident.Name == got.Name {
				gotIndex = i

				continue
			}

			if isSameVariable(want, arg) {
				wantIndex = i

				continue
			}
		}

		if gotIndex < wantIndex {
			continue
		}

		diag := analysis.Diagnostic{
			Pos:      testBlock.TErrorCallExpr().CallExpr().Pos(),
			End:      testBlock.TErrorCallExpr().CallExpr().End(),
			Category: c.category,
			Message: "Test outputs should output the actual value that the function returned before " +
				"printing the value that was expected",
			URL: "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#got-before-want",
		}
		pass.Report(diag)
	}
}

func isSameVariable(a, b ast.Expr) bool {
	switch nodeA := a.(type) {
	case *ast.Ident:
		if nodeB, isIdent := b.(*ast.Ident); isIdent && nodeA.Name == nodeB.Name {
			return true
		}
	case *ast.SelectorExpr:
		if nodeB, isSelectorExpr := b.(*ast.SelectorExpr); isSelectorExpr {
			if nodeA.Sel.Name != nodeB.Sel.Name {
				return false
			}

			return isSameVariable(nodeA.X, nodeB.X)
		}
	}

	return false
}
