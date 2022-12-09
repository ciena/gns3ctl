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

// getApplianceCmd represents the getAppliance command
//
//nolint:exhaustruct
var getAppliancesCmd = &cobra.Command{
	Use:     "appliances [flags] [APPLIANCE...]",
	Short:   "Query the GNS3 server appliances",
	Aliases: []string{"ap", "app", "appliance"},
	Run: func(cmd *cobra.Command, args []string) {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		val, _ := cmd.Flags().GetString("output")
		switch val {
		default:
			fallthrough
		case "columns":
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", "NAME", "CATEGORY", "PRODUCTNAME", "VENDOR", "BUILTIN", "STATUS")
		case "json":
		case "yaml":
		}
		if len(args) == 0 {
			appliances, err := gns3.Connect().Appliances().List()
			if err != nil {
				panic(err)
			}
			val, _ := cmd.Flags().GetString("output")
			switch val {
			default:
				fallthrough
			case "columns":
				for _, a := range appliances {
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\t%s\n", a.Name, a.Category, a.ProductName, a.VendorName, a.Builtin, a.Status)
				}
			case "json":
				j, _ := json.Marshal(appliances)
				fmt.Println(string(j))
			case "yaml":
				y, _ := yaml.Marshal(appliances)
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
					a, err := gns3.Connect().Appliances().Get(id)
					if err != nil {
						fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", id, "", "", "", "", err.Error())
					} else {
						fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%t\t%s\n", a.Name, a.Category, a.ProductName, a.VendorName, a.Builtin, a.Status)
					}
				}
			case "json":
				cvt = json.Marshal
				fallthrough
			case "yaml":
				list := []interface{}{}
				for _, id := range args {
					appliance, err := gns3.Connect().Appliances().Get(id)
					if err != nil {
						nf := map[string]string{
							"name":  id,
							"error": err.Error(),
						}
						list = append(list, nf)
					} else {
						list = append(list, appliance)
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
	getCmd.AddCommand(getAppliancesCmd)
	getAppliancesCmd.Flags().StringP("output", "o", "columns", "Output format. One of json, yaml, columns")
}
