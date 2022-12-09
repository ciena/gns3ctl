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

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
)

const (
	clientVersion = `Client:
  Version: %s
  Commit: %s
`
	serverVersion = `Server:
  Version: %s (%s)
`
)

var (
	Version string = "unknown"
	Commit  string = "unknown"
)

// versionCmd represents the version command
//
//nolint:exhaustruct
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Args:    cobra.NoArgs,
	Short:   "Display the GNS3 server version",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := gns3.Connect().Server().Version()
		fmt.Printf(clientVersion, Version, Commit)
		if err != nil {
			fmt.Printf(serverVersion, "ERROR:", err.Error())
			os.Exit(1)
		}
		local := "local"
		if !v.Local {
			local = "not local"
		}
		fmt.Printf(serverVersion, v.Version, local)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
