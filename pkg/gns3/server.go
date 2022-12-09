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
	"os"
	"os/exec"

	"github.com/shirou/gopsutil/process"
)

const (
	ServerShutdownPath = "v2/shutdown"
	ServerVersionPath  = "v2/version"
)

type Server struct {
	gns3 *Gns3
}

type ServerVersion struct {
	Local   bool   `json:"local,omitempty" yaml:"local"`
	Version string `json:"version,omitempty" yaml:"version"`
}

func (g *Gns3) Server() *Server {
	return &Server{gns3: g}
}

func (s *Server) Shutdown() error {
	return s.gns3.Post(ServerShutdownPath, "", nil, nil)
}

func getProcessFromName(name string) (*process.Process, error) {
	list, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range list {
		name, err := p.Name()
		if err != nil {
			return nil, err
		}

		if name == "gns3server" {
			return p, nil
		}
	}
	return nil, nil
}

func (s *Server) Start(configFile string) (*process.Process, bool, error) {
	// If there is a gns3server process already, then don't attempt to start a new one
	p, err := getProcessFromName("gns3server")
	if err != nil {
		return nil, false, fmt.Errorf("find process by name: %w", err)
	}
	if p != nil {
		return p, false, nil
	}

	// Attempt to start a new gns3server
	devnull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	if err != nil {
		return nil, false, fmt.Errorf("open file: %w", err)
	}
	cmd := []string{"gns3server", "--daemon", "--quiet", "--local", "--allow"}
	if configFile != "" {
		cmd = append(cmd, "--config", configFile)
	}
	osCmd := exec.Command(cmd[0], cmd[1:]...)
	osCmd.Stdout = devnull
	osCmd.Stderr = devnull
	osCmd.Stdin = devnull
	err = osCmd.Run()
	if err != nil {
		return nil, false, fmt.Errorf("run %w", err)
	}

	p, err = getProcessFromName("gns3server")
	if err != nil {
		return nil, false, fmt.Errorf("find after start %w", err)
	}
	return p, true, nil
}

func (s *Server) Version() (*ServerVersion, error) {
	var serverVersion ServerVersion
	err := s.gns3.Get(ServerVersionPath, &serverVersion)
	if err == nil {
		return &serverVersion, nil
	}
	return nil, err
}
