package checks

import (
	"fmt"

	"golang.org/x/tools/go/analysis"

	"github.com/manuelarte/testcommentslint/analyzer/model"
)

const (
	Map   TableDrivenFormatType = "map"
	Slice TableDrivenFormatType = "slice"
)

type (
	TableDrivenFormatType      string
	TableDrivenFormatPredicate func(testFunc model.TestFunction) *analysis.Diagnostic

	TableDrivenFormat struct {
		pred TableDrivenFormatPredicate

		category string
	}

	TableDrivenFormatTypeError struct {
		requestedFormatType TableDrivenFormatType
	}
)

func (e TableDrivenFormatTypeError) Error() string {
	return fmt.Sprintf("table format type not expected: %q", e.requestedFormatType)
}

func AlwaysValid() TableDrivenFormatPredicate {
	return func(testFunc model.TestFunction) *analysis.Diagnostic {
		return nil
	}
}

func OfTypeAndInline(formatType TableDrivenFormatType, inline bool) (TableDrivenFormatPredicate, error) {
	if formatType != Map && formatType != Slice {
		return nil, TableDrivenFormatTypeError{requestedFormatType: formatType}
	}

	inlinedNonInlinedMessage := "inlined"
	if !inline {
		inlinedNonInlinedMessage = "non-inlined"
	}

	expectedMessage := fmt.Sprintf("Expected %s-%s table driven test", formatType, inlinedNonInlinedMessage)

	return func(testFunc model.TestFunction) *analysis.Diagnostic {
		info := testFunc.GetTableDrivenInfo()
		if testFunc.GetTableDrivenInfo().FormatType != string(formatType) || info.Inlined != inline {
			return &analysis.Diagnostic{
				Pos:     info.Range.Pos(),
				End:     info.Range.End(),
				Message: expectedMessage,
			}
		}

		return nil
	}, nil
}

// NewTableDrivenFormat creates a new TableDrivenFormat.
func NewTableDrivenFormat(pred TableDrivenFormatPredicate) (TableDrivenFormat, error) {
	if pred == nil {
		pred = AlwaysValid()
	}

	return TableDrivenFormat{
		pred:     pred,
		category: "Table-Driven Format",
	}, nil
}

func (c TableDrivenFormat) Check(pass *analysis.Pass, testFunc model.TestFunction) {
	info := testFunc.GetTableDrivenInfo()
	if info == nil {
		return
	}

	diag := c.pred(testFunc)
	if diag != nil {
		diag.Category = c.category
		pass.Report(*diag)
	}
}
