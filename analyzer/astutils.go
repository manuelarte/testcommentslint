package analyzer

import (
	"go/ast"
	"strings"
)

// isTestFunction checks if a function declaration is a test function
// A test function must:
// 1. Start with "Test"
// 2. Have exactly one parameter
// 3. Have that parameter be of type *testing.T
// Returns (true, paramName) if it is a test function, (false, "") if it isn't.
//
//nolint:nestif // no need
func isTestFunction(funcDecl *ast.FuncDecl) (bool, string) {
	testMethodPackageType := "testing"
	testMethodStruct := "T"
	testPrefix := "Test"

	if !strings.HasPrefix(funcDecl.Name.Name, testPrefix) {
		return false, ""
	}

	if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) != 1 {
		return false, ""
	}

	param := funcDecl.Type.Params.List[0]
	if starExp, ok := param.Type.(*ast.StarExpr); ok {
		if selectExpr, isSelector := starExp.X.(*ast.SelectorExpr); isSelector {
			if selectExpr.Sel.Name == testMethodStruct {
				if s, isIdent := selectExpr.X.(*ast.Ident); isIdent {
					if len(param.Names) > 0 {
						return s.Name == testMethodPackageType, param.Names[0].Name
					}
				}
			}
		}
	}

	return false, ""
}
