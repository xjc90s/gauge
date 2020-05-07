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

package cmd

import (
	"fmt"

	"github.com/getgauge/gauge/logger"
	"github.com/getgauge/gauge/template"
	"github.com/spf13/cobra"
)

var (
	templateCmd = &cobra.Command{
		Use:     "template [flags] [args]",
		Short:   "Change template configurations",
		Long:    `Change template configurations.`,
		Example: `  gauge template java getgauge/java-template`,
		Run: func(cmd *cobra.Command, args []string) {
			if templateList || machineReadable {
				text, err := template.List(machineReadable)
				if err != nil {
					logger.Fatalf(true, err.Error())
				}
				logger.Infof(true, text)
				return
			}
			if len(args) == 0 {
				exit(fmt.Errorf("Template command needs argument(s)."), cmd.UsageString())
			}
			if len(args) == 1 {
				text, err := template.Get(args[0])
				if err != nil {
					logger.Fatalf(true, err.Error())
				}
				logger.Infof(true, text)
				return
			}
			err := template.Update(args[0], args[1])
			if err != nil {
				exit(fmt.Errorf("Template URL should be a valid link of zip."), cmd.UsageString())
			}
		},
		DisableAutoGenTag: true,
	}
	templateList bool
)

func init() {
	GaugeCmd.AddCommand(templateCmd)
	templateCmd.Flags().BoolVarP(&templateList, "list", "", false, "List all template properties")
}
