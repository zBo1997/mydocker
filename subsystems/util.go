package subsystems

import (
	"bufio"
	"errors"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

const mountPointIndex = 4

/*
* 获取group的绝对路径
 */
func getCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	//获取
	cgroupRoot := findCgroupMountpoint(subsystem)
	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}
	return "", errors.New("")
}

/*
findCgroupMountpoint 通过/proc/self/mountinfo找出挂载了某个subsystem的hierarchy cgroup根节点所在的目录
*/
func findCgroupMountpoint(subsystem string) string {
	// 打开一个系统目录文件
	// /proc/self/mountinfo 为当前进程的 mountinfo 信息
	// 可以直接通过 cat /proc/self/mountinfo 命令查看
	f, error := os.Open("/proc/self/mountinfo")
	if error != nil {
		return ""
	}
	//最后关闭打开的内容
	defer f.Close()
	//开始逐行读取内容
	scanner := bufio.NewScanner(f)
	//一直读取文件内容，直到返回false
	for scanner.Scan() {
		// txt 大概是这样的：104 85 0:20 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,memory
		txt := scanner.Text()
		// 按照空格区分 获取属性
		fields := strings.Split(txt, " ")
		// 其中的的 memory 就表示这是一个 memory subsystem
		subsystems := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range subsystems {
			if opt == subsystem {
				// 如果等于指定的 subsystem，那么就返回这个挂载点跟目录，就是第四个元素，
				// 这里就是`/sys/fs/cgroup/memory`,即我们要找内存的cgroups的根目录
				return fields[mountPointIndex]
			}
		}
	}

	//如果发现错误这返回位空的挂载点
	error = scanner.Err()

	if error != nil {
		log.Error("read err:", error)
		return ""
	}
	return ""
}
