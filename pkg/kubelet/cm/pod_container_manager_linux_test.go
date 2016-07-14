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
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	kubetypes "k8s.io/kubernetes/pkg/types"
)

// getDefaultQOSContainersInfo returns the default top level QOS containers name.
func getDefaultQOSContainersInfo() QOSContainersInfo {
	return QOSContainersInfo{
		Guaranteed: "/",
		BestEffort: "/BestEffort",
		Burstable:  "/Burstable",
	}
}

// newPodWithUID creates and returns a new pod
// with the specified UID and containers
func newPodWithUID(uid string, containers []api.Container) *api.Pod {
	return &api.Pod{
		ObjectMeta: api.ObjectMeta{
			UID: kubetypes.UID(uid),
		},
		Spec: api.PodSpec{
			Containers: containers,
		},
	}
}

// getResourceConfig returns a ResourceConfig
// with the specified cpu and memory resource limits and shares
func getResourceConfig(cpuShares, cpuQuota, memoryLimit string) *ResourceConfig {
	res := &ResourceConfig{}
	if cpuShares != "" {
		cs := getMilliValue(resource.MustParse(cpuShares))
		res.CpuShares = cs
	}
	if cpuQuota != "" {
		cq := getMilliValue(resource.MustParse(cpuQuota))
		res.CpuQuota = cq
	}
	if memoryLimit != "" {
		m := getValue(resource.MustParse(memoryLimit))
		res.Memory = m
	}
	return res
}

// getConfig returns a CgroupConfig object with the
// specified cgroup configuration
func getConfig(name, cpuShares, cpuQuota, memoryLimit string) *CgroupConfig {
	return &CgroupConfig{
		Name:               name,
		ResourceParameters: getResourceConfig(cpuShares, cpuQuota, memoryLimit),
	}
}

// NewFakePodContainerManager is a factory method that
// returns a fake Pod container manager
func NewFakePodContainerManager(nodeInfo *api.Node, mockObject *mockCgroupManager) *podContainerManagerImpl {
	return &podContainerManagerImpl{
		cgroupManager:     mockObject,
		qosContainersInfo: getDefaultQOSContainersInfo(),
		qosPolicy:         CreatePodQOSPolicyMap(),
		nodeInfo:          nodeInfo,
	}
}

func TestPodContainerApplyLimits(t *testing.T) {
	nodeInfo := getNode("10", "10Gi")
	testCases := []struct {
		pod      *api.Pod
		expected *CgroupConfig
	}{
		{
			pod: newPodWithUID("guaranteed", []api.Container{
				newContainer("foo", getResourceList("100m", "100Mi"), getResourceList("100m", "100Mi")),
				newContainer("bar", getResourceList("50m", "100Mi"), getResourceList("50m", "100Mi")),
				newContainer("foobar", getResourceList("50m", "100Mi"), getResourceList("50m", "100Mi")),
			}),
			expected: getConfig("/pod-guaranteed", "200m", "200m", "300Mi"),
		},
		{
			pod: newPodWithUID("burstable", []api.Container{
				newContainer("foo", getResourceList("100m", "100Mi"), getResourceList("100m", "200Mi")),
				newContainer("bar", getResourceList("50m", "100Mi"), getResourceList("50m", "100Mi")),
				newContainer("foobar", getResourceList("50m", "100Mi"), getResourceList("100m", "100Mi")),
			}),
			expected: getConfig("/Burstable/pod-burstable", "200m", "250m", "400Mi"),
		},
		{
			pod: newPodWithUID("burstable", []api.Container{
				newContainer("foo", getResourceList("100m", "100Mi"), getResourceList("", "")),
				newContainer("bar", getResourceList("50m", "100Mi"), getResourceList("50m", "100Mi")),
				newContainer("foobar", getResourceList("50m", "100Mi"), getResourceList("", "100Mi")),
			}),
			expected: getConfig("/Burstable/pod-burstable", "200m", "10", "10Gi"),
		},
		{
			pod: newPodWithUID("besteffort", []api.Container{
				newContainer("foo", getResourceList("", ""), getResourceList("", "")),
				newContainer("bar", getResourceList("", ""), getResourceList("", "")),
			}),
			expected: getConfig("/BestEffort/pod-besteffort", "2m", "", ""),
		},
	}
	for _, tc := range testCases {
		mockObject := &mockCgroupManager{}
		// Set expectation. Check if Cgroup Manager Update() method is being called with
		// the expected CgroupConfig Object
		mockObject.On("Update", tc.expected).Return(nil)
		pm := NewFakePodContainerManager(nodeInfo, mockObject)
		pm.applyLimits(tc.pod)
		mockObject.AssertExpectations(t)
	}
}
