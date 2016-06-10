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

// QOSClass defines the supported qos classes of Pods/Containers.
type QOSClass string

//TODO(dubstack) Move this to qos package use a common definition
// in both qos and cm package
const (
	// GuaranteedQOS is the Guaranteed qos class.
	GuaranteedQOS QOSClass = "Guaranteed"
	// BurstableQOS is the Burstable qos class.
	BurstableQOS QOSClass = "Burstable"
	// BestEffortQOS is the BestEffort qos class.
	BestEffortQOS QOSClass = "BestEffort"
)

// QOSConfig defines how the qos cgroup hierarchy is organized
type QOSConfig struct {
	// RootContainerName is the root for the qos hierarchy
	RootContainerName string
}

// QOSContainersInfo defines the names of containers per qos
type QOSContainersInfo struct {
	GuaranteedContainerName string
	BestEffortContainerName string
	BurstableContainerName  string
}

// QOSContainerManager stores and manages top level qos containers
// Currently its used to create top level qos containers for the Burstable
// and Best Effort qos classes during kubelet bootstrapping.
// We don't have a separate qos class for the Guaranteed qos.
// Guaranteed pods are nested under the RootContainer by default.
type QOSContainerManager interface {
	// Init interfaces with host OS to define
	// qos specific containers per convention
	Init(config *QOSConfig) error

	// @TODO(dubstack) support updates to the Qos containers
	// Update() error

	// GetContainerInfo returns the top level qos containers names.
	// The top level qos container names are different depending
	// upon the driver (systemd or cgroupfs)
	GetContainersInfo() QOSContainersInfo
}
