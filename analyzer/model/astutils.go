package model

import (
	"fmt"
	"go/ast"
	"strings"
)

// IsReflectImport returns whether the ast.ImportSpec is from the reflect package.
func IsReflectImport(i *ast.ImportSpec) bool {
	if i.Path == nil || i.Path.Value != "\"reflect\"" {
		return false
	}

	return true
}

// IsGoCmpImport returns whether the ast.ImportSpec is from go-cmp package.
func IsGoCmpImport(i *ast.ImportSpec) bool {
	if i.Path == nil || i.Path.Value != "\"github.com/google/go-cmp/cmp\"" {
		return false
	}

	return true
}

// importName returns the import name for a package, either the actual name of the alias.
func importName(is *ast.ImportSpec) string {
	if is.Name != nil {
		return is.Name.Name
	}

	unquoted := is.Path.Value[1 : len(is.Path.Value)-1]
	parts := strings.Split(unquoted, "/")

	return parts[len(parts)-1]
}

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

// isTableDrivenTest checks if a test function is a table driven test and returns the BlockStmt.
//
//nolint:gocognit // refactor later
func isTableDrivenTest(funcDecl *ast.FuncDecl) (bool, *ast.BlockStmt) {
	type tableDrivenStruct struct {
		formatType string
		// tests are stored in this identifier, nil if inlined
		param *ast.Ident
	}

	var stmts []ast.Stmt
	if funcDecl.Body != nil {
		stmts = funcDecl.Body.List
	}

	identifiers := make(map[string]*ast.CompositeLit)

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
				identifiers[ident.Name] = node.Rhs[0].(*ast.CompositeLit)
			}
		// possible for loops that can be used in a table-driven test
		case *ast.RangeStmt:
			// the next instruction in a range stmt needs to be a t.Run
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

			if ident, isIdent := selectorExpr.X.(*ast.Ident); !isIdent || ident.Name != "t" || selectorExpr.Sel.Name != "Run" {
				continue
			}

			funcLit, isFuncLit := callExpr.Args[1].(*ast.FuncLit)
			if !isFuncLit {
				continue
			}
			// from here, it's a table-driven test, we need to check whether is map/slice or inlined

			switch n := node.X.(type) {
			case *ast.Ident:
				// identifier must be declared before and be used as range
				if _, isDeclaredBefore := identifiers[n.Name]; !isDeclaredBefore {
					continue
				}

				// returns the second parameter of t.Run with the function
				if len(callExpr.Args) != 2 {
					continue
				}

				// is non-inlined and the param contains whether is map/slice
				formatType := "map"
				if _, isSlice := identifiers[n.Name].Type.(*ast.SliceExpr); isSlice {
					formatType = "slice"
				}

				toReturn := tableDrivenStruct{
					formatType: formatType,
					param:      n,
				}
				fmt.Printf("%v", toReturn)

				return true, funcLit.Body
			case *ast.CompositeLit:
				// is inlined
				if !isMapOrSlice(n) {
					continue
				}

				formatType := "map"
				if _, isSlice := n.Type.(*ast.SliceExpr); isSlice {
					formatType = "slice"
				}

				toReturn := tableDrivenStruct{
					formatType: formatType,
				}
				fmt.Printf("%+v\n", toReturn)

				return true, funcLit.Body

			default:
				continue
			}
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
