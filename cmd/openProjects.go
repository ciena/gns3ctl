/*
Copyright © 2022 Ciena Corporation <info@ciena.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
)

// openProjectCmd represents the openProject command
//
//nolint:exhaustruct
var openProjectsCmd = &cobra.Command{
	Use:     "projects [flags] PROJECT [PROJECT...]",
	Aliases: []string{"project", "proj", "pr"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Open one or more projects on the GNS3 server",
	Run: func(cmd *cobra.Command, args []string) {
		ctl := gns3.Connect()
		projects := ctl.Projects()
		for _, id := range args {
			p, err := projects.Get(id)
			if err != nil {
				fmt.Printf("NOT FOUND: %s: %v\n", id, err)
				continue
			}
			err = projects.Open(p.ProjectId)
			if err != nil {
				fmt.Printf("ERROR: unable to open '%s': %v\n", id, err)
				continue
			}
			fmt.Println(p.ProjectId)
		}
	},
}

func init() {
	openCmd.AddCommand(openProjectsCmd)
}
