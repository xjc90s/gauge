// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

/*Package execution handles gauge's execution of spec/scenario/steps
   Execution can be of two types
	- Simple execution
	- Paralell execution

   Execution Flow :
   	- Checks for updates
    	- Validation
    	- Init Registry
    	- Saving Execution result

   Strategy
    	- Lazy : Lazy is a parallelization strategy for execution. In this case tests assignment will be dynamic during execution, i.e. assign the next spec in line to the stream that has completed it’s previous execution and is waiting for more work.
    	- Eager : Eager is a parallelization strategy for execution. In this case tests are distributed before execution, thus making them an equal number based distribution.
*/
package execution

import (
	"github.com/getgauge/gauge/plugin"

	"github.com/getgauge/gauge/execution/event"
	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/logger"
	"github.com/getgauge/gauge/manifest"
	er "github.com/getgauge/gauge/result"
	"github.com/getgauge/gauge/runner"
)

// Count of iterations
var MaxRetriesCount int

// Tags to filter specs/scenarios to retry
var RetryOnlyTags string

// NumberOfExecutionStreams shows the number of execution streams, in parallel execution.
var NumberOfExecutionStreams int

// InParallel if true executes the specs in parallel else in serial.
var InParallel bool

var TagsToFilterForParallelRun string

// Verbose if true prints additional details about the execution
var Verbose bool

// MachineReadable indicates that the output is in json format
var MachineReadable bool

var ExecutionArgs []*gauge.ExecutionArg

type suiteExecutor interface {
	run() *er.SuiteResult
}

type executor interface {
	execute(i gauge.Item, r er.Result)
}

type executionInfo struct {
	manifest        *manifest.Manifest
	specs           *gauge.SpecCollection
	runner          runner.Runner
	pluginHandler   plugin.Handler
	errMaps         *gauge.BuildErrors
	inParallel      bool
	numberOfStreams int
	tagsToFilter    string
	stream          int
}

func newExecutionInfo(s *gauge.SpecCollection, r runner.Runner, ph plugin.Handler, e *gauge.BuildErrors, p bool, stream int) *executionInfo {
	m, err := manifest.ProjectManifest()
	if err != nil {
		logger.Fatalf(true, err.Error())
	}
	return &executionInfo{
		manifest:        m,
		specs:           s,
		runner:          r,
		pluginHandler:   ph,
		errMaps:         e,
		inParallel:      p,
		numberOfStreams: NumberOfExecutionStreams,
		tagsToFilter:    TagsToFilterForParallelRun,
		stream:          stream,
	}
}

// ExecuteSpecs : Check for updates, validates the specs (by invoking the respective language runners), initiates the registry which is needed for console reporting, execution API and Rerunning of specs
// and finally saves the execution result as binary in .gauge folder.
var ExecuteSpecs = func(res *gauge.ValidationResult, r runner.Runner, specDirs []string) *er.SuiteResult {
	event.InitRegistry()
	ei := newExecutionInfo(res.SpecCollection, r, nil, res.ErrMap, InParallel, 0)

	return newExecution(ei).run()
}

func newExecution(executionInfo *executionInfo) suiteExecutor {
	if executionInfo.inParallel {
		return newParallelExecution(executionInfo)
	}
	return newSimpleExecution(executionInfo, true)
}
