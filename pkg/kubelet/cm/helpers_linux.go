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

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"

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

// GetPodResourceRequest returns the pod requests for the supported resources
// Pod request is defined as the summation of resource requests of all containers
// in the pod.
func GetPodResourceRequests(pod *api.Pod) api.ResourceList {
	requests := api.ResourceList{}
	zeroQuantity := resource.MustParse("0")
	for _, resource := range []api.ResourceName{api.ResourceCPU, api.ResourceMemory} {
		for _, container := range pod.Spec.Containers {
			quantity := container.Resources.Requests[resource]
			if quantity.Cmp(zeroQuantity) == 1 {
				delta := quantity.Copy()
				if _, exists := requests[resource]; !exists {
					requests[resource] = *delta
				} else {
					delta.Add(requests[resource])
					requests[resource] = *delta
				}
			}
		}
	}
	return requests
}

// GetPodResourceLimits returns the pod limits for the supported resources
// Pod limit is defined as the summation of resource limits of all containers
// in the pod. If limit for a particular resource is not specified for
// even a single container then we return the node resource capacity
// as the pod limit for the particular resource.
func GetPodResourceLimits(pod *api.Pod, nodeInfo *api.Node) api.ResourceList {
	capacity := nodeInfo.Status.Capacity
	limits := api.ResourceList{}
	zeroQuantity := resource.MustParse("0")
	for _, resource := range []api.ResourceName{api.ResourceCPU, api.ResourceMemory} {
		for _, container := range pod.Spec.Containers {
			quantity := container.Resources.Limits[resource]
			if quantity.Cmp(zeroQuantity) == 1 {
				delta := quantity.Copy()
				if _, exists := limits[resource]; !exists {
					limits[resource] = *delta
				} else {
					delta.Add(limits[resource])
					limits[resource] = *delta
				}
			} else {
				// if limit not specified for a particular resource in a container
				// we default the pod resource limit to the resource capacity of the node
				if cap, exists := capacity[resource]; exists {
					limits[resource] = *cap.Copy()
					break
				}
			}
		}
	}
	return limits
}
