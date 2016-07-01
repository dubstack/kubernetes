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

package cm

import (
	"fmt"
	"path"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/qos"
)

const (
	podCgroupNamePrefix = "pod-"
)

// podContainerManagerImpl implements podContainerManager interface.
// It is the general implementation which allows pod level container
// management if qos Cgroup is enabled.
type podContainerManagerImpl struct {
	// nodeInfo holds information about the node
	nodeInfo *api.Node
	// qosContainersInfo hold absolute paths of the top level qos containers
	qosContainersInfo QOSContainersInfo
	// subsystems holds all info of all the mounted cgroup subsytems on the node
	subsystems *cgroupSubsystems
}

// Make sure that podContainerManagerImpl implements the PodContainerManager interface
var _ PodContainerManager = &podContainerManagerImpl{}

// applyLimits sets pod cgroup resource limits
// It also updates the resource limits on top level qos containers.
func (m *podContainerManagerImpl) applyLimits(pod *api.Pod, allPods []*api.Pod) error {
	// This function will house the logic for setting the resource parameters
	// on the pod container config and updating top level qos container configs
	return nil
}

// Exists checks if the pod's cgroup already exists
func (m *podContainerManagerImpl) Exists(pod *api.Pod) bool {
	podContainerName := m.GetPodContainerName(pod)
	cm := NewCgroupManager(m.subsystems)
	return cm.Exists(podContainerName)
}

// EnsureExists takes a pod as argument and makes sure that
// pod cgroup exists if qos cgroup hierarchy flag is enabled.
// If the pod level container doesen't already exist it is created.
func (m *podContainerManagerImpl) EnsureExists(pod *api.Pod, allPods []*api.Pod) error {
	podContainerName := m.GetPodContainerName(pod)
	glog.Infof("BAJBDHJKABDJKBAKBDKJBAKBK %v %v", podContainerName)
	cm := NewCgroupManager(m.subsystems)
	// check if container already exist
	alreadyExists := m.Exists(pod)
	if !alreadyExists {
		// Create the pod container
		containerConfig := &CgroupConfig{
			Name:               podContainerName,
			ResourceParameters: &ResourceConfig{},
		}
		if err := cm.Create(containerConfig); err != nil {
			return fmt.Errorf("Failed to created container for %v : %v", podContainerName, err)
		}
	}
	// Apply appropriate resource limits on the pod container
	// Top level qos containers limits are also updated.
	if err := m.applyLimits(pod, allPods); err != nil {
		return fmt.Errorf("Failed to apply resource limits on container for %v : %v", podContainerName, err)
	}
	return nil
}

func (m *podContainerManagerImpl) GetPodContainerName(pod *api.Pod) string {
	// Get the QoS class of the pod.
	podQOS := qos.GetPodQOS(pod)
	// Get the parent QOS container name
	var parentContainer string
	switch podQOS {
	case qos.Guaranteed:
		parentContainer = m.qosContainersInfo.Guaranteed
	case qos.Burstable:
		parentContainer = m.qosContainersInfo.Burstable
	case qos.BestEffort:
		parentContainer = m.qosContainersInfo.BestEffort
	}
	podContainer := podCgroupNamePrefix + string(pod.UID)
	// Get the absolute path of container by join
	return path.Join(parentContainer, podContainer)
}

// Destroy destroys the pod container cgroup paths
func (m *podContainerManagerImpl) Destroy(pod *api.Pod) error {
	cm := NewCgroupManager(m.subsystems)
	// get pod's container name
	podContainerName := m.GetPodContainerName(pod)
	// containerConfig object stores the cgroup specifications
	containerConfig := &CgroupConfig{
		Name:               podContainerName,
		ResourceParameters: &ResourceConfig{},
	}
	if err := cm.Destroy(containerConfig); err != nil {
		return fmt.Errorf("Failed to delete cgroup paths for %v : %v", podContainerName, err)
	}
	return nil
}

// podContainerManagerNoop implements podContainerManager interface.
// It a no op implementation and basically does nothing
type podContainerManagerNoop struct {
	cgroupsRoot string
}

// Make sure that podContainerManagerStub implements the PodContainerManager interface
var _ PodContainerManager = &podContainerManagerNoop{}

func (m *podContainerManagerNoop) Exists(_ *api.Pod) bool {
	return true
}
func (m *podContainerManagerNoop) EnsureExists(_ *api.Pod, _ []*api.Pod) error {
	return nil
}

func (m *podContainerManagerNoop) GetPodContainerName(_ *api.Pod) string {
	return m.cgroupsRoot
}

// Destroy destroys the pod container cgroup paths
func (m *podContainerManagerNoop) Destroy(_ *api.Pod) error {
	return nil
}
