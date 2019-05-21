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

package parallel

import (
	"github.com/getgauge/gauge/execution/item"
	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/logger"
	"github.com/getgauge/gauge/manifest"
	"github.com/getgauge/gauge/plugin"
	"github.com/getgauge/gauge/runner"
)

type executionInfo struct {
	manifest        *manifest.Manifest
	specs           *gauge.SpecCollection
	runner          runner.Runner
	pluginHandler   plugin.Handler
	errMaps         *gauge.BuildErrors
	numberOfStreams int
	tagsToFilter    string
	stream          int
}

func newExecutionInfo(s *gauge.SpecCollection, r runner.Runner, ph plugin.Handler, e *gauge.BuildErrors) *executionInfo {
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
		numberOfStreams: item.NumberOfExecutionStreams,
		tagsToFilter:    TagsToFilterForParallelRun,
	}
}
