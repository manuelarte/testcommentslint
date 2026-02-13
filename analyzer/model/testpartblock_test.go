package model

import (
	"go/ast"
	"testing"
)

func TestIsRecommendedFailureMessage(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testPartBlock TestPartBlock
		want          bool
	}{
		"valid function with one parameter": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "YourFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "YourFunction(%v) = %v, want %v",
				},
				ifComparing: ComparingParamsIfStmt{},
			},
			want: true,
		},
		"valid function with zero parameter": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "YourFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "YourFunction() = %v, want %v",
				},
				ifComparing: ComparingParamsIfStmt{},
			},
			want: true,
		},
		"valid function, no parenthesis": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "YourFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "YourFunction = %v, want %v",
				},
				ifComparing: ComparingParamsIfStmt{},
			},
			want: true,
		},
		"different function name with zero parameter": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "YourFunction() = %v, want %v",
				},
				ifComparing: ComparingParamsIfStmt{},
			},
			want: false,
		},
		"got want, no function name": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "got %v, want %v",
				},
				ifComparing: ComparingParamsIfStmt{},
			},
			want: false,
		},
		"expected, actual, no function name": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "actual: %v, expected %v",
				},
				ifComparing: ComparingParamsIfStmt{},
			},
			want: false,
		},
		"diff, expecting -want +got:\\n%s": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "diff: %s",
				},
				ifComparing: DiffIfStmt{},
			},
			want: false,
		},
		"diff with function name, with (-want +got):\\n%s": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "MyFunction mismatch (-want +got):\n%s",
				},
				ifComparing: DiffIfStmt{},
			},
			want: true,
		},
		"diff with function name, with -want +got:\\n%s": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "MyFunction mismatch -want +got:\n%s",
				},
				ifComparing: DiffIfStmt{},
			},
			want: true,
		},
		"diff without function name, with -want +got:\\n%s": {
			testPartBlock: TestPartBlock{
				testedFunc: TestedCallExpr{
					callExpr: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "MyFunction",
						},
					},
				},
				tErrorCallExpr: TErrorfCallExpr{
					failureMessage: "diff -want +got:\n%s",
				},
				ifComparing: DiffIfStmt{},
			},
			want: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tc.testPartBlock.IsRecommendedFailureMessage()
			if got != tc.want {
				t.Errorf("tc.testPartBlock.isRecommendedFailureMessage() = %t, want %t", got, tc.want)
			}
		})
	}
}
