// Package analyzer contains the analyzer with the business logic of this linter.
package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func New() *analysis.Analyzer {
	f := testcommentslint{}

	a := &analysis.Analyzer{
		Name:     "testcommentslint",
		Doc:      "checks test follow standards",
		URL:      "https://github.com/manuelarte/testcommentslint",
		Run:      f.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return a
}

type testcommentslint struct{}

func (l testcommentslint) run(pass *analysis.Pass) (any, error) {
	insp, found := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !found {
		//nolint:nilnil // impossible case.
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		//nolint:errcheck // impossible case.
		funcDecl := n.(*ast.FuncDecl)
		// Only process _test.go files
		if !strings.HasSuffix(pass.Fset.File(funcDecl.Pos()).Name(), "_test.go") {
			return
		}

		l.analyzeTestFunction(pass, funcDecl)
	})

	//nolint:nilnil //any, error
	return nil, nil
}

func (l testcommentslint) analyzeTestFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	isTest, testVar := isTestFunction(funcDecl)
	if !isTest {
		return
	}

	fmt.Println(testVar)
}
