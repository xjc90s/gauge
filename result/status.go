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

package result

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/getgauge/common"
	"github.com/getgauge/gauge/config"
	"github.com/getgauge/gauge/logger"
)

const (
	executionStatusFile = "executionStatus.json"
)

type executionStatus struct {
	Type          string `json:"type"`
	SpecsExecuted int    `json:"specsExecuted"`
	SpecsPassed   int    `json:"specsPassed"`
	SpecsFailed   int    `json:"specsFailed"`
	SpecsSkipped  int    `json:"specsSkipped"`
	SceExecuted   int    `json:"sceExecuted"`
	ScePassed     int    `json:"scePassed"`
	SceFailed     int    `json:"sceFailed"`
	SceSkipped    int    `json:"sceSkipped"`
}

func (status *executionStatus) getJSON() (string, error) {
	j, err := json.Marshal(status)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func Print(suiteResult *SuiteResult, isParsingOk bool) int {
	nSkippedSpecs := suiteResult.SpecsSkippedCount
	var nExecutedSpecs int
	if len(suiteResult.SpecResults) != 0 {
		nExecutedSpecs = len(suiteResult.SpecResults) - nSkippedSpecs
	}
	nFailedSpecs := suiteResult.SpecsFailedCount
	nPassedSpecs := nExecutedSpecs - nFailedSpecs

	nExecutedScenarios := 0
	nFailedScenarios := 0
	nPassedScenarios := 0
	nSkippedScenarios := 0
	for _, specResult := range suiteResult.SpecResults {
		nExecutedScenarios += specResult.ScenarioCount
		nFailedScenarios += specResult.ScenarioFailedCount
		nSkippedScenarios += specResult.ScenarioSkippedCount
	}
	nExecutedScenarios -= nSkippedScenarios
	nPassedScenarios = nExecutedScenarios - nFailedScenarios
	if nExecutedScenarios < 0 {
		nExecutedScenarios = 0
	}

	if nPassedScenarios < 0 {
		nPassedScenarios = 0
	}

	s := statusJSON(nExecutedSpecs, nPassedSpecs, nFailedSpecs, nSkippedSpecs, nExecutedScenarios, nPassedScenarios, nFailedScenarios, nSkippedScenarios)
	logger.Infof(true, "Specifications:\t%d executed\t%d passed\t%d failed\t%d skipped", nExecutedSpecs, nPassedSpecs, nFailedSpecs, nSkippedSpecs)
	logger.Infof(true, "Scenarios:\t%d executed\t%d passed\t%d failed\t%d skipped", nExecutedScenarios, nPassedScenarios, nFailedScenarios, nSkippedScenarios)
	logger.Infof(true, "\nTotal time taken: %s", time.Millisecond*time.Duration(suiteResult.ExecutionTime))
	writeExecutionResult(s)

	if !isParsingOk {
		return ParseFailed
	}
	if suiteResult.IsFailed {
		return ExecutionFailed
	}
	return Success
}

func statusJSON(executedSpecs, passedSpecs, failedSpecs, skippedSpecs, executedScenarios, passedScenarios, failedScenarios, skippedScenarios int) string {
	executionStatus := &executionStatus{}
	executionStatus.Type = "out"
	executionStatus.SpecsExecuted = executedSpecs
	executionStatus.SpecsPassed = passedSpecs
	executionStatus.SpecsFailed = failedSpecs
	executionStatus.SpecsSkipped = skippedSpecs
	executionStatus.SceExecuted = executedScenarios
	executionStatus.ScePassed = passedScenarios
	executionStatus.SceFailed = failedScenarios
	executionStatus.SceSkipped = skippedScenarios
	s, err := executionStatus.getJSON()
	if err != nil {
		logger.Fatalf(true, "Unable to parse execution status information : %v", err.Error())
	}
	return s
}

func writeExecutionResult(content string) {
	executionStatusFile := filepath.Join(config.ProjectRoot, common.DotGauge, executionStatusFile)
	dotGaugeDir := filepath.Join(config.ProjectRoot, common.DotGauge)
	if err := os.MkdirAll(dotGaugeDir, common.NewDirectoryPermissions); err != nil {
		logger.Fatalf(true, "Failed to create directory in %s. Reason: %s", dotGaugeDir, err.Error())
	}
	err := ioutil.WriteFile(executionStatusFile, []byte(content), common.NewFilePermissions)
	if err != nil {
		logger.Fatalf(true, "Failed to write to %s. Reason: %s", executionStatusFile, err.Error())
	}
}

// ReadLastExecutionResult returns the result of previous execution in JSON format
// This is stored in $GAUGE_PROJECT_ROOT/.gauge/executionStatus.json file after every execution
func ReadLastExecutionResult() (interface{}, error) {
	contents, err := common.ReadFileContents(filepath.Join(config.ProjectRoot, common.DotGauge, executionStatusFile))
	if err != nil {
		logger.Fatalf(true, "Failed to read execution status information. Reason: %s", err.Error())
	}
	meta := &executionStatus{}
	if err = json.Unmarshal([]byte(contents), meta); err != nil {
		logger.Fatalf(true, "Invalid execution status information. Reason: %s", err.Error())
		return meta, err
	}
	return meta, nil
}
