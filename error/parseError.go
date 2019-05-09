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

package error

import "fmt"

// ParseError holds information about a parse failure
type ParseError struct {
	FileName string
	LineNo   int
	Message  string
	LineText string
}

// Error prints error with filename, line number, error message and step text.
func (se ParseError) Error() string {
	if se.LineNo == 0 && se.FileName == "" {
		return fmt.Sprintf("%s", se.Message)
	}
	return fmt.Sprintf("%s:%d %s => '%s'", se.FileName, se.LineNo, se.Message, se.LineText)
}
