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
	"os"
	"path"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
)

// importApplianceCmd represents the importAppliance command
//
//nolint:exhaustruct
var importApplianceCmd = &cobra.Command{
	Use:     "appliances",
	Aliases: []string{"ap", "app", "appliance"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Import appliance definitions from files",
	Run: func(cmd *cobra.Command, args []string) {
		ctl := gns3.Connect()
		apps := ctl.Appliances()
		templates := ctl.Templates()
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Printf("ERROR: '%s': %v\n", filename, err)
				continue
			}
			a, t, err := apps.Import(file, path.Dir(filename))
			if err != nil {
				fmt.Printf("ERROR: '%s': %v\n", filename, err)
			} else {
				fmt.Println(a.Name)
				// try creating the template for the appliance
				existing, err := templates.Get(t.Name)
				if err == nil {
					fmt.Println("Template", existing.Name, "already present")
				} else {
					// create the template
					created, err := templates.Create(t)
					if err != nil {
						fmt.Printf("ERROR: creating template '%s': %v\n", t.Name, err)
					} else {
						fmt.Println("Template", created.Name, "type", created.TemplateType, "created")
					}
				}
			}
			file.Close()

		}
	},
}

func init() {
	importCmd.AddCommand(importApplianceCmd)
}
