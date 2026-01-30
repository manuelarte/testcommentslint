package checks

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// FailureMessage check that the failure messages in t.Errorf follow the format expected.
// The format expected can be as the following:
//   - when the condition is reflect.DeepEqual, cmp.Equal or got != want: "YourFunction(%v) = %v, want %v"
//   - when the condition is cmp.Diff: YourFunction() mismatch (-want +got):\n%s
type FailureMessage struct {
	category string
}

// NewFailureMessage creates a new FailureMessage.
func NewFailureMessage() FailureMessage {
	return FailureMessage{
		category: "Failure Message",
	}
}

func (c FailureMessage) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	blStmt := testFunc.GetActualTestBlockStmt()
	testVar := testFunc.GetTestVar()

	var stmts []ast.Stmt
	if blStmt != nil {
		stmts = blStmt.List
	}

	for i, stmt := range stmts {
		switch node := stmt.(type) {
		case *ast.IfStmt:
			// Check the condition to see if it's a comparison like got != want
			gwStmt, ok := newGotWantIfStmt(testVar, node)
			if !ok {
				continue
			}

			fmt.Printf("%v, %v\n", gwStmt, ok)

			if c.isEqualityOrDiffCondition(node.Cond) {
				// Check the body for t.Errorf calls
				if node.Body != nil {
					// Extract function name from preceding statements
					functionName := c.extractFunctionName(stmts[:i])
					c.checkStmtListWithContext(pass, node.Body.List, testVar, functionName)
				}
			}
		}
	}
}

// isEqualityOrDiffCondition checks if a condition is an equality/inequality check or DeepEqual call.
func (c FailureMessage) isEqualityOrDiffCondition(cond ast.Expr) bool {
	switch node := cond.(type) {
	case *ast.BinaryExpr:
		// Check for operators like !=, ==
		return node.Op.String() == "!=" || node.Op.String() == "=="
	case *ast.UnaryExpr:
		// Check for ! prefix (e.g., !reflect.DeepEqual)
		return true
	case *ast.CallExpr:
		// Check for function calls like reflect.DeepEqual, cmp.Equal, cmp.Diff
		return true
	}

	return false
}

// checkStmtList checks a list of statements for t.Errorf calls.
func (c FailureMessage) checkStmtList(pass *analysis.Pass, stmts []ast.Stmt, testVar string) {
	c.checkStmtListWithContext(pass, stmts, testVar, "")
}

// checkStmtListWithContext checks a list of statements for t.Errorf calls with function context.
func (c FailureMessage) checkStmtListWithContext(pass *analysis.Pass, stmts []ast.Stmt, testVar, functionName string) {
	for _, stmt := range stmts {
		switch node := stmt.(type) {
		case *ast.ExprStmt:
			if callExpr, ok := node.X.(*ast.CallExpr); ok {
				c.checkErrorfCallWithContext(pass, callExpr, testVar, functionName)
			}
		}
	}
}

// extractFunctionName extracts the function name and parameter count from preceding statements.
func (c FailureMessage) extractFunctionName(stmts []ast.Stmt) string {
	// Look for assignment statements that contain function calls
	// Prefer the most recent assignment that might be related to the test
	for i := len(stmts) - 1; i >= 0; i-- {
		switch node := stmts[i].(type) {
		case *ast.AssignStmt:
			if len(node.Rhs) > 0 {
				if callExpr, ok := node.Rhs[0].(*ast.CallExpr); ok {
					funcName := c.extractCallExprName(callExpr)

					paramCount := len(callExpr.Args)
					if funcName != "" {
						// Create format like "sum(%v, %v)" based on parameter count
						placeholders := make([]string, paramCount)
						for j := range paramCount {
							placeholders[j] = "%v"
						}

						return funcName + "(" + strings.Join(placeholders, ", ") + ")"
					}
				}
			}
		}
	}

	return ""
}

// extractCallExprName extracts the function name from a call expression.
func (c FailureMessage) extractCallExprName(callExpr *ast.CallExpr) string {
	switch fn := callExpr.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		if ident, ok := fn.X.(*ast.Ident); ok {
			return ident.Name + "." + fn.Sel.Name
		}
	}

	return ""
}

// checkErrorfCall checks if a t.Errorf call has the correct format.
func (c FailureMessage) checkErrorfCall(pass *analysis.Pass, callExpr *ast.CallExpr, testVar string) {
	c.checkErrorfCallWithContext(pass, callExpr, testVar, "")
}

// checkErrorfCallWithContext checks if a t.Errorf call has the correct format with function context.
func (c FailureMessage) checkErrorfCallWithContext(pass *analysis.Pass, callExpr *ast.CallExpr, testVar, functionName string) {
	// Check if this is a t.Errorf call
	if !c.isErrorfCall(callExpr, testVar) {
		return
	}

	if len(callExpr.Args) < 1 {
		return
	}

	// Get the format string
	formatArg := callExpr.Args[0]

	formatStr, ok := formatArg.(*ast.BasicLit)
	if !ok || formatStr.Kind.String() != "STRING" {
		return
	}

	// Remove quotes from the format string
	format := strings.Trim(formatStr.Value, "\"")

	// Check if the format follows the expected pattern
	// Expected: "functionName(...) = ..., want ..."
	if !c.isValidFormatString(format) {
		// Generate a suggested format
		suggestedFormat := c.generateSuggestedFormatWithContext(format, functionName)
		// Build message: Prefer "suggestedFormat" format for failure message
		var sb strings.Builder
		sb.WriteString("Prefer \"")
		sb.WriteString(suggestedFormat)
		sb.WriteString("\" format for failure message")
		message := sb.String()
		diag := &analysis.Diagnostic{
			Pos:      callExpr.Pos(),
			End:      callExpr.End(),
			Category: c.category,
			Message:  message,
			URL:      "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#failure-message",
		}
		pass.Report(*diag)
	}
}

// isErrorfCall checks if a call expression is a t.Errorf call.
func (c FailureMessage) isErrorfCall(callExpr *ast.CallExpr, testVar string) bool {
	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := selectorExpr.X.(*ast.Ident); ok {
			return ident.Name == testVar && selectorExpr.Sel.Name == "Errorf"
		}
	}

	return false
}

// isValidFormatString checks if a format string matches the expected pattern.
func (c FailureMessage) isValidFormatString(format string) bool {
	// Pattern should include function name with parameters and "want"
	// e.g., "sum(%v, %v) = %v, want %v" or "foo() mismatch (-want +got):\n%s"

	// Check if it contains the expected pattern
	if strings.Contains(format, "want") {
		// Should have format like "functionName(...) = ..., want ..."
		return strings.Contains(format, "=")
	}

	return false
}

// generateSuggestedFormat generates a suggested format based on the provided format.
func (c FailureMessage) generateSuggestedFormat(format string) string {
	return c.generateSuggestedFormatWithContext(format, "")
}

// generateSuggestedFormatWithContext generates a suggested format with function context.
func (c FailureMessage) generateSuggestedFormatWithContext(format, functionName string) string {
	// Parse the format to extract function call pattern
	// If format is "got %v, want %v", suggest "functionName(%v, %v) = %v, want %v"

	// Simple heuristic: if it has "got %v, want %v", suggest the proper format
	if strings.Contains(format, "got") && strings.Contains(format, "want") {
		// Extract the pattern - try to find how many %v are used
		pattern := regexp.MustCompile(`%[a-z]`)
		matches := pattern.FindAllString(format, -1)

		// For a simple got/want pattern with 2 matches, suggest the proper format
		if len(matches) == 2 {
			if functionName != "" {
				return functionName + " = %v, want %v"
			}

			return "functionName(%v, %v) = %v, want %v"
		}
	}

	return format
}

// Auxiliary structs to facilitate the business logic.
type (
	// gotWantIfStatement struct holding an if statement that contains a comparison of got and want.
	// Conditions allowed:
	// - got != want
	// - reflect.DeepEqual
	// - cmp.Equal
	// - cmp.Diff
	// and inside the if statement there is a call to t.Errorf.
	gotWantIfStmt struct {
		// original ifStmt
		ifStmt *ast.IfStmt

		params        []*ast.Ident
		errorCallExpr tErrorfCallExpr
	}

	tErrorfCallExpr struct {
		callExpr *ast.CallExpr

		params []*ast.Ident
	}
)

// newGotWantIfStmt creates a new gotWantIfStmt.
// only if the condition applies.
func newGotWantIfStmt(testVar string, ifStmt *ast.IfStmt) (gotWantIfStmt, bool) {
	if ifStmt == nil || ifStmt.Body == nil {
		return gotWantIfStmt{}, false
	}

	if len(ifStmt.Body.List) != 1 {
		return gotWantIfStmt{}, false
	}

	teCallExpr, istErrorf := newTErrorfCallExpr(testVar, ifStmt.Body.List[0])
	if !istErrorf {
		return gotWantIfStmt{}, false
	}

	params := make([]*ast.Ident, 2)

	switch node := ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		// check ident1 != ident2 and both are used in the failure message t.Errorf
		if node.Op.String() != "!=" {
			return gotWantIfStmt{}, false
		}

		xIdent, isXIdent := node.X.(*ast.Ident)

		yIdent, isYIdent := node.Y.(*ast.Ident)
		if !isXIdent || !isYIdent {
			return gotWantIfStmt{}, false
		}

		params[0] = xIdent
		params[1] = yIdent
	default:
		return gotWantIfStmt{}, false
	}

	// check params match
	if params[0].Name != teCallExpr.params[0].Name && params[0].Name != teCallExpr.params[1].Name {
		return gotWantIfStmt{}, false
	}

	if params[1].Name != teCallExpr.params[0].Name && params[1].Name != teCallExpr.params[1].Name {
		return gotWantIfStmt{}, false
	}

	return gotWantIfStmt{
		ifStmt:        ifStmt,
		params:        params,
		errorCallExpr: teCallExpr,
	}, true
}

func newTErrorfCallExpr(testVar string, stmt ast.Stmt) (tErrorfCallExpr, bool) {
	exprStmt, isExprStmt := stmt.(*ast.ExprStmt)
	if !isExprStmt {
		return tErrorfCallExpr{}, false
	}

	callExpr, isCallExpr := exprStmt.X.(*ast.CallExpr)
	if !isCallExpr {
		return tErrorfCallExpr{}, false
	}

	selectorExpr, isSelectorExpr := callExpr.Fun.(*ast.SelectorExpr)
	if !isSelectorExpr {
		return tErrorfCallExpr{}, false
	}

	ident, isIdent := selectorExpr.X.(*ast.Ident)
	if !isIdent {
		return tErrorfCallExpr{}, false
	}

	if ident.Name != testVar || selectorExpr.Sel.Name != "Errorf" {
		return tErrorfCallExpr{}, false
	}

	if len(callExpr.Args) != 3 {
		return tErrorfCallExpr{}, false
	}

	firstIdent, isFirstIdent := callExpr.Args[1].(*ast.Ident)

	secondIdent, isSecondIdent := callExpr.Args[2].(*ast.Ident)
	if !isFirstIdent || !isSecondIdent {
		return tErrorfCallExpr{}, false
	}

	return tErrorfCallExpr{
		callExpr: callExpr,
		params:   []*ast.Ident{firstIdent, secondIdent},
	}, true
}
