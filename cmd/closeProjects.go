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

// closeProjectCmd represents the closeProject command
//
//nolint:exhaustruct
var closeProjectsCmd = &cobra.Command{
	Use:     "projects [flags] PROJECT [PROJECT...]",
	Aliases: []string{"project", "proj", "pr"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Closes the named GNS3 projects",
	Long: `
Closes the list or specified projects. A project can be specified as either
the name of the project or as its UUID.
`,
	Run: func(cmd *cobra.Command, args []string) {
		projects := gns3.Connect().Projects()
		for _, id := range args {
			p, err := projects.Get(id)
			if err != nil {
				fmt.Printf("ERROR: %s: %v\n", id, err)
			} else {
				_, err = projects.Close(p.ProjectId)
				if err == nil {
					fmt.Println(p.ProjectId)
				} else {
					fmt.Printf("ERROR: %s: %v\n", id, err)
				}
			}
		}
	},
}

func init() {
	closeCmd.AddCommand(closeProjectsCmd)
}
