package model

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewTableDrivenTestInfoIsTableDrivenTest(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		content   string
		wantBlock func(*ast.FuncDecl) *ast.BlockStmt
	}{
		"map non-inline table driven test": {
			content: `
package main

func TestExample(t *testing.T) {
  tests := map[string]struct {
    input  string
    want   int
  }{
    "example": {
      input:  "1",
      want:   1,
    },
  }
  for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parse(tc.input)
			if got != tc.want {
				t.Errorf("parse got %v, want %v", got, tc.want)
			}
		})
	}
}
			`[1:],
			wantBlock: func(funcDecl *ast.FuncDecl) *ast.BlockStmt {
				//nolint:lll
				return funcDecl.Body.List[1].(*ast.RangeStmt).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Args[1].(*ast.FuncLit).Body
			},
		},
		"map inline table driven test": {
			content: `
package main

func TestExample(t *testing.T) {
	for name, test := range map[string]struct {
		in int
		out int
	} {
		"test1": {
			in: 1,
			out: 1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := abs(test.in)
			if got != test.out {
				t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
			}
		})
	}
}
			`[1:],
			wantBlock: func(funcDecl *ast.FuncDecl) *ast.BlockStmt {
				//nolint:lll
				return funcDecl.Body.List[0].(*ast.RangeStmt).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Args[1].(*ast.FuncLit).Body
			},
		},
		"slice non-inlined table driven test": {
			content: `
package main

func TestExample(t *testing.T) {
  tests := []struct {
    desc   string
    input  string
    want   int
  }{
    {
      desc:   "example",
      input:  "1",
      want:   1,
    },
  }
  for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := parse(tc.input)
			if got != tc.want {
				t.Errorf("parse got %v, want %v", got, tc.want)
			}
		})
	}
}
			`[1:],
			wantBlock: func(funcDecl *ast.FuncDecl) *ast.BlockStmt {
				//nolint:lll
				return funcDecl.Body.List[1].(*ast.RangeStmt).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Args[1].(*ast.FuncLit).Body
			},
		},
		"slice inlined table driven test": {
			content: `
package main

func TestExample(t *testing.T) {
  for _, test := range []struct {
		name string
		in int
		out int
	} {
    	{
			name: "test1",
			in: 1,
			out: 1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := abs(test.in)
			if got != test.out {
				t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
			}
		})
	}
}
			`[1:],
			wantBlock: func(funcDecl *ast.FuncDecl) *ast.BlockStmt {
				//nolint:lll
				return funcDecl.Body.List[0].(*ast.RangeStmt).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Args[1].(*ast.FuncLit).Body
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()

			node, err := parser.ParseFile(fset, "test.go", tc.content, parser.ParseComments)
			if err != nil {
				t.Fatalf("error parsing file: %v", err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					got := newTableDrivenInfo("t", funcDecl)

					gotBlock := got.Block
					if tc.wantBlock != nil && !cmp.Equal(gotBlock, tc.wantBlock(funcDecl)) {
						t.Errorf("IsTableDrivenTest() mismatch (-want +got):\n%s", cmp.Diff(tc.wantBlock(funcDecl), gotBlock))
					}
				}

				return true
			})
		})
	}
}

func TestNewTableDrivenTestInfoNoTableDrivenTest(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		content string
	}{
		"no table driven test": {
			content: `
package main

func TestExample(t *testing.T) {
  got := parse("1")
  if got != tc.want {
    t.Errorf("parse got %v, want %v", got, tc.want)
  }
}
			`[1:],
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()

			node, err := parser.ParseFile(fset, "test.go", tc.content, parser.ParseComments)
			if err != nil {
				t.Fatalf("error parsing file: %v", err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					got := newTableDrivenInfo("t", funcDecl)
					if got != nil {
						t.Errorf("newTableDrivenInfo() = %v, want nil", got)
					}
				}

				return true
			})
		})
	}
}
