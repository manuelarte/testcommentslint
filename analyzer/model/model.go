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

		// funcDecl the original function declaration.
		funcDecl *ast.FuncDecl

		// testVar is the name given to the testing.T parameter
		testVar string
	}

	// ImportGroup contains the imports that are important for the test.
	ImportGroup struct {
		// GoCmp import spec containing go-cmp package. Nil if go-cmp is not imported.
		GoCmp *ast.ImportSpec
		// Reflect import spec containing the "reflect" package. Nil if reflect is not imported.
		Reflect *ast.ImportSpec
	}
)

func NewTestFunction(importGroup ImportGroup, funcDecl *ast.FuncDecl) (TestFunction, bool) {
	ok, testVar := isTestFunction(funcDecl)
	if !ok {
		return TestFunction{}, false
	}

	return TestFunction{
		importGroup: importGroup,
		funcDecl:    funcDecl,
		testVar:     testVar,
	}, true
}

func (t TestFunction) ReflectImportName() (string, bool) {
	if t.importGroup.Reflect == nil {
		return "", false
	}

	return importName(t.importGroup.Reflect), true
}

func (t TestFunction) GoCmpImportName() (string, bool) {
	if t.importGroup.GoCmp == nil {
		return "", false
	}

	return importName(t.importGroup.GoCmp), true
}

// GetActualTestBlockStmt returns the actual block test logic, if it's not a table-driven test
// it returns the actual body of the function, and if it's table-driven test it returns
// the content inside the t.Run function.
func (t TestFunction) GetActualTestBlockStmt() *ast.BlockStmt {
	if ok, bl := isTableDrivenTest(t.funcDecl); ok {
		return bl
	}

	return t.funcDecl.Body
}

// GetTestVar returns the name of the testing.T parameter.
func (t TestFunction) GetTestVar() string {
	return t.testVar
}
