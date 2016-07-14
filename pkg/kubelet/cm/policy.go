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

import (
	"github.com/opencontainers/runc/libcontainer/configs"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/cm"
)

func CreatePodQOSPolicyMap() map[cm.QOSClass]func(api.ResourceList, api.ResourceList) *configs.Resources {
	return map[cm.QOSClass]func(api.ResourceList, api.ResourceList) *configs.Resources{
		cm.GuaranteedQOS: GuaranteedPodQOSPolicy,
		cm.BurstableQOS:  BurstablePodQOSPolicy,
		cm.BestEffortQOS: BestEffortPodQOSPolicy,
	}
}

func GuaranteedPodQOSPolicy(requests api.ResourceList, limits api.ResourceList) *configs.Resources {
	return &configs.Resources{
		CpuShares: requests.Cpu().MilliValue(),
		CpuQuota:  limits.Cpu().MilliValue(),
		Memory:    limits.Memory().Value(),
	}
}

func BurstablePodQOSPolicy(requests api.ResourceList, limits api.ResourceList) *configs.Resources {
	return &configs.Resources{
		CpuShares: requests.Cpu().MilliValue(),
		CpuQuota:  limits.Cpu().MilliValue(),
		Memory:    limits.Memory().Value(),
	}
}

func BestEffortPodQOSPolicy(requests api.ResourceList, limits api.ResourceList) *configs.Resources {
	return &configs.Resources{}
}