package model

import (
	"go/ast"
	"testing"
)

func TestFunctionName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testedCallExpr TestedCallExpr
		want           string
	}{
		"MyFunction()": {
			testedCallExpr: TestedCallExpr{
				callExpr: &ast.CallExpr{
					Fun: &ast.Ident{
						Name: "MyFunction",
					},
				},
			},
			want: "MyFunction",
		},
		"mystruct.MyFunction()": {
			testedCallExpr: TestedCallExpr{
				callExpr: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "mystruct",
						},
						Sel: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
			},
			want: "mystruct.MyFunction",
		},
		"test.mystruct.MyFunction()": {
			testedCallExpr: TestedCallExpr{
				callExpr: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "test"},
							Sel: &ast.Ident{Name: "mystruct"},
						},
						Sel: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
			},
			want: "test.mystruct.MyFunction",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := test.testedCallExpr.FunctionName()
			if got != test.want {
				t.Errorf("FunctionName() = %q, want %q", got, test.want)
			}
		})
	}
}
