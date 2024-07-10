package subsystems

import (
	"fmt"
	"mydocker/constant"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MemorySubSystem struct {
}

func (s *MemorySubSystem) Name() string {
	return "memory"
}

// Set 设置cgroupPath对应的cgroup的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	logrus.Info("memeory res:", res.MemoryLimit)
	if res.MemoryLimit == "" {
		return nil
	}
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, true)
	logrus.Info("subsysCgroupPath:", subsysCgroupPath)
	if err != nil {
		return err
	}
	// 设置这个cgroup的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes 文件中。
	if err := os.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), constant.Perm0644); err != nil {
		logrus.Error("write error:", err)
		return fmt.Errorf("set cgroup memory fail %v", err)
	}
	logrus.Info("cgroup set success:", res.MemoryLimit)
	return nil
}

// Apply 将pid加入到cgroupPath对应的cgroup中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return errors.Wrapf(err, "get cgroup %s", cgroupPath)
	}
	if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), constant.Perm0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}

// Remove 删除cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := getCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		logrus.Error("cgroup remove fail:", cgroupPath)
		return err
	}
	//删处则个目录下得所有文件
	return os.RemoveAll(subsysCgroupPath)
}
