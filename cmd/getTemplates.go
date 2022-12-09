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

// getTemplateCmd represents the getTemplate command
//
//nolint:exhaustruct
var getTemplatesCmd = &cobra.Command{
	Use:     "templates [flags] [TEMPLATE...]",
	Short:   "Query templates from the GNS3 server",
	Aliases: []string{"t", "te", "temp", "temps", "template"},
	Run: func(cmd *cobra.Command, args []string) {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		val, _ := cmd.Flags().GetString("output")
		switch val {
		default:
			fallthrough
		case "columns":
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", "UUID", "NAME", "CATEGORY", "TYPE", "BUILTIN")
		case "json":
		case "yaml":
		}
		ctl := gns3.Connect()
		if len(args) == 0 {
			templates, err := ctl.Templates().List()
			if err != nil {
				panic(err)
			}
			val, _ := cmd.Flags().GetString("output")
			switch val {
			default:
				fallthrough
			case "columns":
				for _, t := range templates {
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\n", t.TemplateId, t.Name, t.Category, t.TemplateType, t.Builtin)
				}
			case "json":
				j, _ := json.Marshal(templates)
				fmt.Println(string(j))
			case "yaml":
				y, _ := yaml.Marshal(templates)
				fmt.Println(string(y))
			}
		} else {
			val, _ := cmd.Flags().GetString("output")
			cvt := yaml.Marshal
			switch val {
			default:
				fallthrough
			case "columns":
				for _, id := range args {
					template, err := ctl.Templates().Get(id)
					if err != nil {
						fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", id, "", "", "", err.Error())
					} else {
						fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\n", template.TemplateId, template.Name, template.Category, template.TemplateType, template.Builtin)
					}
				}
			case "json":
				cvt = json.Marshal
				fallthrough
			case "yaml":
				list := []interface{}{}
				for _, id := range args {
					template, err := ctl.Templates().Get(id)
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
			}
		}
		val, _ = cmd.Flags().GetString("output")
		switch val {
		default:
			fallthrough
		case "columns":
			tw.Flush()
		case "json":
		case "yaml":
		}
	},
}

func init() {
	getCmd.AddCommand(getTemplatesCmd)
	getTemplatesCmd.Flags().StringP("output", "o", "columns", "Output format. One of json, yaml, columns")
}
