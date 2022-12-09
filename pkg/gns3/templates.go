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
	TemplatesPath = "v2/templates"
)

//nolint:tagliatelle
type Template struct {
	Builtin           bool            `json:"builtin,omitempty" yaml:"builtin"`
	Category          string          `json:"category,omitempty" yaml:"category"`
	ComputeId         string          `json:"compute_id,omitempty" yaml:"compute_id"`
	DefaultNameFormat string          `json:"default_name_format,omitempty" yaml:"default_name_format"`
	FirstPortName     string          `json:"first_port_name,omitempty" yaml:"first_port_name"`
	Name              string          `json:"name,omitempty" yaml:"name"`
	Symbol            string          `json:"symbol,omitempty" yaml:"symbol"`
	TemplateId        string          `json:"template_id,omitempty" yaml:"template_id"`
	TemplateType      string          `json:"template_type,omitempty" yaml:"template_type"`
	Usage             string          `json:"usage,omitempty" yaml:"usage"`
	Qemu              *TemplateQemu   `json:"-" yaml:"qemu"`
	Docker            *TemplateDocker `json:"-" yaml:"docker"`
}

//nolint:tagliatelle
type TemplateQemu struct {
	Options           string `json:"options,omitempty" yaml:"options"`
	KernelCommandLine string `json:"kernel_command_line,omitempty" yaml:"kernel_command_line"`
	AdapterType       string `json:"adapter_type,omitempty" yaml:"adapter_type"`
	Adapters          int    `json:"adapters,omitempty" yaml:"adapters"`
	BootPriority      string `json:"boot_priority,omitempty" yaml:"boot_priority"`
	ConsoleType       string `json:"console_type,omitempty" yaml:"console_type"`
	HdaDiskInterface  string `json:"hda_disk_interface,omitempty" yaml:"hda_disk_interface"`
	Ram               int64  `json:"ram,omitempty" yaml:"ram"`
	Path              string `json:"qemu_path,omitempty" yaml:"qemu_path"`
	PortNameFormat    string `json:"port_name_format,omitempty" yaml:"port_name_format"`
	BiosImage         string `json:"bios_image,omitempty" yaml:"bios_image"`
	HdaDiskImage      string `json:"hda_disk_image,omitempty" yaml:"hda_disk_image"`
	CdromImage        string `json:"cdrom_image,omitempty" yaml:"cdrom_image"`
}

//nolint:tagliatelle
type TemplateDocker struct {
	Image             string   `json:"image,omitempty" yaml:"image,omitempty"`
	Usage             string   `json:"usage,omitempty" yaml:"usage,omitempty"`
	Adapters          int      `json:"adapters,omitempty" yaml:"adapters,omitempty"`
	StartCommand      string   `json:"start_command,omitempty" yaml:"start_command,omitempty"`
	Environment       string   `json:"environment,omitempty" yaml:"environment,omitempty"`
	ConsoleType       string   `json:"console_type,omitempty" yaml:"console_type,omitempty"`
	ConsoleAutoStart  bool     `json:"console_auto_start,omitempty" yaml:"console_auto_start,omitempty"`
	ConsoleHttpPort   int      `json:"console_http_port,omitempty" yaml:"console_http_port,omitempty"`
	ConsoleHttpPath   string   `json:"console_http_path,omitempty" yaml:"console_http_path,omitempty"`
	ConsoleResolution string   `json:"console_resolution,omitempty" yaml:"console_resolution,omitempty"`
	ExtraHosts        string   `json:"extra_hosts,omitempty" yaml:"extra_hosts,omitempty"`
	ExtraVolumes      []string `json:"extra_volumes,omitempty" yaml:"extra_volumes,omitempty"`
	CustomAdapters    []*struct {
		Schema               string   `json:"$schema,omitempty" yaml:"$schema,omitempty"`
		Description          string   `json:"description,omitempty" yaml:"description,omitempty"`
		Type                 any      `json:"type,omitempty" yaml:"type,omitempty"`
		Properties           any      `json:"properties,omitempty" yaml:"properties,omitempty"`
		Required             []string `json:"required,omitempty" yaml:"required,omitempty"`
		AdditionalProperties bool     `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	} `json:"custom_adapters,omitempty" yaml:"custom_adapters,omitempty"`
}

type Templates struct {
	gns3 *Gns3
}

func (g *Gns3) Templates() *Templates {
	return &Templates{gns3: g}
}

func (t *Templates) List() ([]Template, error) {
	var list []Template
	err := t.gns3.Get(TemplatesPath, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (t *Templates) Get(id string) (*Template, error) {
	var template Template

	_, err := uuid.Parse(id)
	if err == nil {
		// Think this may be a UUID, so try to delete directly
		err = t.gns3.Get(fmt.Sprintf("%s/%s", TemplatesPath, id), &template)
		if err == nil {
			return &template, nil
		}
	}
	// not a UUID, so get a list of templates and search based on name
	list, err := t.List()
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

func (t *Templates) Close(id string) (string, error) {
	template, err := t.Get(id)
	if err != nil {
		return "", err
	}
	return template.TemplateId, t.gns3.Post(fmt.Sprintf("%s/%s/close", TemplatesPath, id), "", nil, nil)
}

func (t *Templates) Create(template *Template) (*Template, error) {
	var out Template
	err := t.gns3.Post(TemplatesPath, "application/json", template, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (t *Templates) Delete(id string) (string, error) {
	template, err := t.Get(id)
	if err != nil {
		return "", err
	}
	return template.TemplateId, t.gns3.Delete(fmt.Sprintf("%s/%s", TemplatesPath, template.TemplateId))
}

func (t Template) MarshalJSON() ([]byte, error) {
	var m map[string]interface{}
	type _Template Template

	res, err := json.Marshal(_Template(t))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &m)
	if err != nil {
		return nil, err
	}

	var typeData map[string]interface{}

	// for now marshal the qemu options
	if t.TemplateType == TemplateTypeQemu {
		res, err = json.Marshal(t.Qemu)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(res, &typeData)
		if err != nil {
			return nil, err
		}
	}

	if t.TemplateType == TemplateTypeDocker {
		res, err = json.Marshal(t.Docker)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(res, &typeData)
		if err != nil {
			return nil, err
		}
	}

	// add the type data for qemu to the template map
	for k, v := range typeData {
		m[k] = v
	}

	return json.Marshal(m)
}

func (t *Template) UnmarshalJSON(data []byte) error {
	type _Template Template
	var template _Template

	err := json.Unmarshal(data, &template)
	if err != nil {
		return err
	}

	*t = Template(template)

	if t.TemplateType == TemplateTypeQemu {
		t.Qemu = &TemplateQemu{}
		err = json.Unmarshal(data, &t.Qemu)
		if err != nil {
			return err
		}
	}

	if t.TemplateType == TemplateTypeDocker {
		t.Docker = &TemplateDocker{}
		err = json.Unmarshal(data, &t.Docker)
		if err != nil {
			return err
		}
	}

	return nil
}
