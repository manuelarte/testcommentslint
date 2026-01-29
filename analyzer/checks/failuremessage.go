package checks

import (
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
	// TODO(manuelarte): implement check
}
