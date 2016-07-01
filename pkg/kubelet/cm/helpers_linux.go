/*
Copyright 2016 The Kubernetes Authors.

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
	"fmt"

	"github.com/golang/glog"

	libcontainercgroups "github.com/opencontainers/runc/libcontainer/cgroups"
)

// cgroupSubsystems holds information about the mounted cgroup subsytems
type cgroupSubsystems struct {
	// Cgroup subsystem mounts.
	// e.g.: "/sys/fs/cgroup/cpu" -> ["cpu", "cpuacct"]
	mounts []libcontainercgroups.Mount

	// Cgroup subsystem to their mount location.
	// e.g.: "cpu" -> "/sys/fs/cgroup/cpu"
	mountPoints map[string]string
}

// GetCgroupSubsystems returns information about the mounted cgroup subsystems
func getCgroupSubsystems() (*cgroupSubsystems, error) {
	// Get all cgroup mounts.
	allCgroups, err := libcontainercgroups.GetCgroupMounts()
	if err != nil {
		return &cgroupSubsystems{}, fmt.Errorf("Failed to get the cgroup mounts on the system: %v", err)
	}
	if len(allCgroups) == 0 {
		return &cgroupSubsystems{}, fmt.Errorf("Failed to find the cgroup mounts")
	}

	// subsystems, err := libcontainercgroups.GetAllSubsystems()
	// if err != nil {
	// 	return &cgroupSubsystems{}, fmt.Errorf("Failed to get all subsystems supported by the kernel: %v", err)
	// }
	// subsystemsMap := make(map[string]bool, len(subsystems))
	// for _, subsystem := range subsystems {
	// 	subsystemsMap[subsystem] = true
	// }

	mountPoints := make(map[string]string, len(allCgroups))
	for _, mount := range allCgroups {
		for _, subsystem := range mount.Subsystems {
			// if exists, _ := subsystemsMap[subsystem]; exists {
			mountPoints[subsystem] = mount.Mountpoint
			// }
		}
	}
	glog.Infof("%v", mountPoints)
	return &cgroupSubsystems{
		mounts:      allCgroups,
		mountPoints: mountPoints,
	}, nil
}
