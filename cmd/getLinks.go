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
	"strings"
	"text/tabwriter"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	LinksPath = "v2/projects/%s/links"
)

// getLinksCmd represents the getLinks command
//
//nolint:exhaustruct
var getLinksCmd = &cobra.Command{
	Use:     "links [flags] [LINK...]",
	Aliases: []string{"li", "link"},
	Short:   "Query a GNS3 server network links",
	Run: func(cmd *cobra.Command, args []string) {
		ctl := gns3.Connect()

		pname := viper.GetString("project")
		if pname == "" {
			fmt.Printf("ERROR: a project context must be specified")
			return
		}
		project, err := ctl.Projects().Get(pname)
		if err != nil {
			fmt.Printf("ERROR: project `%s` not found\n", pname)
			return
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		output, _ := cmd.Flags().GetString("output")
		switch output {
		default:
			fallthrough
		case "columns":
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", "UUID", "TYPE", "SUSPEND", "NODES")
		case "json", "yaml", "name", "id":
		}
		if len(args) == 0 {
			links, err := ctl.Links(project.ProjectId).List()
			if err != nil {
				panic(err)
			}
			switch output {
			default:
				fallthrough
			case "columns":
				for _, l := range links {
					var nodes []string
					for _, n := range l.Nodes {
						info, err := ctl.Nodes(project.ProjectId).Get(n.NodeId)
						if err == nil {
							nodes = append(nodes, fmt.Sprintf("%s(%d)", info.Name, n.PortNumber))
						}
					}
					fmt.Fprintf(tw, "%s\t%s\t%t\t%s\n", l.LinkId, l.LinkType, l.Suspend, strings.Join(nodes, ","))
				}
			case "json":
				j, _ := json.Marshal(links)
				fmt.Println(string(j))
			case "yaml":
				y, _ := yaml.Marshal(links)
				fmt.Println(string(y))
			case "id", "name":
				for _, l := range links {
					fmt.Println(l.LinkId)
				}
			}
		} else {
			cvt := yaml.Marshal
			switch output {
			default:
				fallthrough
			case "columns":
				for _, id := range args {
					l, err := ctl.Links(project.ProjectId).Get(id)
					if err != nil {
						fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", id, "", "", err.Error())
					} else {
						var nodes []string
						for _, n := range l.Nodes {
							info, err := ctl.Nodes(project.ProjectId).Get(n.NodeId)
							if err == nil {
								nodes = append(nodes, fmt.Sprintf("%s(%d)", info.Name, n.PortNumber))
							}
						}
						fmt.Fprintf(tw, "%s\t%s\t%t\t%s\n", l.LinkId, l.LinkType, l.Suspend, strings.Join(nodes, ","))
					}
				}
			case "name", "id":
				for _, id := range args {
					l, err := ctl.Links(project.ProjectId).Get(id)
					if err != nil {
						fmt.Printf("%s => %s\n", id, err.Error())
					} else {
						fmt.Println(l.LinkId)
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
					out, _ := cvt(list)
					fmt.Println(string(out))
				} else {
					out, _ := cvt(list[0])
					fmt.Println(string(out))
				}
			}
		}
		switch output {
		default:
			fallthrough
		case "columns":
			tw.Flush()
		case "json", "yaml", "name", "id":
		}
	},
}

func init() {
	getCmd.AddCommand(getLinksCmd)
	getLinksCmd.Flags().StringP("output", "o", "columns", "Output format. One of json yaml, columns")
}
