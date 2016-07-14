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
	"github.com/stretchr/testify/mock"
)

// mockCgroupManager is a mock object which implements the cm.CgroupManager interface
type mockCgroupManager struct {
	mock.Mock
}

// Make sure that mockLibcontainerManager implements CgroupManager interface
var _ CgroupManager = &mockCgroupManager{}

func (m *mockCgroupManager) Exists(name string) bool {
	args := m.Called(name)
	return args.Bool(0)
}

func (m *mockCgroupManager) Update(cgroupConfig *CgroupConfig) error {
	args := m.Called(cgroupConfig)
	return args.Error(0)
}
func (m *mockCgroupManager) Destroy(cgroupConfig *CgroupConfig) error {
	args := m.Called(cgroupConfig)
	return args.Error(0)
}

func (m *mockCgroupManager) Create(cgroupConfig *CgroupConfig) error {
	args := m.Called(cgroupConfig)
	return args.Error(0)
}
