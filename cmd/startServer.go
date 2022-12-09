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

// startServerCmd represents the server command
//
//nolint:exhaustruct
var startServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a GNS3 network simulation server, if one is not already running.",
	Long: `If no GNS3 network simulation server process can be found,
this command wil start a server as a daemon process.`,
	Run: func(cmd *cobra.Command, args []string) {
		gnsConfig, _ := cmd.Flags().GetString("gns-config")
		p, started, err := gns3.Connect().Server().Start(gnsConfig)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}
		if !started {
			running, err := p.IsRunning()
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				return
			}
			if running {
				fmt.Printf("INFO: a server was already started as PID '%v'\n", p.Pid)
				return
			}
			fmt.Printf("ERROR: a server was already started as PID '%v', but is not running", p.Pid)
			return
		}
		fmt.Printf("PID: %v\n", p.Pid)

	},
}

func init() {
	startCmd.AddCommand(startServerCmd)

	startServerCmd.Flags().String("gns-config", "", "Configuration file to pass to server when starting")
}
