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

package e2e_node

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

var _ = framework.KubeDescribe("Kubelet Cgroup Manager", func() {
	f := NewDefaultFramework("kubelet-cgroup-manager")
	Describe("QOS containers", func() {
		Context("On enabling QOS cgroup hierarchy", func() {
			It("Top level QoS containers should have been created", func() {
				podName := "qos-container" + string(util.NewUUID())
				pod := &api.Pod{
					ObjectMeta: api.ObjectMeta{
						Name:      podName,
						Namespace: f.Namespace.Name,
					},
					Spec: api.PodSpec{
						// Force the Pod to schedule to the node without a scheduler running
						NodeName: *nodeName,
						// Don't restart the Pod since it is expected to exit
						RestartPolicy: api.RestartPolicyNever,
						Containers: []api.Container{
							{
								Image:   "gcr.io/google_containers/busybox:1.24",
								Name:    podName,
								Command: []string{"sh", "-c", "if [ -d /tmp/memory/Burstable ]; then echo Found; else echo Failed; fi"},
								VolumeMounts: []api.VolumeMount{
									{
										Name:      "sysfscgroup",
										MountPath: "/tmp",
									},
								},
							},
						},
						Volumes: []api.Volume{
							{
								Name: "sysfscgroup",
								VolumeSource: api.VolumeSource{
									HostPath: &api.HostPathVolumeSource{Path: "/sys/fs/cgroup"},
								},
							},
						},
					},
				}
				f.TestContainerOutput("top level qos creation", pod, 0, []string{"Found"})
			})
		})
		Context("On scheduling a ", func() {
			It("Top level QoS containers should have been created", func() {
				podName := "qos-container" + string(util.NewUUID())
				pod := &api.Pod{
					ObjectMeta: api.ObjectMeta{
						Name:      podName,
						Namespace: f.Namespace.Name,
					},
					Spec: api.PodSpec{
						// Force the Pod to schedule to the node without a scheduler running
						NodeName: *nodeName,
						// Don't restart the Pod since it is expected to exit
						RestartPolicy: api.RestartPolicyNever,
						Containers: []api.Container{
							{
								Image:   "gcr.io/google_containers/busybox:1.24",
								Name:    podName,
								Command: []string{"sh", "-c", "if [ -d /tmp/memory/Burstable ]; then echo Found; else echo Failed; fi"},
								VolumeMounts: []api.VolumeMount{
									{
										Name:      "sysfscgroup",
										MountPath: "/tmp",
									},
								},
							},
						},
						Volumes: []api.Volume{
							{
								Name: "sysfscgroup",
								VolumeSource: api.VolumeSource{
									HostPath: &api.HostPathVolumeSource{Path: "/sys/fs/cgroup"},
								},
							},
						},
					},
				}
				f.TestContainerOutput("top level qos creation", pod, 0, []string{"Found"})
			})
		})
	})
})

// newTestPodForQuota returns a pod that has the specified requests and limits
func newTestPodForQuota(f *framework.Framework, name string, requests api.ResourceList, limits api.ResourceList) *api.Pod {
	return &api.Pod{
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:  "pause",
					Image: framework.GetPauseImageName(f.Client),
					Resources: api.ResourceRequirements{
						Requests: requests,
						Limits:   limits,
					},
				},
			},
		},
	}
}
