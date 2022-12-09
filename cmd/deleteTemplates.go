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
	"errors"
	"fmt"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteTemplatesCmd represents the deleteTemplate command
//
//nolint:exhaustruct
var deleteTemplatesCmd = &cobra.Command{
	Use:     "template [flags] TEMPLATE [TEMPLATE...]",
	Aliases: []string{"templates", "templ", "temp", "te", "tpl"},
	Short:   "Delete an existing template",
	Long: `
Delete the list of named templates. A template can be specified either by the
name or the UUID of the template.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templates := gns3.Connect().Templates()
		for _, id := range args {
			uuid, err := templates.Delete(id)
			if err == nil {
				fmt.Println(uuid)
			} else if !errors.Is(err, gns3.ErrNotFound) || !viper.GetBool("ignore-not-found") {
				fmt.Printf("ERROR: %s: %v\n", id, err)
			}
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteTemplatesCmd)
}
