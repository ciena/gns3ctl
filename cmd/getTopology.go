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
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Topology struct {
	Nodes []*gns3.Node `json:"nodes" yaml:"nodes"`
	Links []*gns3.Link `json:"links" yaml:"links"`
}

// getNodesCmd represents the getNodes command
//
//nolint:exhaustruct
var getTopologyCmd = &cobra.Command{
	Use:     "topology [flags]",
	Aliases: []string{"to", "topo"},
	Short:   "Query the topology (nodes and links) of a GNS3 network",
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

		// Fetch nodes and links
		nodes, err := ctl.Nodes(project.ProjectId).List()
		if err != nil {
			return fmt.Errorf("unable to retrieve nodes: %w", err)
		}
		links, err := ctl.Links(project.ProjectId).List()
		if err != nil {
			return fmt.Errorf("unable to retrieve links: %w", err)
		}

		topo := Topology{
			Nodes: nodes,
			Links: links,
		}

		switch output {
		case "json":
			j, _ := json.Marshal(topo)
			fmt.Println(string(j))
		default:
			fallthrough
		case "yaml":
			y, _ := yaml.Marshal(topo)
			fmt.Println(string(y))
		case "template":
			funcMap := template.FuncMap{
				"toLower": strings.ToLower,
				"toUpper": strings.ToUpper,
				"filterLinks": func(nodeId string, links []*gns3.Link) []*gns3.Link {
					var filtered []*gns3.Link

					for _, l := range links {
						for _, n := range l.Nodes {
							if n.NodeId == nodeId {
								filtered = append(filtered, l)
								break
							}
						}
					}
					return filtered
				},
				"lookupNode": func(nodes []*gns3.Node, id string) *gns3.Node {
					for _, n := range nodes {
						if n.NodeId == id {
							return n
						}
					}
					return nil
				},
				"getPort": func(node *gns3.Node, adapterNum, portNum int) *gns3.Port {
					for _, p := range node.Ports {
						if p.AdapterNumber == adapterNum {
							return p
						}
					}
					return nil
				},
			}
			ut, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).ParseFiles(templateFile)
			if err != nil {
				return fmt.Errorf("failed to parse nodes template files '%s': %w", templateFile, err)
			}
			err = ut.Execute(os.Stdout, topo)
			if err != nil {
				return fmt.Errorf("failed to process nodes template: %w", err)
			}
		}

		return nil
	},
}

func init() {
	getCmd.AddCommand(getTopologyCmd)
	getTopologyCmd.Flags().StringP("output", "o", "columns", "Output format. One of json, yaml, columns, template=FILE")
}
