/*
Copyright 2022 Ciena Corporation

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

package gns3

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	NodesPath = "v2/projects/%s/nodes"
	NodePath  = "v2/projects/%s/nodes/%s"
)

//nolint:tagliatelle
type Port struct {
	AdapterNumber int                    `json:"adapter_number,omitempty" yaml:"adapter_number,omitempty"`
	DataLinkTypes map[string]interface{} `json:"data_link_types,omitempty" yaml:"data_link_types,omitempty"`
	LinkType      string                 `json:"link_type,omitempty" yaml:"link_type,omitempty"`
	Name          string                 `json:"name,omitempty" yaml:"name,omitempty"`
	PortNumber    int                    `json:"port_number,omitempty" yaml:"port_number,omitempty"`
	ShortName     string                 `json:"short_name,omitempty" yaml:"short_name,omitempty"`
	MacAddress    string                 `json:"mac_address,omitempty" yaml:"mac_address,omitempty"`
	IpAddress     string                 `json:"ip_address,omitempty" yaml:"ip_address,omitempty"`
}

//nolint:tagliatelle
type Node struct {
	CommandLine      string                   `json:"command_line,omitempty" yaml:"command_line,omitempty"`
	ComputeId        string                   `json:"compute_id,omitempty" yaml:"compute_id,omitempty"`
	Console          int                      `json:"console,omitempty" yaml:"console,omitempty"`
	ConsoleAutoStart bool                     `json:"console_auto_start,omitempty" yaml:"console_auto_start,omitempty"`
	ConsoleHost      string                   `json:"console_host,omitempty" yaml:"console_host,omitempty"`
	ConsoleType      string                   `json:"console_type,omitempty" yaml:"console_type,omitempty"`
	CustomAdapters   []map[string]interface{} `json:"custom_adapters,omitempty" yaml:"custom_adapters,omitempty"`
	FirstPortName    string                   `json:"first_port_name,omitempty" yaml:"first_port_name,omitempty"`
	Height           int                      `json:"height,omitempty" yaml:"height,omitempty"`
	Label            map[string]interface{}   `json:"label,omitempty" yaml:"label,omitempty"`
	Locked           bool                     `json:"locked,omitempty" yaml:"locked,omitempty"`
	Name             string                   `json:"name,omitempty" yaml:"name,omitempty"`
	NodeDirectory    string                   `json:"node_directory,omitempty" yaml:"node_directory,omitempty"`
	NodeId           string                   `json:"node_id,omitempty" yaml:"node_id,omitempty"`
	NodeType         string                   `json:"node_type,omitempty" yaml:"node_type,omitempty"`
	PortNameFormat   string                   `json:"port_name_format,omitempty" yaml:"port_name_format,omitempty"`
	PortSegmentSize  int                      `json:"port_segment_size,omitempty" yaml:"port_segment_size,omitempty"`
	Ports            []*Port                  `json:"ports,omitempty" yaml:"ports,omitempty"`
	ProjectId        string                   `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	Properties       map[string]interface{}   `json:"properties,omitempty" yaml:"properties,omitempty"`
	Status           string                   `json:"status,omitempty" yaml:"status,omitempty"`
	Symbol           string                   `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	TemplateId       string                   `json:"template_id,omitempty" yaml:"template_id,omitempty"`
	Width            int                      `json:"width,omitempty" yaml:"width,omitempty"`
	X                int                      `json:"x,omitempty" yaml:"x,omitempty"`
	Y                int                      `json:"y,omitempty" yaml:"y,omitempty"`
	Z                int                      `json:"z,omitempty" yaml:"z,omitempty"`
}

type Nodes struct {
	gns3      *Gns3
	projectID string
}

func (g *Gns3) Nodes(id string) *Nodes {
	return &Nodes{gns3: g, projectID: id}
}

func (n *Nodes) List() ([]*Node, error) {
	list := []*Node{}
	err := n.gns3.Get(fmt.Sprintf(NodesPath, n.projectID), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (n *Nodes) Get(id string) (*Node, error) {
	_, err := uuid.Parse(id)
	var node Node
	if err == nil {
		// Think this may be a UUID, so try to delete directly
		err = n.gns3.Get(fmt.Sprintf(NodePath, n.projectID, id), &node)
		if err == nil {
			return &node, nil
		}
	}
	// not a UUID, so get a list of projects and search based on name
	list, err := n.List()
	if err != nil {
		return nil, err
	}
	for _, val := range list {
		if val.Name == id {
			return val, nil
		}
	}
	return nil, ErrNotFound
}

func (n *Nodes) Create(node *Node) (*Node, error) {
	var out Node
	err := n.gns3.Post(fmt.Sprintf(NodesPath, n.projectID), "application/json", node, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (n *Nodes) CreateUsingTemplate(node *Node, template *Template) (*Node, error) {
	in := *node
	in.NodeType = template.TemplateType
	in.Symbol = template.Symbol
	//in := Node{Name: name, ComputeId: template.ComputeId, NodeType: template.TemplateType}
	if template.TemplateType == TemplateTypeQemu && template.Qemu != nil {
		// fill the node properties
		qemuData, err := json.Marshal(template.Qemu)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(qemuData, &in.Properties)
		if err != nil {
			return nil, err
		}
	}
	if template.TemplateType == TemplateTypeDocker && template.Docker != nil {
		// fill the node properties
		dockerData, err := json.Marshal(template.Docker)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(dockerData, &in.Properties)
		if err != nil {
			return nil, err
		}
	}

	return n.Create(&in)
}

func (n *Nodes) Start(id string) error {
	no, err := n.Get(id)
	if err != nil {
		return err
	}
	return n.gns3.Post(fmt.Sprintf(NodePath+"/start", n.projectID, no.NodeId), "application/json", nil, nil)
}

func (n *Nodes) Stop(id string) error {
	no, err := n.Get(id)
	if err != nil {
		return err
	}
	var node Node
	fmt.Printf(NodePath+"/stop\n", n.projectID, no.NodeId)
	err = n.gns3.Post(fmt.Sprintf(NodePath+"/stop", n.projectID, no.NodeId), "application/json", nil, &node)
	return err
}
