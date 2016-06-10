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

// QosClass defines the supported QoS classes of Pods/Containers.
type QosClass string

const (
	// GuaranteedQoS is the Guaranteed QoS class.
	GuaranteedQoS = "Guaranteed"
	// BurstableQoS is the Burstable QoS class.
	BurstableQoS = "Burstable"
	// BestEffortQoS is the BestEffort QoS class.
	BestEffortQoS = "BestEffort"
)

// ResourceConfig holds information about all the supported cgroup resource parameters.
type ResourceConfig struct {
	// Memory limit (in bytes).
	Memory int64 `json:"memory"`
	// CPU shares (relative weight vs. other containers).
	CpuShares int64 `json:"cpu_shares"`
	// CPU hardcap limit (in usecs). Allowed cpu time in a given period.
	CpuQuota int64 `json:"cpu_quota"`
}

// CgroupConfig holds the cgroup configuration Information.
// This is common object which is used to specify
// cgroup information to both systemd and raw cgroup fs
// implementation of the Cgroup Manager interface.
type CgroupConfig struct {
	// Cgroup or slice Name.
	Name string
	// name of parent cgroup or slice.
	Parent string
	// ResourceParameters contains various cgroups settings to apply.
	ResourceParameters *ResourceConfig
}

// CgroupManger is an interface which defines the manager which allows for cgroup management.
// Supports Cgroup Creation , Deletion and Updates.
type CgroupManager interface {
	// Creates and apply the cgroup configuartions on the cgroup.
	Create(*CgroupConfig) error
	// Destroys the cgroup.
	Destroy() error
	// Update cgroup configuration.
	Update(*CgroupConfig) error
}
