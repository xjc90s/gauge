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

// +build js

package main

import (
	"fmt"
	"net"
	"syscall/js"

	"github.com/getgauge/gauge/gauge_messages"

	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/runner"

	"github.com/getgauge/gauge/parser"
	"github.com/getgauge/gauge/validation"
)

var signal = make(chan int)

func keepAlive() {
	for {
		<-signal
	}
}

func parse() {
	document := js.Global().Get("document")
	specEl := document.Call("getElementById", "specText")
	content := specEl.Get("innerText").String()
	s, r := new(parser.SpecParser).ParseSpecText(content, "browser")
	console := js.Global().Get("console")
	if !r.Ok {
		for _, e := range r.Errors() {
			console.Call("error", e)
		}
		return
	}
	vErrs := validation.NewValidator([]*gauge.Specification{s}, newInBrowserRunner(), gauge.NewConceptDictionary()).Validate()
	for _, e := range vErrs[s] {
		console.Call("error", e.Error())
	}

}

func main() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		parse()
		return nil
	})
	document := js.Global().Get("document")
	parseButton := document.Call("getElementById", "parseButton")
	parseButton.Call("addEventListener", "click", cb)

	keepAlive()
}

type inBrowserRunner struct{}

func newInBrowserRunner() runner.Runner {
	return inBrowserRunner{}
}

func (r inBrowserRunner) Alive() bool {
	return true
}

func (r inBrowserRunner) Kill() error {
	return fmt.Errorf("cannot kill inbrowser runner")
}

func (r inBrowserRunner) Connection() net.Conn {
	return nil
}

func (r inBrowserRunner) IsMultithreaded() bool {
	return false
}

func (r inBrowserRunner) Pid() int {
	return -1
}

func (r inBrowserRunner) ExecuteAndGetStatus(m *gauge_messages.Message) *gauge_messages.ProtoExecutionResult {
	fmt.Printf("received %s\n", m.MessageType)
	return &gauge_messages.ProtoExecutionResult{}
}

func (r inBrowserRunner) ExecuteMessageWithTimeout(m *gauge_messages.Message) (*gauge_messages.Message, error) {
	return &gauge_messages.Message{}, nil
}
