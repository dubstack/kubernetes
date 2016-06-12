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

package testing

import (
	"github.com/opencontainers/runc/libcontainer/configs"

	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/stretchr/testify/mock"
)

type MockLibcontainerManager struct {
	mock.Mock
}

// Make sure that MockLibcontainerManager implements cgroups.Manager interface
var _ cgroups.Manager = &MockLibcontainerManager{}

func (m *MockLibcontainerManager) Apply(pid int) error {
	args := m.Called(pid)
	return args.Error(1)
}
func (m *MockLibcontainerManager) Destroy() error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockLibcontainerManager) GetPaths() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockLibcontainerManager) GetStats() (*cgroups.Stats, error) {
	args := m.Called()
	return args.Get(0).(*cgroups.Stats), args.Error(1)
}

func (m *MockLibcontainerManager) Set(container *configs.Config) error {
	args := m.Called(container)
	return args.Error(1)
}

func (m *MockLibcontainerManager) Freeze(state configs.FreezerState) error {
	args := m.Called(state)
	return args.Error(1)
}

func (m *MockLibcontainerManager) GetPids() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockLibcontainerManager) GetAllPids() ([]int, error) {
	args := m.Called()
	return args.Get(0).([]int), args.Error(1)
}
