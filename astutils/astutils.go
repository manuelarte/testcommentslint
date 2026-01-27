package astutils

import (
	"go/ast"
	"strings"
)

// IsTestFunction checks if a function declaration is a test function
// A test function must:
// 1. Start with "Test"
// 2. Have exactly one parameter
// 3. Have that parameter be of type *testing.T
// Returns (true, paramName) if it is a test function, (false, "") if it isn't.
//
//nolint:nestif // no need
func IsTestFunction(funcDecl *ast.FuncDecl) (bool, string) {
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

// IsTableDrivenTest checks if a test function is a table driven test and returns the BlockStmt.
//
//nolint:gocognit // refactor later
func IsTableDrivenTest(funcDecl *ast.FuncDecl) (bool, *ast.BlockStmt) {
	isTestFunction, testVar := IsTestFunction(funcDecl)
	if !isTestFunction {
		return false, nil
	}

	identifiers := make(map[string]struct{})

	var stmts []ast.Stmt
	if funcDecl.Body != nil {
		stmts = funcDecl.Body.List
	}

	for _, stmt := range stmts {
		switch node := stmt.(type) {
		// possible identifiers that can be used in a table-driven test
		case *ast.AssignStmt:
			if len(node.Rhs) != 1 {
				continue
			}

			if !isMapOrSlice(node.Rhs[0]) {
				continue
			}

			if len(node.Lhs) != 1 {
				continue
			}

			if ident, ok := node.Lhs[0].(*ast.Ident); ok {
				identifiers[ident.Name] = struct{}{}
			}
		// possible for loops that can be used in a table-driven test
		case *ast.RangeStmt:
			// identifier must be declared before, and be used as range
			rangeIdent, ok := node.X.(*ast.Ident)
			if !ok {
				continue
			}

			if _, isDeclaredBefore := identifiers[rangeIdent.Name]; !isDeclaredBefore {
				continue
			}
			// next instruction in a range stmt needs to be a t.Run
			if node.Body != nil && len(node.Body.List) != 1 {
				continue
			}

			exprStmt, isExprStmt := node.Body.List[0].(*ast.ExprStmt)
			if !isExprStmt {
				continue
			}

			callExpr, isCallExpr := exprStmt.X.(*ast.CallExpr)
			if !isCallExpr {
				continue
			}

			selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
			if !isSelectorExpr {
				continue
			}

			if selectorExpr.Sel.Name != "Run" {
				continue
			}

			if ident, isIdent := selectorExpr.X.(*ast.Ident); isIdent && ident.Name != testVar {
				continue
			}

			return true, node.Body
		}
	}

	return false, nil
}

func isMapOrSlice(expr ast.Expr) bool {
	if compositeLit, ok := expr.(*ast.CompositeLit); ok {
		switch compositeLit.Type.(type) {
		case *ast.MapType, *ast.ArrayType:
			return true
		default:
			return false
		}
	}

	return false
}
