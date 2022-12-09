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

	"github.com/google/uuid"
)

const (
	ProjectsPath     = "v2/projects"
	ProjectPath      = "v2/projects/%s"
	ProjectOpenPath  = "v2/projects/%s/open"
	ProjectClosePath = "v2/projects/%s/close"
)

//nolint:tagliatelle
type Project struct {
	AutoClose           bool   `json:"auto_close,omitempty"`
	AutoOpen            bool   `json:"auto_open,omitempty"`
	AutoStart           bool   `json:"auto_start,omitempty"`
	DrawingGridSize     int    `json:"drawing_grid_size,omitempty"`
	Filename            string `json:"filename,omitempty"`
	GridSize            int    `json:"grid_size,omitempty"`
	Name                string `json:"name,omitempty"`
	Path                string `json:"path,omitempty"`
	ProjectId           string `json:"project_id,omitempty"`
	SceneHeight         int    `json:"scene_height,omitempty"`
	SceneWidth          int    `json:"scene_width,omitempty"`
	ShowGrid            bool   `json:"show_grid,omitempty"`
	ShowInterfaceLabels bool   `json:"show_interface_labels,omitempty"`
	ShowLayers          bool   `json:"show_layers,omitempty"`
	SnapToGrid          bool   `json:"snap_to_grid,omitempty"`
	Status              string `json:"status,omitempty"`
	Zoom                int    `json:"zoom,omitempty"`
}

type Projects struct {
	gns3 *Gns3
}

func (g *Gns3) Projects() *Projects {
	return &Projects{gns3: g}
}

func (p *Projects) List() ([]Project, error) {
	list := []Project{}
	err := p.gns3.Get(ProjectsPath, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *Projects) Get(id string) (*Project, error) {
	_, err := uuid.Parse(id)
	var project Project
	if err == nil {
		// Think this may be a UUID, so try to delete directly
		err = p.gns3.Get(fmt.Sprintf("%s/%s", ProjectsPath, id), &project)
		if err == nil {
			return &project, nil
		}
	}

	// not a UUID, so get a list of projects and search based on name
	list, err := p.List()
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

func (p *Projects) Create(project *Project) (*Project, error) {
	var out Project
	err := p.gns3.Post(ProjectsPath, "application/json", project, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (p *Projects) Delete(id string) (string, error) {
	project, err := p.Get(id)
	if err != nil {
		return "", err
	}
	return project.ProjectId, p.gns3.Delete(fmt.Sprintf(ProjectPath, project.ProjectId))
}

func (p *Projects) Close(id string) (string, error) {
	project, err := p.Get(id)
	if err != nil {
		return "", err
	}
	return project.ProjectId, p.gns3.Post(fmt.Sprintf(ProjectClosePath, project.ProjectId), "", nil, nil)
}

func (p *Projects) Open(id string) error {
	project, err := p.Get(id)
	if err != nil {
		return err
	}
	return p.gns3.Post(fmt.Sprintf(ProjectOpenPath, project.ProjectId), "", nil, nil)
}
