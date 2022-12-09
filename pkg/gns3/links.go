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
	"fmt"
)

const (
	LinksPath = "v2/projects/%s/links"
	LinkPath  = "v2/projects/%s/links/%s"
)

type LinkLabel struct {
	Style string `json:"style,omitempty" yaml:"style"`
	Text  string `json:"text,omitempty" yaml:"text"`
	X     int    `json:"x,omitempty" yaml:"x"`
	Y     int    `json:"y,omitempty" yaml:"y"`
}

//nolint:tagliatelle
type NodeRef struct {
	AdapterNumber int        `json:"adapter_number" yaml:"adapter_number"`
	Label         *LinkLabel `json:"label,omitempty" yaml:"label,omitempty"`
	NodeId        string     `json:"node_id,omitempty" yaml:"node_id"`
	PortNumber    int        `json:"port_number" yaml:"port_number"`
}

//nolint:tagliatelle
type Link struct {
	CaptureComputeId string                 `json:"capture_compute_id,omitempty" yaml:"capture_compute_id"`
	CaptureFileName  string                 `json:"capture_file_name,omitempty" yaml:"capture_file_name"`
	CaptureFilePath  string                 `json:"capture_file_path,omitempty" yaml:"capture_file_path"`
	Capturing        bool                   `json:"capturing,omitempty" yaml:"capturing"`
	Filters          map[string]interface{} `json:"filters,omitempty" yaml:"filters"`
	LinkId           string                 `json:"link_id,omitempty" yaml:"link_id"`
	LinkType         string                 `json:"link_type,omitempty" yaml:"link_type"`
	Nodes            []NodeRef              `json:"nodes,omitempty" yaml:"nodes"`
	ProjectId        string                 `json:"project_id,omitempty" yaml:"project_id"`
	Suspend          bool                   `json:"suspend,omitempty" yaml:"suspend"`
}

type Links struct {
	gns3      *Gns3
	projectID string
}

func (g *Gns3) Links(id string) *Links {
	return &Links{gns3: g, projectID: id}
}

func (l *Links) List() ([]*Link, error) {
	list := []*Link{}
	err := l.gns3.Get(fmt.Sprintf(LinksPath, l.projectID), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (l *Links) Get(id string) (*Link, error) {
	var link Link
	err := l.gns3.Get(fmt.Sprintf(LinkPath, l.projectID, id), &link)
	return &link, err
}

func (l *Links) Create(link *Link) (*Link, error) {
	var out Link
	err := l.gns3.Post(fmt.Sprintf(LinksPath, l.projectID), "application/json", link, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (l *Links) Delete(id string) (string, error) {
	li, err := l.Get(id)
	if err != nil {
		return "", err
	}
	return li.LinkId, l.gns3.Delete(fmt.Sprintf(LinkPath, l.projectID, li.LinkId))
}

var resumePatch = map[string]interface{}{
	"suspend": false,
}
var suspendPatch = map[string]interface{}{
	"suspend": true,
}

func (l *Links) Resume(id string) (string, error) {
	li, err := l.Get(id)
	if err != nil {
		return "", err
	}
	return li.LinkId, l.gns3.Put(fmt.Sprintf(LinkPath, l.projectID, li.LinkId), "application/json", &resumePatch, nil)
}

func (l *Links) Suspend(id string) (string, error) {
	li, err := l.Get(id)
	if err != nil {
		return "", err
	}
	return li.LinkId, l.gns3.Put(fmt.Sprintf(LinkPath, l.projectID, li.LinkId), "application/json", &suspendPatch, nil)
}
