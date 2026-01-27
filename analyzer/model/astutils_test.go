package model

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestIsTableDrivenTest(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		content string
		want    bool
	}{
		"map table driven test": {
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
			want: true,
		},
		"slice table driven test": {
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
			want: true,
		},
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
			want: false,
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
					got, _ := isTableDrivenTest(funcDecl)
					if got != tc.want {
						t.Errorf("IsTableDrivenTest() got %v, want %v", got, tc.want)
					}
				}

				return true
			})
		})
	}
}
