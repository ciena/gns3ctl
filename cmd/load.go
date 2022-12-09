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
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ciena/gns3ctl/pkg/gns3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	vpcsStartupTempl = `# Startup Configuration
set pcname %s
ip %s %s %s
`
)

type linkRef struct {
	nodeId  string
	port    int
	adapter int
}

type linkIndex struct {
	aEnd, zEnd linkRef
}

type linkState struct {
	suspend bool
	linkId  string
}

type linkInfo struct {
	aEndName, zEndName string
}

// loadCmd represents the load command
//
//nolint:exhaustruct
var loadCmd = &cobra.Command{
	Use:     "load [flags] FILE [FILE...]",
	Aliases: []string{"apply"},
	Short:   "Loads a project into the GNS3 environment",
	Long: `
From a YAML formated description of a network, this comamnd crate a GNS3
project, the nodes in that project, and the links between the nodes, according
to the specification in the YAML network document.

The name of the project is specified in the YAML network document, and if a
project with this name already exists, then the loading of the YAML network
document will report and error.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, filename := range args {
			pr, err := doLoad(filename)
			if err == nil {
				fmt.Println(pr.ProjectId)
			} else {
				fmt.Printf("ERROR: %s: %v\n", filename, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

func doLoad(filename string) (*gns3.Project, error) {
	var network gns3.Network
	// Open and parse the file as YAML
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&network)
	if err != nil {
		return nil, err
	}

	ctl := gns3.Connect()
	var project *gns3.Project
	project, err = ctl.Projects().Get(network.Metadata.Name)
	if err != nil {
		project, err = ctl.Projects().Create(&gns3.Project{Name: network.Metadata.Name})
		if err != nil {
			return nil, fmt.Errorf("create: %w", err)
		}
		fmt.Printf("PROJECT: %s (%s) created\n", project.Name, project.ProjectId)
	} else {
		fmt.Printf("PROJECT: %s (%s) exists\n", project.Name, project.ProjectId)
	}

	if len(network.Spec.Appliances) > 0 {
		apps := ctl.Appliances()
		templates := ctl.Templates()
		var reader io.ReadCloser
		for _, ref := range network.Spec.Appliances {
			// Attempt to parse as url
			u, err := url.Parse(ref)
			if err != nil || u.Scheme == "" {
				// not able to parse the URL, or contains
				// no schema, then assume a local file reference
				file, err := os.Open(ref)
				if err != nil {
					return nil, fmt.Errorf("ERROR: '%s': %w", ref, err)
				}
				reader = file
			} else if err != nil {
				return nil, fmt.Errorf("ERROR: unable to parse appliance reference '%s': %w", ref, err)
			} else {
				if strings.HasPrefix(strings.ToLower(u.Scheme), "http") {
					resp, err := http.Get(u.String())
					if err != nil {
						return nil, fmt.Errorf("ERROR: unable to fetch appliance '%s': %w", ref, err)
					}
					reader = resp.Body
				} else {
					return nil, fmt.Errorf("ERROR: unsupported appliance URL '%s'", ref)
				}
			}
			defer reader.Close()
			a, t, err := apps.Import(reader, path.Dir(ref))
			if err != nil {
				return nil, fmt.Errorf("ERROR: '%s': %w\n", ref, err)
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
						return nil, fmt.Errorf("ERROR: creating template '%s': %w\n", t.Name, err)
					} else {
						fmt.Println("Template", created.Name, "type", created.TemplateType, "created")
					}
				}
			}
		}
	}

	nctl := ctl.Nodes(project.ProjectId)
	for _, node := range network.Spec.Nodes {
		computeID := node.ComputeId
		if computeID == "" {
			computeID = viper.GetString("compute")
		}
		resp, err := nctl.Get(node.Name)
		if err != nil {
			if node.Template != "" {
				var t *gns3.Template
				t, err = ctl.Templates().Get(node.Template)
				if err != nil {
					return nil, fmt.Errorf("unknown template: %w", err)
				}
				resp, err = nctl.CreateUsingTemplate(&gns3.Node{
					Name:      node.Name,
					NodeType:  node.Type,
					ComputeId: computeID,
					X:         node.X,
					Y:         node.Y,
					Z:         node.Z}, t)
			} else {
				symbol := ""
				switch strings.ToLower(node.Type) {
				case gns3.TypeVpcs:
					symbol = gns3.SymbolVpcs
				case gns3.TypeNat:
					symbol = gns3.SymbolCloud
				case gns3.TypeRouter:
					symbol = gns3.SymbolRouter
				case gns3.TypeEthernetSwitch:
					symbol = gns3.SymbolEthernetSwitch
				case gns3.TypeMultilayerSwitch:
					symbol = gns3.SymbolMultilayerSwitch
				case gns3.TypeFirewall:
					symbol = gns3.SymbolFirewall
				default:
					symbol = fmt.Sprintf(":/symbols/classic/%s.svg", strings.ToLower(node.Type))
				}
				resp, err = nctl.Create(&gns3.Node{
					Name:      node.Name,
					NodeType:  node.Type,
					ComputeId: computeID,
					Symbol:    symbol,
					X:         node.X,
					Y:         node.Y,
					Z:         node.Z})
			}
			if err != nil {
				return nil, fmt.Errorf("node create: %w", err)
			}
			switch strings.ToLower(node.Type) {
			case gns3.TypeVpcs:
				if node.Config != nil {
					// If we are a VPCS and have a config, we will write out a
					// startup.vpc files
					filename := fmt.Sprintf("%s/projects/%s/project-files/vpcs/%s/startup.vpc",
						viper.GetString("base-directory"), project.ProjectId, resp.NodeId)
					err := os.MkdirAll(path.Dir(filename), 0755)
					if err != nil {
						return nil, fmt.Errorf("create VPCS startup file: %w", err)
					}
					data := fmt.Sprintf(vpcsStartupTempl, node.Config.Name, node.Config.Address, node.Config.Netmask, node.Config.Gateway)
					err = os.WriteFile(filename, []byte(data), 0644)
					fmt.Printf("STARTUP: %s\n", filename)
					if err != nil {
						return nil, fmt.Errorf("create VPCS startup file: %w", err)
					}
				}
			default:
			}
			fmt.Printf("NODE: %s (%s) created\n", resp.Name, resp.NodeId)
		} else {
			fmt.Printf("NODE: %s (%s) exists\n", resp.Name, resp.NodeId)
		}
	}

	lctl := ctl.Links(project.ProjectId)
	links, err := lctl.List()
	if err != nil {
		return nil, fmt.Errorf("%w: error listing links", err)
	}

	presentLinks := make(map[linkIndex]linkState, len(links))

	for _, link := range links {
		if len(link.Nodes) != 2 {
			continue
		}

		aEnd := linkRef{
			nodeId:  link.Nodes[0].NodeId,
			port:    link.Nodes[0].PortNumber,
			adapter: link.Nodes[0].AdapterNumber,
		}

		zEnd := linkRef{
			nodeId:  link.Nodes[1].NodeId,
			port:    link.Nodes[1].PortNumber,
			adapter: link.Nodes[1].AdapterNumber,
		}

		state := linkState{suspend: link.Suspend, linkId: link.LinkId}
		presentLinks[linkIndex{aEnd: aEnd, zEnd: zEnd}] = state
	}

	createLinks := make(map[linkIndex]linkInfo, len(network.Spec.Links))

	var keep []string
	for _, link := range network.Spec.Links {
		a, err := nctl.Get(link.AEnd.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to find a-end: %w", err)
		}
		z, err := nctl.Get(link.ZEnd.Name)
		if err != nil {
			return nil, fmt.Errorf("unable to find z-end: %w", err)
		}
		aEnd := linkRef{nodeId: a.NodeId, port: link.AEnd.Port, adapter: link.AEnd.Adapter}
		zEnd := linkRef{nodeId: z.NodeId, port: link.ZEnd.Port, adapter: link.ZEnd.Adapter}
		key := linkIndex{aEnd: aEnd, zEnd: zEnd}
		keyToggle := linkIndex{aEnd: zEnd, zEnd: aEnd}

		state, ok := presentLinks[key]
		if !ok {
			state, ok = presentLinks[keyToggle]
		}

		if ok {
			if state.suspend {
				// resume this link
				keep = append(keep, state.linkId)
			}

			delete(presentLinks, key)
			delete(presentLinks, keyToggle)

			fmt.Printf("Link: %s already exists. Suspended state: %v\n", state.linkId, state.suspend)
			continue
		}

		createLinks[key] = linkInfo{aEndName: link.AEnd.Name, zEndName: link.ZEnd.Name}
	}

	// delete all invalid links not created for this project
	for _, state := range presentLinks {
		fmt.Printf("Deleting stale link: %s\n", state.linkId)
		if _, err := lctl.Delete(state.linkId); err != nil {
			fmt.Printf("Error %v deleting link %s\n", err, state.linkId)
		}
	}

	for index, inf := range createLinks {
		var resp *gns3.Link
		var e error
		resp, e = lctl.Create(&gns3.Link{ProjectId: project.ProjectId, LinkType: "ethernet", Suspend: true, Nodes: []gns3.NodeRef{
			gns3.NodeRef{NodeId: index.aEnd.nodeId, AdapterNumber: index.aEnd.adapter, PortNumber: index.aEnd.port},
			gns3.NodeRef{NodeId: index.zEnd.nodeId, AdapterNumber: index.zEnd.adapter, PortNumber: index.zEnd.port},
		}})
		if e != nil {
			return nil, fmt.Errorf("link create: %w", e)
		}

		fmt.Printf("LINK: %s (%s-%s) created\n", resp.LinkId, inf.aEndName, inf.zEndName)
		keep = append(keep, resp.LinkId)
	}

	// start all nodes
	for _, node := range network.Spec.Nodes {
		err := nctl.Start(node.Name)
		if err != nil {
			fmt.Printf("NODE: %s: start failed\n", node.Name)
		} else {
			fmt.Printf("NODE: %s: started\n", node.Name)
		}
	}

	for _, id := range keep {
		_, e := lctl.Resume(id)
		if e != nil {
			return nil, fmt.Errorf("resume failed: %w", e)
		}
	}

	return project, nil
}
