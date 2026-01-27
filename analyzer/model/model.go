package model

import (
	"go/ast"
)

type ReflectImport struct {
	is *ast.ImportSpec
}

func NewReflectImport(i *ast.ImportSpec) (*ReflectImport, bool) {
	if i.Path == nil || i.Path.Value != "\"reflect\"" {
		return nil, false
	}

	return &ReflectImport{is: i}, true
}

func (i ReflectImport) ImportName() string {
	if i.is.Name != nil {
		return i.is.Name.Name
	}

	return i.is.Path.Value
}

// TestFunction is the holder of a test function declaration.
// A test function must:
// 1. Start with "Test".
// 2. Have exactly one parameter.
// 3. Have that parameter be of type *testing.T.
type TestFunction struct {
	// reflectImport import spec containing the "reflect" package
	reflectImport *ReflectImport
	// funcDecl the original function declaration.
	funcDecl *ast.FuncDecl

	// testVar is the name given to the testing.T parameter
	testVar string
}

func NewTestFunction(reflectImport *ReflectImport, funcDecl *ast.FuncDecl) (TestFunction, bool) {
	ok, testVar := isTestFunction(funcDecl)
	if !ok {
		return TestFunction{}, false
	}

	return TestFunction{
		reflectImport: reflectImport,
		funcDecl:      funcDecl,
		testVar:       testVar,
	}, true
}

func (t TestFunction) ReflectImportName() (string, bool) {
	if t.reflectImport == nil {
		return "", false
	}

	return t.reflectImport.ImportName(), true
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
