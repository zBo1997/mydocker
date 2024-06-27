package cgroups

import (
	"mydocker/cgroups/subsystems"

	logrus "github.com/sirupsen/logrus"
)

// 设置一个Cgroup的管理器
type CgroupManager struct {
	Path    string // cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	resouce *subsystems.ResourceConfig
}

// 构造器 返回一个CgroupManager的指针
func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 批量设置对应pid进程的cgroup资源限制
func (c *CgroupManager) Apply(pid int) error {
	for _, subsystemsIns := range subsystems.SubsystemsIns {
		err := subsystemsIns.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply subsystem:%s err:%s", subsystemsIns.Name(), err)
		}
	}
	return nil
}

// 这里才是批量设置资源限制的方法
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		err := subSysIns.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("Set subsystem:%s err:%s", subSysIns.Name(), err)
		}
	}
	return nil
}

// Destroy 释放cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
