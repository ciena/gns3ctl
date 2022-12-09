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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
//
//nolint:exhaustruct
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm", "remove"},
	Short:   "Deletes a subresource",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().Bool("ignore-not-found", false, "ignore and don't report not found errors")
	_ = viper.BindPFlag("ignore-not-found", deleteCmd.PersistentFlags().Lookup("ignore-not-found"))
}
