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

package util

import (
	"os"

	"github.com/getgauge/gauge/config"
	"github.com/getgauge/gauge/env"
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestSpecDirEnvVariableAllowsCommaSeparatedList(c *C) {
	os.Clearenv()
	config.ProjectRoot = "_testdata/proj1"

	e := env.LoadEnv("multipleSpecs")
	c.Assert(e, Equals, nil)
	c.Assert(GetSpecDirs(), DeepEquals, []string{"spec1", "spec2", "spec3"})
}

func (s *MySuite) TestUnescapedString(c *C) {
	unEscapedString := GetUnescapedString("hello \n world")
	c.Assert(unEscapedString, Equals, `hello \n world`)

	unEscapedString = GetUnescapedString("hello \n \"world")
	c.Assert(unEscapedString, Equals, `hello \n \"world`)

	unEscapedString = GetUnescapedString("\"hello \n \"world\"\"")
	c.Assert(unEscapedString, Equals, `\"hello \n \"world\"\"`)

	unEscapedString = GetUnescapedString("\"\"")
	c.Assert(unEscapedString, Equals, `\"\"`)

	unEscapedString = GetUnescapedString("")
	c.Assert(unEscapedString, Equals, "")

}
