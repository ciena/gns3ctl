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
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	NodesPath = "v2/projects/%s/nodes"
	virshTmpl = `virsh net-dhcp-leases default --mac {{.MacAddress}} | tail -2 | head -1 | awk '{print $5}' | sed -e 's;/.*;;'`
)

var ErrNoTemplateFile = errors.New("a template file must be specified")
var ErrNoProjectSpecified = errors.New("a project must be specified")
var ErrProjectNotFound = errors.New("project not found")

// getNodesCmd represents the getNodes command
//
//nolint:exhaustruct
var getNodesCmd = &cobra.Command{
	Use:     "nodes [flags] [NODE...]",
	Aliases: []string{"no", "node"},
	Short:   "Query the nodes of a GNS3 network",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If the output option template was specified then make sure
		// an actual template file was part of the option
		var templateFile string
		output, _ := cmd.Flags().GetString("output")
		if strings.HasPrefix(output, "template") {
			parts := strings.SplitN(output, "=", 2)
			if len(parts) != 2 {
				return ErrNoTemplateFile
			}
			output = parts[0]
			templateFile = parts[1]
		}

		// No project, no nodes
		pname := viper.GetString("project")
		if pname == "" {
			return ErrNoProjectSpecified
		}

		ctl := gns3.Connect()
		project, err := ctl.Projects().Get(pname)
		if err != nil {
			return fmt.Errorf("project '%s' not found: %w", pname, err)
		}
		var nodes []*gns3.Node
		var errs []error

		// Fetch nodes, either whole list or by ID/name
		if len(args) == 0 {
			nodes, err = ctl.Nodes(project.ProjectId).List()
			if err != nil {
				return fmt.Errorf("unable to retrieve nodes: %w", err)
			}
		} else {
			for _, id := range args {
				n, err := ctl.Nodes(project.ProjectId).Get(id)

				// If a query for a named node fails, save the error to
				// report at the end of the output, mimics kubectl
				if err != nil {
					errs = append(errs, fmt.Errorf("node '%s': %w", id, err))
				} else {
					nodes = append(nodes, n)
				}
			}
		}

		// The nodes as they come from GNS3 don't container the IP address for the
		// interfaces on the nodes. This is actually useful information and so we
		// augment the node with this information if we have a command to fetch
		// the IP address.
		var ipTmpl string
		ipTmpl, _ = cmd.Flags().GetString("get-ip-command")
		switch ipTmpl {
		case "virsh":
			ipTmpl = virshTmpl
		case "none":
			ipTmpl = ""
		default:
		}

		// If a template command was specified, then parse it and augment the default
		// nodes record(s) from GNS3 with IP address information
		if ipTmpl != "" {
			ut, err := template.New("ip-lookup").Parse(ipTmpl)
			if err != nil {
				return fmt.Errorf("unable to parse IP address query template '%s': %w", ipTmpl, err)
			}
			for _, n := range nodes {
				for _, p := range n.Ports {
					var b strings.Builder
					err = ut.Execute(&b, p)
					if err != nil {
						return fmt.Errorf("failed to process IP address query template for '%s/%s': %w", n.Name, p.Name, err)
					}
					shCmd, _ := cmd.Flags().GetString("shell-command")
					cmd := strings.Fields(shCmd)
					cmd = append(cmd, b.String())
					out, err := exec.Command(cmd[0], cmd[1:]...).Output()
					if err != nil {
						return fmt.Errorf("failed to execute IP address query command for '%s/%s': %w", n.Name, p.Name, err)
					}
					p.IpAddress = strings.TrimSpace(string(out))
				}
			}
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		switch output {
		default:
			fallthrough
		case "columns":
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", "UUID", "NAME", "PORT", "TYPE", "STATUS")
			for _, n := range nodes {
				fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n", n.NodeId, n.Name, n.Console, n.NodeType, n.Status)
			}
			tw.Flush()
		case "json":
			j, _ := json.Marshal(nodes)
			fmt.Println(string(j))
		case "yaml":
			y, _ := yaml.Marshal(nodes)
			fmt.Println(string(y))
		case "id":
			for _, n := range nodes {
				fmt.Println(n.NodeId)
			}
		case "name":
			for _, n := range nodes {
				fmt.Println(n.Name)
			}
		case "template":
			ut, err := template.ParseFiles(templateFile)
			if err != nil {
				return fmt.Errorf("failed to parse nodes template files '%s': %w", templateFile, err)
			}
			err = ut.Execute(os.Stdout, nodes)
			if err != nil {
				return fmt.Errorf("failed to process nodes template: %w", err)
			}
		}

		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "Error from server: %s\n", err)
			}
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getNodesCmd)
	getNodesCmd.Flags().String("get-ip-command", "virsh", "command template to query IP address for nodes. One of none, virsh, CUSTOM")
	getNodesCmd.Flags().String("shell-command", "sh -c", "shell command ")
	getNodesCmd.Flags().StringP("output", "o", "columns", "Output format. One of json, yaml, columns, template=FILE")
}
