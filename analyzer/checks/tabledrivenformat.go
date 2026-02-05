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

	return func(testFunc model.TestFunction) *analysis.Diagnostic {
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
	// TODO implement me
}
