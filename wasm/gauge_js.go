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

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"syscall/js"

	"github.com/getgauge/gauge/execution/event"
	"github.com/getgauge/gauge/execution/item"
	"github.com/getgauge/gauge/formatter"
	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/gauge_messages"
	"github.com/getgauge/gauge/parser"
	"github.com/getgauge/gauge/plugin"
	"github.com/getgauge/gauge/resolver"
	"github.com/getgauge/gauge/result"
	"github.com/getgauge/gauge/runner"
	"github.com/getgauge/gauge/validation"
)

var signal = make(chan int)

func keepAlive() {
	for {
		<-signal
	}
}

func parse() *gauge.Specification {
	document := js.Global().Get("document")
	specEl := document.Call("getElementById", "specText")
	content := specEl.Get("innerText").String()
	s, r := new(parser.SpecParser).ParseSpecText(content, "browser")
	if !r.Ok {
		for _, e := range r.Errors() {
			js.Global().Get("console").Call("error", e)
		}
		return nil
	}
	return s
}

func validate(s *gauge.Specification, r runner.Runner) {
	js.Global().Call("parse")
	specs := []*gauge.Specification{s}
	vErrs := validation.NewValidator(specs, r, gauge.NewConceptDictionary()).Validate()
	for _, e := range vErrs[s] {
		js.Global().Get("console").Call("error", e.Error())
	}
}

func setReporter(wg *sync.WaitGroup) {
	ch := make(chan event.ExecutionEvent, 0)
	event.Register(ch, event.SuiteStart, event.SpecStart, event.SpecEnd, event.ScenarioStart, event.ScenarioEnd, event.StepStart, event.StepEnd, event.ConceptStart, event.ConceptEnd, event.SuiteEnd)
	wg.Add(1)

	go func() {
		for {
			e := <-ch
			switch e.Topic {
			case event.SpecStart:
				fmt.Printf("# %s\n", (*e.Item.(*gauge.Specification)).Heading.Value)
			case event.ScenarioStart:
				if e.Result.(*result.ScenarioResult).ProtoScenario.GetExecutionStatus() == gauge_messages.ExecutionStatus_SKIPPED {
					continue
				}
				sce := e.Item.(*gauge.Scenario)
				if sce.SpecDataTableRow.GetRowCount() != 0 {
					fmt.Println(formatter.FormatTable(&sce.SpecDataTableRow))
				}
				if sce.ScenarioDataTableRow.GetRowCount() != 0 {
					fmt.Println(formatter.FormatTable(&sce.ScenarioDataTableRow))
				}
				fmt.Printf("## %s\n", sce.Heading.Value)
			case event.StepEnd:
				step := e.Item.(gauge.Step)
				stepRes := e.Result.(*result.StepResult)
				if stepRes.GetStepFailed() {
					fmt.Printf("Failed Step: %s\nLine No: %d\nError Message: %s\nStacktrace: %s\n",
						step.LineText,
						step.LineNo,
						stepRes.ProtoStepExecResult().GetExecutionResult().GetErrorMessage(),
						stepRes.ProtoStepExecResult().GetExecutionResult().GetStackTrace())
				}
			case event.SpecEnd:
				wg.Done()
				return
			}
		}
	}()
}

func main() {
	document := js.Global().Get("document")
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		console := document.Call("getElementById", "out")
		console.Set("innerHTML", "")
		s := parse()
		r := newInBrowserRunner()
		validate(s, r)
		resolver.GetResolvedDataTablerows(s.DataTable.Table)
		item.MaxRetriesCount = 1
		event.InitRegistry()
		wg := &sync.WaitGroup{}
		setReporter(wg)
		event.Notify(event.NewExecutionEvent(event.SpecStart, s, nil, 1, gauge_messages.ExecutionInfo{}))
		eRes := item.NewSpecExecutor(s, r, &plugin.GaugePlugins{}, gauge.NewBuildErrors(), 1).Execute(false, true, false)
		event.Notify(event.NewExecutionEvent(event.SpecEnd, s, nil, 1, gauge_messages.ExecutionInfo{}))
		wg.Wait()
		executed := eRes.ScenarioCount - eRes.ScenarioSkippedCount
		passed := executed - eRes.ScenarioFailedCount
		fmt.Printf("Scenarios:\t%d executed\t%d passed\t%d failed\t%d skipped\n", executed, passed, eRes.ScenarioFailedCount, eRes.ScenarioSkippedCount)
		return nil
	})

	runButton := document.Call("getElementById", "runButton")
	runButton.Call("addEventListener", "click", cb)

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
	type row struct {
		Cells []string `json:"cells"`
	}
	
	type table struct {
		Headers []string   `json:"headers"`
		Rows    []row `json:"rows"`
	}

	type param struct {
		Kind  string `json:"kind"`
		Value string `json:"value"`
		Table table  `json:"table"`
	}

	switch m.MessageType {
	case gauge_messages.Message_ScenarioExecutionStarting,
		gauge_messages.Message_ScenarioExecutionEnding,
		gauge_messages.Message_StepExecutionStarting,
		gauge_messages.Message_StepExecutionEnding:
		return &gauge_messages.ProtoExecutionResult{}
	case gauge_messages.Message_ExecuteStep:
		params := []param{}
		for _, p := range m.ExecuteStepRequest.Parameters {
			if p.ParameterType == gauge_messages.Parameter_Table {
				t := table{Headers: p.Table.Headers.Cells, Rows: make([]row, 0)}
				for _, r := range p.Table.Rows {
					t.Rows = append(t.Rows, row{Cells: r.Cells})
				}
				params = append(params, param{Kind: "table", Table: t})
			} else {
				params = append(params, param{Kind: "static", Value: p.Value})
			}
		}
		b, err := json.Marshal(params)
		if err != nil {
			js.Global().Get("console").Call("error", err.Error())
			return &gauge_messages.ProtoExecutionResult{Failed: true, ErrorMessage: err.Error(), ErrorType: gauge_messages.ProtoExecutionResult_VERIFICATION}
		}
		jsErr := js.Global().Call("execute", m.ExecuteStepRequest.ParsedStepText, string(b))
		if jsErr != js.Undefined() {
			return &gauge_messages.ProtoExecutionResult{Failed: true, ErrorMessage: jsErr.String(), ErrorType: gauge_messages.ProtoExecutionResult_VERIFICATION}
		}
		return &gauge_messages.ProtoExecutionResult{}
	}
	return nil
}

func (r inBrowserRunner) ExecuteMessageWithTimeout(m *gauge_messages.Message) (*gauge_messages.Message, error) {
	if m.MessageType == gauge_messages.Message_StepValidateRequest {
		implemented := js.Global().Call("stepImplemented", m.StepValidateRequest.StepText).Bool()
		if !implemented {
			return &gauge_messages.Message{
				MessageId:   m.MessageId,
				MessageType: gauge_messages.Message_StepValidateResponse,
				StepValidateResponse: &gauge_messages.StepValidateResponse{
					ErrorType: gauge_messages.StepValidateResponse_STEP_IMPLEMENTATION_NOT_FOUND,
				},
			}, nil
		}
		impl := js.Global().Call("stepImplementationLocations", m.StepValidateRequest.StepText).String()
		impls := strings.Split(impl, "|")
		if len(impls) > 1 {
			return &gauge_messages.Message{
				MessageId:   m.MessageId,
				MessageType: gauge_messages.Message_StepValidateResponse,
				StepValidateResponse: &gauge_messages.StepValidateResponse{
					ErrorType:    gauge_messages.StepValidateResponse_DUPLICATE_STEP_IMPLEMENTATION,
					ErrorMessage: fmt.Sprintf("Duplicate step implementations at: %s", impls),
				},
			}, nil

		}
	}
	return &gauge_messages.Message{
		MessageId:            m.MessageId,
		MessageType:          gauge_messages.Message_StepValidateResponse,
		StepValidateResponse: &gauge_messages.StepValidateResponse{IsValid: true},
	}, nil

}
