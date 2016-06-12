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
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"
)

// struct libcontainerCgroupManager implements CgroupManager interface
// Uses the Libcontainer raw fs cgroup manager for cgroup management
type libcontainerCgroupManager struct {
	// Libcontainer raw fs cgroup manager
	fsCgroupManager *fs.Manager
}

// Make sure that libcontainerCgroupManager implements the CgroupManager interface
var _ CgroupManager = &libcontainerCgroupManager{}

// Returns libcontainer's cgroups.config{} struct given the general cgroupConfig
func getLibcontainerCgroupConfig(cgroupConfig *CgroupConfig) *configs.Cgroup {
	resourceConfig := cgroupConfig.ResourceParameters
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
	cgroupLibcontainer := &configs.Cgroup{
		Parent:    cgroupConfig.Parent,
		Name:      cgroupConfig.Name,
		Resources: resources,
	}
	return cgroupLibcontainer
}

// Factory Method that returns a configured CgroupManager
func NewLibcontainerCgroupManager(cgroupConfig *CgroupConfig) *libcontainerCgroupManager {
	libcontainerCgroupConfig := getLibcontainerCgroupConfig(cgroupConfig)
	return &libcontainerCgroupManager{
		fsCgroupManager: &fs.Manager{
			Cgroups: libcontainerCgroupConfig,
		},
	}
}

// 'Destroy' destroys the associated cgroup set
func (m *libcontainerCgroupManager) Destroy() error {
	if err := m.fsCgroupManager.Destroy(); err != nil {
		return err
	}
	return nil
}

// 'Update' updates the cgroup set with the specified Cgroup Configuration
func (m *libcontainerCgroupManager) Update(cgroupConfig *CgroupConfig) error {
	libcontainerCgroupConfig := getLibcontainerCgroupConfig(cgroupConfig)
	m.fsCgroupManager.Cgroups = libcontainerCgroupConfig
	config := &configs.Config{
		Cgroups: m.fsCgroupManager.Cgroups,
	}
	if err := m.fsCgroupManager.Set(config); err != nil {
		return err
	}
	return nil
}

// Create creates the cgroup
func (m *libcontainerCgroupManager) Create() error {
	config := &configs.Config{
		Cgroups: m.fsCgroupManager.Cgroups,
	}
	if err := m.fsCgroupManager.Apply(0); err != nil {
		return err
	}
	if err := m.fsCgroupManager.Set(config); err != nil {
		return err
	}
	return nil
}
