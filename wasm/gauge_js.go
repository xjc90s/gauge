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
	"syscall/js"

	"github.com/getgauge/gauge/parser"
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
	_, r := new(parser.SpecParser).ParseSpecText(content, "browser")
	if !r.Ok {
		for _, e := range r.Errors() {
			fmt.Println(e)
		}
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
