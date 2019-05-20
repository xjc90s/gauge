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

package resolver

import (
	"testing"

	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/util"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestParsingFileSpecialType(c *C) {
	resolver := NewSpecialTypeResolver()
	resolver.predefinedResolvers["file"] = func(value string) (*gauge.StepArg, error) {
		return &gauge.StepArg{Value: "dummy", ArgType: gauge.Static}, nil
	}

	stepArg, _ := resolver.Resolve("file:foo")
	c.Assert(stepArg.Value, Equals, "dummy")
	c.Assert(stepArg.ArgType, Equals, gauge.Static)
	c.Assert(stepArg.Name, Equals, "file:foo")
}

func (s *MySuite) TestParsingFileAsSpecialParamWithWindowsPathAsValue(c *C) {
	resolver := NewSpecialTypeResolver()
	resolver.predefinedResolvers["file"] = func(value string) (*gauge.StepArg, error) {
		return &gauge.StepArg{Value: "hello", ArgType: gauge.SpecialString}, nil
	}

	stepArg, _ := resolver.Resolve(`file:C:\Users\abc`)
	c.Assert(stepArg.Value, Equals, "hello")
	c.Assert(stepArg.ArgType, Equals, gauge.SpecialString)
	if util.IsWindows() {
		c.Assert(stepArg.Name, Equals, `file:C:\\Users\\abc`)
	} else {
		c.Assert(stepArg.Name, Equals, `file:C:\Users\abc`)
	}
}

func (s *MySuite) TestParsingInvalidSpecialType(c *C) {
	resolver := NewSpecialTypeResolver()

	_, err := resolver.Resolve("unknown:foo")
	c.Assert(err.Error(), Equals, "Resolver not found for special param <unknown:foo>")
}

func (s *MySuite) TestConvertCsvToTable(c *C) {
	table, _ := convertCsvToTable("id,name\n1,foo\n2,bar")

	idColumn, _ := table.Get("id")
	c.Assert(idColumn[0].Value, Equals, "1")
	c.Assert(idColumn[1].Value, Equals, "2")

	nameColumn, _ := table.Get("name")
	c.Assert(nameColumn[0].Value, Equals, "foo")
	c.Assert(nameColumn[1].Value, Equals, "bar")
}

func (s *MySuite) TestConvertEmptyCsvToTable(c *C) {
	table, _ := convertCsvToTable("")
	c.Assert(len(table.Columns), Equals, 0)
}

func (s *MySuite) TestParsingUnknownSpecialType(c *C) {
	resolver := NewSpecialTypeResolver()

	_, err := resolver.getStepArg("unknown", "foo", "unknown:foo")
	c.Assert(err.Error(), Equals, "Resolver not found for special param <unknown:foo>")
}
