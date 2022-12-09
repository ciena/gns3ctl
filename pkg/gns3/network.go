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

//nolint:tagliatelle
type Network struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name string `json:"name" yaml:"name"`
	} `json:"metadata" yaml:"metadata"`
	Spec struct {
		Appliances []string `json:"appliances,omitempty" yaml:"appliances,omitempty"`
		Nodes      []struct {
			Name      string `json:"name,omitempty" yaml:"name"`
			Type      string `json:"type,omitempty" yaml:"type,omitempty"`
			Template  string `json:"template,omitempty" yaml:"template,omitempty"`
			ComputeId string `json:"compute_id,omitempty" yaml:"compute_id"`
			X         int    `json:"x,omitempty" yaml:"x"`
			Y         int    `json:"y,omitempty" yaml:"y"`
			Z         int    `json:"z,omitempty" yaml:"z"`
			Config    *struct {
				Name    string `json:"name,omitempty" yaml:"name"`
				Address string `json:"address,omitempty" yaml:"address"`
				Netmask string `json:"netmask,omitempty" yaml:"netmask"`
				Gateway string `json:"gateway,omitempty" yaml:"gateway"`
			} `json:"config,omitempty" yaml:"config"`
		} `json:"nodes" yaml:"nodes"`
		Links []struct {
			AEnd struct {
				Name    string `json:"name,omitempty" yaml:"name"`
				Adapter int    `json:"adapter,omitempty" yaml:"adapter"`
				Port    int    `json:"port,omitempty" yaml:"port"`
			} `json:"aEnd,omitempty" yaml:"aEnd"`
			ZEnd struct {
				Name    string `json:"name,omitempty" yaml:"name"`
				Adapter int    `json:"adapter,omitempty" yaml:"adapter"`
				Port    int    `json:"port,omitempty" yaml:"port"`
			} `json:"zEnd,omitempty" yaml:"zEnd"`
		} `json:"links" yaml:"links"`
	} `json:"spec" yaml:"spec"`
}
