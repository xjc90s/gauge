package item

import (
	"github.com/getgauge/gauge/gauge"
	er "github.com/getgauge/gauge/result"
)

// Count of iterations
var MaxRetriesCount int

// NumberOfExecutionStreams shows the number of execution streams, in parallel execution.
var NumberOfExecutionStreams int

// Tags to filter specs/scenarios to retry
var RetryOnlyTags string

var ExecutionArgs []*gauge.ExecutionArg

var tableRowsIndexes []int

// SetTableRows is used to limit data driven execution to specific rows
func SetTableRows(tableRows string) {
	tableRowsIndexes = getDataTableRows(tableRows)
}

type executor interface {
	Execute(i gauge.Item, r er.Result)
}
