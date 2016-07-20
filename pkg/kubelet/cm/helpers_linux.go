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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"

	"github.com/golang/glog"
	libcontainercgroups "github.com/opencontainers/runc/libcontainer/cgroups"
)

// GetCgroupSubsystems returns information about the mounted cgroup subsystems
func GetCgroupSubsystems() (*CgroupSubsystems, error) {
	// get all cgroup mounts.
	allCgroups, err := libcontainercgroups.GetCgroupMounts()
	if err != nil {
		return &CgroupSubsystems{}, err
	}
	if len(allCgroups) == 0 {
		return &CgroupSubsystems{}, fmt.Errorf("failed to find cgroup mounts")
	}

	mountPoints := make(map[string]string, len(allCgroups))
	for _, mount := range allCgroups {
		for _, subsystem := range mount.Subsystems {
			mountPoints[subsystem] = mount.Mountpoint
		}
	}
	return &CgroupSubsystems{
		Mounts:      allCgroups,
		MountPoints: mountPoints,
	}, nil
}

// readProcsFile takes a cgroup directory name as an argument
// reads through the cgroup's procs file and returns a list of tgid's.
// It returns an empty list if a procs file doesn't exists
func readProcsFile(dir string) ([]int, error) {
	procsFile := filepath.Join(dir, "cgroup.procs")
	_, err := os.Stat(procsFile)
	if os.IsNotExist(err) {
		// The procsFile does not exist, So no pids attached to this directory
		return []int{}, nil
	}
	f, err := os.Open(procsFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		s   = bufio.NewScanner(f)
		out = []int{}
	)

	for s.Scan() {
		if t := s.Text(); t != "" {
			pid, err := strconv.Atoi(t)
			if err != nil {
				return nil, err
			}
			out = append(out, pid)
		}
	}
	if len(out) != 0 {
		glog.V(3).Infof("XOXOXOXOXOXOXO found pid in dir %v", dir)
	}
	return out, nil
}

// GetPodResourceRequest returns the pod requests for the supported resources.
// Pod request is the summation of resource requests of all containers in the pod.
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
// Pod limit is the summation of resource limits of all containers
// in the pod. If limit for a particular resource is not specified for
// even a single container then we return the node resource capacity
// as the pod limit for the particular resource.
func GetPodResourceLimits(pod *api.Pod, nodeInfo *api.Node) api.ResourceList {
	capacity := nodeInfo.Status.Capacity
	glog.Infof("The capacity: %v", capacity)
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
