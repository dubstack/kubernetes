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
	"github.com/opencontainers/runc/libcontainer/configs"
	"k8s.io/kubernetes/pkg/kubelet/cm/testing"
)

func TestLibcontainerCgroupManager_Create(t *Testing.T) {
	libcontainerManagerMock := &testing.MockLibcontainerManager{}
	//Setup expectations
	cgroupConfig := &CgroupConfig{
		Name:   "foo",
		Parent: "/",
		ResourceParameters: &ResourceConfig{
			Memory:    31457280,
			CpuQuota:  1000,
			CpuShares: 2,
		},
	}
	libcontainerConfig := &configs.Config{
		Cgroups: configs.Cgroup{
			Name:   "foo",
			Parent: "/",
			Resources: &configs.Resources{
				Memory:    31457280,
				CpuQuota:  1000,
				CpuShares: 2,
			},
		},
	}
	libcontainerManagerMock.On("Set", libcontainerConfig).Return(nil)
	libcontainerManagerMock.On("Apply", 0).Return(nil)
	libcontainerCgroupManager := NewLibcontainerCgroupManager(cgroupConfig)
	libcontainerCgroupManager.Create()
	libcontainerManagerMock.AssertExpectations(t)
}

func TestLibcontainerCgroupManager_Update(t *Testing.T) {
	libcontainerManagerMock := &testing.MockLibcontainerManager{}
	//Setup expectations
	cgroupConfig := &CgroupConfig{
		Name:   "foo",
		Parent: "/",
		ResourceParameters: &ResourceConfig{
			Memory:    31457280,
			CpuQuota:  1000,
			CpuShares: 2,
		},
	}
	libcontainerConfig := &configs.Config{
		Cgroups: configs.Cgroup{
			Name:   "foo",
			Parent: "/",
			Resources: &configs.Resources{
				Memory:    31457280,
				CpuQuota:  1000,
				CpuShares: 2,
			},
		},
	}
	libcontainerManagerMock.On("Set", libcontainerConfig).Return(nil)
	libcontainerCgroupManager := NewLibcontainerCgroupManager(cgroupConfig)
	libcontainerCgroupManager.Update(cgroupConfig)
	libcontainerManagerMock.AssertExpectations(t)
}
