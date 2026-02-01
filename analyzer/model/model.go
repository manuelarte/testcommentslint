package model

import (
	"go/ast"
)

// TestFunction is the holder of a test function declaration.
// A test function must:
// 1. Start with "Test".
// 2. Have exactly one parameter.
// 3. Have that parameter be of type *testing.T.
type TestFunction struct {
	// goCmpImport import spec containing go-cmp package. Nil if go-cmp is not imported.
	goCmpImport *ast.ImportSpec
	// reflectImport import spec containing the "reflect" package. Nil if reflect is not imported.
	reflectImport *ast.ImportSpec

	// funcDecl the original function declaration.
	funcDecl *ast.FuncDecl

	// testVar is the name given to the testing.T parameter
	testVar string
}

func NewTestFunction(goCmpImport, reflectImport *ast.ImportSpec, funcDecl *ast.FuncDecl) (TestFunction, bool) {
	ok, testVar := isTestFunction(funcDecl)
	if !ok {
		return TestFunction{}, false
	}

	return TestFunction{
		goCmpImport:   goCmpImport,
		reflectImport: reflectImport,
		funcDecl:      funcDecl,
		testVar:       testVar,
	}, true
}

func (t TestFunction) ReflectImportName() (string, bool) {
	if t.reflectImport == nil {
		return "", false
	}

	return importName(t.reflectImport), true
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
