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
	ComputesPath = "v2/computes"
	ComputePath  = "v2/computes/%s"
)

//nolint:tagliatelle
type Compute struct {
	Capabilities       map[string]interface{} `json:"capabilities,omitempty" yaml:"capabilities"`
	ComputeId          string                 `json:"compute_id,omitempty" yaml:"compute_id"`
	Connected          bool                   `json:"connected,omitempty" yaml:"connected"`
	CpuUsagePercent    float64                `json:"cpu_usage_percent,omitempty" yaml:"cpu_usage_percent"`
	Host               string                 `json:"host,omitempty" yaml:"host"`
	LastError          string                 `json:"last_error,omitempty" yaml:"last_error"`
	MemoryUsagePercent float64                `json:"memory_usage_percent,omitempty" yaml:"memory_usage_percent"`
	Name               string                 `json:"name,omitempty" yaml:"name"`
	Port               int                    `json:"port,omitempty" yaml:"port"`
	Protocol           string                 `json:"protocol,omitempty" yaml:"protocol"`
	User               string                 `json:"user,omitempty" yaml:"user"`
}

type Computes struct {
	gns3 *Gns3
}

func (g *Gns3) Computes() *Computes {
	return &Computes{gns3: g}
}

func (c *Computes) List() ([]Compute, error) {
	list := []Compute{}
	err := c.gns3.Get(ComputesPath, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Computes) Get(id string) (*Compute, error) {
	var compute Compute
	// Think this may be a UUID, so try to delete directly
	err := c.gns3.Get(fmt.Sprintf(ComputePath, id), &compute)
	if err == nil {
		return &compute, nil
	}

	// not a UUID, so get a list of computes and search based on name
	list, err := c.List()
	if err != nil {
		return nil, err
	}
	for _, val := range list {
		if val.Name == id {
			return &val, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Computes) Create(compute *Compute) (*Compute, error) {
	var out Compute
	err := c.gns3.Post(ComputesPath, "application/json", compute, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Computes) Delete(id string) (string, error) {
	compute, err := c.Get(id)
	if err != nil {
		return "", err
	}
	return compute.ComputeId, c.gns3.Delete(fmt.Sprintf("%s/%s", ComputesPath, compute.ComputeId))
}
