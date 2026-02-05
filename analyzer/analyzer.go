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
	FailureMessageCheckName     = "failure-message"
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
	a.Flags.BoolVar(&l.failureMessage, FailureMessageCheckName, true,
		"Check that the failure messages in t.Errorf follow the format expected.")

	return a
}

type testcommentslint struct {
	equalityComparison bool
	failureMessage     bool
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

	var importGroup model.ImportGroup

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Only process _test.go files
		if !strings.HasSuffix(pass.Fset.File(n.Pos()).Name(), "_test.go") {
			importGroup = model.ImportGroup{}

			return
		}

		switch node := n.(type) {
		case *ast.ImportSpec:
			if model.IsReflectImport(node) {
				importGroup.Reflect = node
			}

			if model.IsGoCmpImport(node) {
				importGroup.GoCmp = node
			}
		case *ast.FuncDecl:
			testFunc, ok := model.NewTestFunction(importGroup, node)
			if !ok {
				return
			}

			if l.equalityComparison {
				checks.NewEqualityComparison().Check(pass, testFunc)
			}

			if l.failureMessage {
				checks.NewFailureMessage().Check(pass, testFunc)
			}
		}
	})

	//nolint:nilnil //any, error
	return nil, nil
}
