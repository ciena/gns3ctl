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
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// getProjectCmd represents the getProject command
//
//nolint:exhaustruct
var getProjectsCmd = &cobra.Command{
	Use:     "projects [PROJECT...]",
	Short:   "Query the projects from a GNS3 server",
	Aliases: []string{"project", "proj", "pr"},
	Run: func(cmd *cobra.Command, args []string) {
		ctl := gns3.Connect()

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		output, _ := cmd.Flags().GetString("output")
		switch output {
		default:
			fallthrough
		case "columns":
			fmt.Fprintf(tw, "%s\t%s\t%s\n", "UUID", "NAME", "STATUS")
		case "json", "yaml", "name", "id":
		}
		if len(args) == 0 {
			projects, err := ctl.Projects().List()
			if err != nil {
				panic(err)
			}
			switch output {
			default:
				fallthrough
			case "columns":
				for _, p := range projects {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", p.ProjectId, p.Name, p.Status)
				}
			case "json":
				j, _ := json.Marshal(projects)
				fmt.Println(string(j))
			case "yaml":
				y, _ := yaml.Marshal(projects)
				fmt.Println(string(y))
			case "name":
				for _, p := range projects {
					fmt.Println(p.Name)
				}
			case "id":
				for _, p := range projects {
					fmt.Println(p.ProjectId)
				}
			}
		} else {
			cvt := yaml.Marshal
			switch output {
			default:
				fallthrough
			case "columns":
				for _, id := range args {
					project, err := ctl.Projects().Get(id)
					if err != nil {
						fmt.Fprintf(tw, "%s\t%s\t%s\n", id, "", err.Error())
					} else {
						fmt.Fprintf(tw, "%s\t%s\t%s\n", project.ProjectId, project.Name, project.Status)
					}
				}
			case "json":
				cvt = json.Marshal
				fallthrough
			case "yaml":
				list := []interface{}{}
				for _, id := range args {
					template, err := ctl.Projects().Get(id)
					if err != nil {
						nf := map[string]string{
							"name":  id,
							"error": err.Error(),
						}
						list = append(list, nf)
					} else {
						list = append(list, template)
					}
				}
				if len(list) != 1 {
					output, _ := cvt(list)
					fmt.Println(string(output))
				} else {
					output, _ := cvt(list[0])
					fmt.Println(string(output))
				}
			case "name":
				for _, id := range args {
					project, err := ctl.Projects().Get(id)
					if err != nil {
						fmt.Printf("%s => %s\n", id, err.Error())
					} else {
						fmt.Println(project.Name)
					}
				}
			case "id":
				for _, id := range args {
					project, err := ctl.Projects().Get(id)
					if err != nil {
						fmt.Printf("%s => %s\n", id, err.Error())
					} else {
						fmt.Println(project.ProjectId)
					}
				}
			}
		}
		switch output {
		default:
			fallthrough
		case "columns":
			tw.Flush()
		case "json", "yaml", "name", "id":
		}
	},
}

func init() {
	getCmd.AddCommand(getProjectsCmd)
	getProjectsCmd.Flags().StringP("output", "o", "columns", "Output format. One of json yaml, columns")
}
