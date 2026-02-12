package model

import (
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

func isMapOrSliceCompositeLit(expr ast.Expr) *ast.CompositeLit {
	if compositeLit, ok := expr.(*ast.CompositeLit); ok {
		switch compositeLit.Type.(type) {
		case *ast.MapType, *ast.ArrayType:
			return compositeLit
		default:
			return nil
		}
	}

	return nil
}

func isNotBlankIdent(expr ast.Expr) (*ast.Ident, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return nil, false
	}

	if ident.Name == "_" {
		return nil, false
	}

	return ident, true
}

func isGoCmpDiff(goCmpImportAlias string, se *ast.SelectorExpr) bool {
	if se == nil {
		return false
	}

	if ident, ok := se.X.(*ast.Ident); ok {
		if ident.Name == goCmpImportAlias && se.Sel.Name == "Diff" {
			return true
		}
	}

	return false
}

func isGoCmpEqual(goCmpImportAlias string, se *ast.SelectorExpr) bool {
	if se == nil {
		return false
	}

	if ident, ok := se.X.(*ast.Ident); ok {
		if ident.Name == goCmpImportAlias && se.Sel.Name == "Equal" {
			return true
		}
	}

	return false
}

func isReflectEqual(reflectImportAlias string, se *ast.SelectorExpr) bool {
	if se == nil {
		return false
	}

	if ident, ok := se.X.(*ast.Ident); ok {
		if ident.Name == reflectImportAlias && se.Sel.Name == "DeepEqual" {
			return true
		}
	}

	return false
}
