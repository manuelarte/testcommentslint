// Package analyzer contains the analyzer with the business logic of this linter.
package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/testcommentslint/analyzer/checks"
	"github.com/manuelarte/testcommentslint/analyzer/model"
)

const (
	EqualityComparisonCheckName = "equality-comparison"
)

func New() *analysis.Analyzer {
	l := testcommentslint{}

	a := &analysis.Analyzer{
		Name:     "testcommentslint",
		Doc:      "checks test follow standards",
		URL:      "https://github.com/manuelarte/testcommentslint",
		Run:      l.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	a.Flags.BoolVar(&l.equalityComparison, EqualityComparisonCheckName, true,
		"Checks reflect.DeepEqual can be replaced by newer cmp.Equal.")

	return a
}

type testcommentslint struct {
	equalityComparison bool
}

func (l *testcommentslint) run(pass *analysis.Pass) (any, error) {
	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
		(*ast.FuncDecl)(nil),
	}

	var reflectImport *model.ReflectImport

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Only process _test.go files
		if !strings.HasSuffix(pass.Fset.File(n.Pos()).Name(), "_test.go") {
			reflectImport = nil

			return
		}

		switch node := n.(type) {
		case *ast.ImportSpec:
			if rf, ok := model.NewReflectImport(node); ok {
				reflectImport = rf
			}
		case *ast.FuncDecl:
			testFunc, ok := model.NewTestFunction(reflectImport, node)
			if !ok {
				return
			}

			if l.equalityComparison {
				checks.NewEqualityComparisonCheck().Check(pass, testFunc)
			}
		}
	})

	//nolint:nilnil //any, error
	return nil, nil
}
