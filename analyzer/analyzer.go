// Package analyzer contains the analyzer with the business logic of this linter.
package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/manuelarte/testcommentslint/analyzer/checks"
	"github.com/manuelarte/testcommentslint/analyzer/model"
)

const (
	EqualityComparisonCheckName       = "equality-comparison"
	GotBeforeWantCheck                = "got-before-want"
	IdentifyTheFunctionCHeck          = "identify-function"
	TableDrivenFormatCheckTypeName    = "table-driven-format.type"
	TableDrivenFormatCheckInlinedName = "table-driven-format.inlined"
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
	a.Flags.BoolVar(&l.gotBeforeWant, GotBeforeWantCheck, true,
		"Check that output the actual value that the function returned before printing the value that was expected.")
	a.Flags.BoolVar(&l.identifyFunction, IdentifyTheFunctionCHeck, true,
		"Check that the failure messages in t.Errorf contains the function name.")
	a.Flags.StringVar(&l.tableDrivenFormat.formatType, TableDrivenFormatCheckTypeName, "",
		"Check that the table-driven tests are either Map or Slice.")
	a.Flags.BoolVar(&l.tableDrivenFormat.inlined, TableDrivenFormatCheckInlinedName, false,
		"Check that the table-driven tests are either inline or declared before.")

	return a
}

type (
	testcommentslint struct {
		equalityComparison bool
		gotBeforeWant      bool
		identifyFunction   bool
		tableDrivenFormat  tableDrivenFormat
	}
	tableDrivenFormat struct {
		formatType string
		inlined    bool
	}
)

func (t tableDrivenFormat) getTableDrivenFormatPredicate() checks.TableDrivenFormatPredicate {
	f := checks.TableDrivenFormatType(t.formatType)
	if f != checks.Map && f != checks.Slice {
		return checks.AlwaysValid()
	}

	pred, _ := checks.OfTypeAndInline(f, t.inlined)

	return pred
}

//nolint:gocognit // refactor later
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

	tbfCheck, err := checks.NewTableDrivenFormat(l.tableDrivenFormat.getTableDrivenFormatPredicate())
	if err != nil {
		return nil, fmt.Errorf("error creating table driven format checker: %w", err)
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

			tbfCheck.Check(pass, testFunc)

			if l.equalityComparison {
				checks.NewEqualityComparison().Check(pass, testFunc)
			}

			if l.identifyFunction {
				checks.NewFailureMessage().Check(pass, testFunc)
			}
		}
	})

	//nolint:nilnil //any, error
	return nil, nil
}
