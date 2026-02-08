package model

import (
	"go/ast"
)

type (
	// TestFunction is the holder of a test function declaration.
	// A test function must:
	// 1. Start with "Test".
	// 2. Have exactly one parameter.
	// 3. Have that parameter be of type *testing.T.
	TestFunction struct {
		// importGroup contains the import important on this test.
		importGroup ImportGroup

		// testVar is the name given to the testing.T parameter
		testVar string

		// funcDecl the original function declaration.
		funcDecl *ast.FuncDecl

		// tableDrivenInfo information about table-driven test, nil if not a table-driven test.
		tableDrivenInfo *TableDrivenInfo
	}

	// ImportGroup contains the imports that are important for the test.
	ImportGroup struct {
		// GoCmp import spec containing go-cmp package. Nil if go-cmp is not imported.
		GoCmp *ast.ImportSpec
		// Reflect import spec containing the "reflect" package. Nil if reflect is not imported.
		Reflect *ast.ImportSpec
	}

	// TableDrivenInfo contains information about table-driven test.
	TableDrivenInfo struct {
		// Range that iterates over the tests and call t.Run
		Range *ast.RangeStmt
		// FormatType is either "map" or "slice".
		FormatType string
		// Inlined is true if the table is declared in the range statement.
		Inlined bool
		// Block is the body of the t.Run function.
		Block *ast.BlockStmt
	}
)

func NewTestFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (TestFunction, bool) {
	ok, testVar := isTestFunction(funcDecl)
	if !ok {
		return TestFunction{}, false
	}

	tbi := newTableDrivenInfo(testVar, funcDecl)

	return TestFunction{
		importGroup:     importGroup,
		testVar:         testVar,
		funcDecl:        funcDecl,
		tableDrivenInfo: tbi,
	}, true
}

func (t TestFunction) ImportGroup() ImportGroup {
	return t.importGroup
}

// GetActualTestBlockStmt returns the actual block test logic, if it's not a table-driven test
// it returns the actual body of the function, and if it's table-driven test it returns
// the content inside the t.Run function.
func (t TestFunction) GetActualTestBlockStmt() *ast.BlockStmt {
	if t.tableDrivenInfo != nil {
		return t.tableDrivenInfo.Block
	}

	return t.funcDecl.Body
}

// GetTestVar returns the name of the testing.T parameter.
func (t TestFunction) GetTestVar() string {
	return t.testVar
}

func (t TestFunction) GetTableDrivenInfo() *TableDrivenInfo {
	return t.tableDrivenInfo
}

func (i ImportGroup) ReflectImportName() (string, bool) {
	if i.Reflect == nil {
		return "", false
	}

	return importName(i.Reflect), true
}

func (i ImportGroup) GoCmpImportName() (string, bool) {
	if i.GoCmp == nil {
		return "", false
	}

	return importName(i.GoCmp), true
}

// newTableDrivenInfo returns information about a table driven test or nil if it's not a table-driven test.
//
//nolint:gocognit,funlen // refactor later
func newTableDrivenInfo(testVar string, funcDecl *ast.FuncDecl) *TableDrivenInfo {
	var stmts []ast.Stmt
	if funcDecl.Body != nil {
		stmts = funcDecl.Body.List
	}

	identifiers := make(map[string]*ast.CompositeLit)

	var rangeStmt *ast.RangeStmt

	for _, stmt := range stmts {
		switch node := stmt.(type) {
		// possible identifiers that can be used in a table-driven test
		case *ast.AssignStmt:
			if len(node.Rhs) != 1 {
				continue
			}

			mapOrSliceCompositeLit := isMapOrSliceCompositeLit(node.Rhs[0])
			if mapOrSliceCompositeLit == nil {
				continue
			}

			if len(node.Lhs) != 1 {
				continue
			}

			if ident, ok := node.Lhs[0].(*ast.Ident); ok {
				identifiers[ident.Name] = mapOrSliceCompositeLit
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

			//nolint:lll // long line
			if ident, isIdent := selectorExpr.X.(*ast.Ident); !isIdent || ident.Name != testVar || selectorExpr.Sel.Name != "Run" {
				continue
			}

			funcLit, isFuncLit := callExpr.Args[1].(*ast.FuncLit)
			if !isFuncLit {
				continue
			}
			// from here, it's a table-driven test, we need to check whether is map/slice or inlined
			rangeStmt = node

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
				if _, isSlice := identifiers[n.Name].Type.(*ast.ArrayType); isSlice {
					formatType = "slice"
				}

				return &TableDrivenInfo{
					Range:      rangeStmt,
					FormatType: formatType,
					Inlined:    false,
					Block:      funcLit.Body,
				}
			case *ast.CompositeLit:
				// is inlined
				if isMapOrSliceCompositeLit(n) == nil {
					continue
				}

				formatType := "map"
				if _, isSlice := n.Type.(*ast.ArrayType); isSlice {
					formatType = "slice"
				}

				return &TableDrivenInfo{
					Range:      rangeStmt,
					FormatType: formatType,
					Inlined:    true,
					Block:      funcLit.Body,
				}
			default:
				continue
			}
		}
	}

	return nil
}
