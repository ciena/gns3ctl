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
	"github.com/spf13/viper"
)

// startNodeCmd represents the startNode command
//
//nolint:exhaustruct
var startLinksCmd = &cobra.Command{
	Use:     "links",
	Aliases: []string{"li", "link"},
	Short:   "Start or resume a network link",
	Run: func(cmd *cobra.Command, args []string) {
		pname := viper.GetString("project")
		if pname == "" {
			fmt.Printf("ERROR: a project context must be specified")
			return
		}
		ctl := gns3.Connect()
		project, err := ctl.Projects().Get(pname)
		if err != nil {
			fmt.Printf("ERROR: project `%s` not found\n", pname)
			return
		}

		for _, li := range args {
			uuid, err := ctl.Links(project.ProjectId).Resume(li)
			if err != nil {
				fmt.Printf("ERROR: %s: %v\n", li, err)
			} else {
				fmt.Println(uuid)
			}
		}
	},
}

func init() {
	startCmd.AddCommand(startLinksCmd)
}
