package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
	//nolint:nilnil //any, error
	return nil, nil
}
