// +build linux

/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package cm

import (
	"sync"

	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type CgroupManager interface {
	// Creates and apply the cgroup configuartions on the cgroup
	Create() error
	// Destroys the cgroup
	Destroy() error
	// Update cgroup configuration
	Update(*CgroupConfig) error
}

type ResourceConfig struct {
	// Memory limit (in bytes)
	Memory int64 `json:"memory"`

	// CPU shares (relative weight vs. other containers)
	CpuShares int64 `json:"cpu_shares"`

	// CPU hardcap limit (in usecs). Allowed cpu time in a given period.
	CpuQuota int64 `json:"cpu_quota"`
}

type CgroupConfig struct {
	Name string

	// name of parent cgroup or slice
	Parent string

	// Paths represent the cgroups paths to join
	// Paths map[string]string

	// ResourceParameters contains various cgroups settings to apply
	ResourceParameters *ResourceConfig
}

type libcontainerCgroupManager struct {
	// mu      sync.Mutex
	cgroup *CgroupConfig
	// Paths   map[string]string
}

type systemdCgroupManager struct {
	mu     sync.Mutex
	cgroup *CgroupConfig
	// Paths   map[string]string
}

func NewLibcontainerCgroupManager(cgroupConfig *CgroupConfig) *libcontainerCgroupManager {
	return &libcontainerCgroupManager{
		cgroup: cgroupConfig,
	}
}

func (m *libcontainerCgroupManager) Destroy() error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	// if err := cgroups.RemovePaths(m.Paths); err != nil {
	// 	return err
	// }
	// m.Paths = make(map[string]string)
	return nil
}

func getLibcontainerResourceConfig(resourceConfig *ResourceConfig) *configs.Resources {
	resources := &configs.Resources{}
	if resourceConfig.Memory != 0 {
		resources.Memory = resourceConfig.Memory
	}
	if resourceConfig.CpuShares != 0 {
		resources.CpuShares = resourceConfig.CpuShares
	}
	if resourceConfig.CpuQuota != 0 {
		resources.CpuQuota = resourceConfig.CpuQuota
	}
	return resources
}

func (m *libcontainerCgroupManager) Update(c *CgroupConfig) error {
	resources := getLibcontainerResourceConfig(c.ResourceParameters)
	cgroupLibcontainer := &configs.Cgroup{
		Parent:    c.Parent,
		Name:      c.Name,
		Resources: resources,
	}
	cgroupManager := &fs.Manager{
		Cgroups: cgroupLibcontainer,
	}
	config := &configs.Config{
		Cgroups: cgroupLibcontainer,
	}
	if err := cgroupManager.Set(config); err != nil {
		return err
	}
	return nil
}

// Creates
func (m *libcontainerCgroupManager) Create() error {
	resources := getLibcontainerResourceConfig(m.cgroup.ResourceParameters)
	cgroupLibcontainer := &configs.Cgroup{
		Parent:    m.cgroup.Parent,
		Name:      m.cgroup.Name,
		Resources: resources,
	}
	cgroupManager := &fs.Manager{
		Cgroups: cgroupLibcontainer,
	}
	config := &configs.Config{
		Cgroups: cgroupLibcontainer,
	}
	if err := cgroupManager.Apply(0); err != nil {
		return err
	}
	if err := cgroupManager.Set(config); err != nil {
		return err
	}
	return nil
}
func NewSystemdCgroupManager(cg *CgroupConfig) *systemdCgroupManager {
	return nil
}
