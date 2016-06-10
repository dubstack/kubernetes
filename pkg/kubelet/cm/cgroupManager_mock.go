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

import (
	"github.com/stretchr/testify/mock"
)

// MockCgroupManager is a mock object which implements the cm.CgroupManager interface
type MockCgroupManager struct {
	mock.Mock
}

// Make sure that MockLibcontainerManager implements CgroupManager interface
var _ CgroupManager = &MockCgroupManager{}

func (m *MockCgroupManager) Update(cgroupConfig *CgroupConfig) error {
	args := m.Called(cgroupConfig)
	return args.Error(0)
}
func (m *MockCgroupManager) Destroy() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCgroupManager) Create(cgroupConfig *CgroupConfig) error {
	args := m.Called(cgroupConfig)
	return args.Error(0)
}
