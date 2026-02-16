package checks

import (
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

// IdentifyFunction check that the failure messages in t.Errorf/Fatalf contains the function name.
type IdentifyFunction struct {
	category string
}

// NewIdentifyFunction creates a new IdentifyFunction.
func NewIdentifyFunction() IdentifyFunction {
	return IdentifyFunction{
		category: "Identify The Function",
	}
}

// Check checks that the failure messages in t.Errorf/Fatalf follow the format expected.
func (c IdentifyFunction) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	for _, testBlock := range testFunc.TestPartBlocks() {
		if containsFunctionName(testBlock) {
			continue
		}

		diag := analysis.Diagnostic{
			Pos:      testBlock.TErrorCallExpr().CallExpr().Pos(),
			End:      testBlock.TErrorCallExpr().CallExpr().End(),
			Category: c.category,
			Message:  "Failure messages should include the name of the function that failed",
			URL:      "https://github.com/manuelarte/testcommentslint/tree/main?tab=readme-ov-file#identify-the-function",
		}
		pass.Report(diag)
	}
}

// containsFunctionName returns whether the failure message contains the function name.
func containsFunctionName(t model.TestPartBlock) bool {
	currentFailureMessage := t.TErrorCallExpr().FailureMessage()

	unquoted, err := strconv.Unquote(currentFailureMessage)
	if err != nil {
		// It's not a string literal that can be unquoted, maybe a raw string literal.
		// We'll check the content as is.
		unquoted = currentFailureMessage
	}

	funName := t.TestedFunc().FunctionName()

	return containsFunctionNameString(funName, unquoted)
}

func containsFunctionNameString(functionName, failureMessage string) bool {
	// Extract the last part of the function name (in case it's a selector expression like "test.YourFunction")
	parts := strings.Split(functionName, ".")
	lastFunctionName := parts[len(parts)-1]

	// Check if the failure message contains the full function name
	if strings.Contains(failureMessage, functionName) {
		return true
	}

	// Check if the failure message contains just the last part
	if strings.Contains(failureMessage, lastFunctionName) {
		// If the function name has multiple parts (has a selector), we need to be more careful
		// We should reject if there's a different selector in the message
		if len(parts) > 1 {
			// Check if there's a selector expression in the message
			// Pattern: something.lastFunctionName where something is not the expected selector
			pattern := `\w+\.` + regexp.QuoteMeta(lastFunctionName)
			if matched, _ := regexp.MatchString(pattern, failureMessage); matched {
				// There's a selector expression in the message, check if it matches the expected one
				expectedPrefix := strings.Join(parts[:len(parts)-1], ".")
				// Check if the expected prefix is in the message
				return strings.Contains(failureMessage, expectedPrefix+"."+lastFunctionName)
			}
		}

		return true
	}

	return false
}
