// +build linux
package cm

import (
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"
)

type CgroupManager interface {
	// Create and set the Cgroup
	Create() error
	// Destroys the cgroup set
	Destroy() error

	// Update configuration
	Update(c *configs.Cgroup) error
}

type CgroupManagerFs struct {
	// mu      sync.Mutex
	Cgroups *configs.Cgroup
	// Paths   map[string]string
}

func NewCgroupManagerFs(cg *configs.Cgroup) *CgroupManagerFs {
	return &CgroupManagerFs{
		Cgroups: cg,
	}
}

func (m *CgroupManagerFs) Destroy() error {
	// if m.Cgroups.Paths != nil {
	// 	return nil
	// }
	// m.mu.Lock()
	// defer m.mu.Unlock()
	// if err := cgroups.RemovePaths(m.Paths); err != nil {
	// 	return err
	// }sudo s
	// m.Paths = make(map[string]string)
	return nil
}

func (m *CgroupManagerFs) Update(c *configs.Cgroup) error {
	cg := &fs.Manager{
		Cgroups: c,
	}
	fakeConfig := &configs.Config{
		Cgroups: c,
	}
	if err := cg.Set(fakeConfig); err != nil {
		return err
	}
	return nil
}

func (m *CgroupManagerFs) Create() error {
	cg := &fs.Manager{
		Cgroups: m.Cgroups,
	}
	fakeConfig := &configs.Config{
		Cgroups: m.Cgroups,
	}
	if err := cg.Apply(0); err != nil {
		return err
	}
	if err := cg.Set(fakeConfig); err != nil {
		return err
	}
	return nil
}
func NewCgroupManagerSystemd(cg *configs.Cgroup) *CgroupManagerFs {
	return nil
}
