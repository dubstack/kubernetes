// +build linux

/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package cgroupfs

import (
	"fmt"
	"path"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"
	"k8s.io/kubernetes/pkg/kubelet/cm"
)

const (
	defaultContainerRoot string = "/"
)

// qosContainerManagerImpl implements QOSContainerManager interface.
// It internally uses libcontainer's cgroup manager implementation
// for qos container creation, and updates.
type qosContainerManagerImpl struct {
	manager   cgroups.Manager
	qosConfig *cm.QOSConfig
}

// Make sure that it implements the qosContainerManager interface
var _ cm.QOSContainerManager = &qosContainerManagerImpl{}

// NewQOSContainerManager is a factory method and returns a qosContainerManager object
func NewQOSContainerManager(qosConfig *cm.QOSConfig) *qosContainerManagerImpl {
	return &qosContainerManagerImpl{
		manager:   &fs.Manager{},
		qosConfig: qosConfig,
	}
}

// Create creates the top level qos cgroup containers
// We currently create containers for the Burstable and Best Effort containers
// All guaranteed pods are nested under the RootContainer by default
// Init is called only once during kubelet bootstrapping.
func (m *qosContainerManagerImpl) Init(qosConfig *cm.QOSConfig) error {
	// The rootContainer under which all qos containers are brought up is configurable
	// and can be specified through the --cgroup-root flag.
	// We default to the system root / in case no root is specified
	rootContainerName := defaultContainerRoot
	if qosConfig.RootContainerName != "" {
		rootContainerName = qosConfig.RootContainerName
	}
	// Top level for Qos containers are created only for Burstable
	// and Best Effort classes
	qosClasses := [2]cm.QOSClass{cm.BurstableQOS, cm.BestEffortQOS}

	// Create containers for both qos classes
	for _, qosClass := range qosClasses {
		// containerConfig object stores the cgroup specifications
		containerConfig := &configs.Config{
			Cgroups: &configs.Cgroup{
				Name:      string(qosClass),
				Parent:    rootContainerName,
				Resources: &configs.Resources{},
			},
		}
		//Apply(0) is a hack to create the cgroup directories for each resource
		// subsystem. The function [cgroups.Manager.apply()] applies cgroup
		// configuration to the process with the specified pid.
		// It creates cgroup files for each subsytems and writes the pid
		// in the tasks file. We use the function to create all the required
		// cgroup files but not attach any "real" pid to the cgroup.
		if err := m.manager.Apply(0); err != nil {
			return fmt.Errorf("Failed to create cgroups for the %v qos class : %v", qosClass, err)
		}
		if err := m.manager.Set(containerConfig); err != nil {
			return fmt.Errorf("Failed to Set container Config for the %v qos class : %v", qosClass, err)
		}
	}
	return nil
}

func (m *qosContainerManagerImpl) GetContainersInfo() cm.QOSContainersInfo {
	// qosContainersInfo object stores the absolute name of the created qos containers
	qosContainersInfo := cm.QOSContainersInfo{}
	rootContainerName := defaultContainerRoot
	if m.qosConfig.RootContainerName != "" {
		rootContainerName = m.qosConfig.RootContainerName
	}
	// Guaranteed container is not created so guaranteed pods are brought
	// directly under the Root container.
	qosContainersInfo.GuaranteedContainerName = rootContainerName
	qosContainersInfo.BurstableContainerName = path.Join(rootContainerName, string(cm.BurstableQOS))
	qosContainersInfo.BestEffortContainerName = path.Join(rootContainerName, string(cm.BestEffortQOS))
	return qosContainersInfo
}

//@TODO(@dubstack)
// func (m *qosContainerManagerImpl) Update() error{}
