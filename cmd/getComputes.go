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

// getComputeCmd represents the getCompute command
//
//nolint:exhaustruct
var getComputesCmd = &cobra.Command{
	Use:     "computes [flags] [COMPUTE...]",
	Short:   "Query the GNS3 compute nodes",
	Aliases: []string{"co", "comp", "compute"},
	Run: func(cmd *cobra.Command, args []string) {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		output, _ := cmd.Flags().GetString("output")
		switch output {
		default:
			fallthrough
		case "columns":
			fmt.Fprintf(tw, "%s\t%s\t%s\n", "UUID", "NAME", "HOST")
		case "json":
		case "yaml":
		}
		if len(args) == 0 {
			computes, err := gns3.Connect().Computes().List()
			if err != nil {
				panic(err)
			}
			output, _ = cmd.Flags().GetString("output")
			switch output {
			default:
				fallthrough
			case "columns":
				for _, p := range computes {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", p.ComputeId, p.Name, p.Host)
				}
			case "json":
				j, _ := json.Marshal(computes)
				fmt.Println(string(j))
			case "yaml":
				y, _ := yaml.Marshal(computes)
				fmt.Println(string(y))
			}
		} else {
			output, _ = cmd.Flags().GetString("output")
			cvt := yaml.Marshal
			switch output {
			default:
				fallthrough
			case "columns":
				for _, id := range args {
					compute, err := gns3.Connect().Computes().Get(id)
					if err != nil {
						fmt.Fprintf(tw, "%s\t%s\t%s\n", id, "", err.Error())
					} else {
						fmt.Fprintf(tw, "%s\t%s\t%s\n", compute.ComputeId, compute.Name, compute.Host)
					}
				}
			case "json":
				cvt = json.Marshal
				fallthrough
			case "yaml":
				list := []interface{}{}
				for _, id := range args {
					template, err := gns3.Connect().Templates().Get(id)
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
		output, _ = cmd.Flags().GetString("output")
		switch output {
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
	getCmd.AddCommand(getComputesCmd)
	getComputesCmd.Flags().StringP("output", "o", "columns", "Output format. One of json yaml, columns")
}
