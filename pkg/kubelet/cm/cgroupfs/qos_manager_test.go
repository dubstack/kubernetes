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
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/kubelet/cm"
)

// TestGetContainersInfo
func TestGetContainersInfo(t *testing.T) {
	cases := []struct {
		qosConfig     *cm.QOSConfig
		expectedValue cm.QOSContainersInfo
	}{
		{
			qosConfig: &cm.QOSConfig{
				RootContainerName: "",
			},
			expectedValue: cm.QOSContainersInfo{
				GuaranteedContainerName: "/",
				BurstableContainerName:  "/Burstable",
				BestEffortContainerName: "/BestEffort",
			},
		},
		{
			qosConfig: &cm.QOSConfig{
				RootContainerName: "/root-container",
			},
			expectedValue: cm.QOSContainersInfo{
				GuaranteedContainerName: "/root-container",
				BurstableContainerName:  "/root-container/Burstable",
				BestEffortContainerName: "/root-container/BestEffort",
			},
		},
	}
	as := assert.New(t)
	for idx, tc := range cases {
		c := NewQOSContainerManager(tc.qosConfig)
		actual := c.GetContainersInfo()
		as.Equal(tc.expectedValue, actual, "expected test case [%d] to return %q; got %q instead", idx, tc.expectedValue, actual)
	}
}
