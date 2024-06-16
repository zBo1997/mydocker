package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/urfave/cli"
)

/*
一个 Go 中调用 namespace 和 Cgroups 的例子。
注: 运行时需要 root 权限。
*/

func main() {
	// namespace()
	// cgroups()
	//urfaveCli()
	//获取Cgroup的
	testmemorylimit := findCgroupMountpoint("memory")
	fmt.Println(testmemorylimit)
}

// 这里是为了测试cgroup v1版本 看是否可以拿到队友的挂载点
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
		fmt.Println(fields)
		// 其中的的 memory 就表示这是一个 memory subsystem
		subsystems := strings.Split(fields[len(fields)-1], ",")
		for _, opt := range subsystems {
			if opt == subsystem {
				// 如果等于指定的 subsystem，那么就返回这个挂载点跟目录，就是第四个元素，
				// 这里就是`/sys/fs/cgroup/memory`,即我们要找内存的cgroups的根目录
				return fields[4]
			}
		}
	}

	//如果发现错误这返回位空的挂载点
	error = scanner.Err()
	return "empty"
}

// namespace 如何在Go中新建Namespace
func namespace() {
	cmd := exec.Command("bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

// cgroupMemoryHierarchyMount 挂载了memory subsystem的hierarchy的根目录位置
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

// cgroups cgroups初体验
func cgroups() {
	// /proc/self/exe是一个符号链接，代表当前程序的绝对路径
	if os.Args[0] == "/proc/self/exe" {
		// 第一个参数就是当前执行的文件名，所以只有fork出的容器进程才会进入该分支
		fmt.Printf("容器进程内部 PID %d\n", syscall.Getpid())
		// 需要先在宿主机上安装 stress 比如 apt-get install stress
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// 主进程会走这个分支
		cmd := exec.Command("/proc/self/exe")
		cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// 得到 fork 出来的进程在外部namespace 的 pid
		fmt.Println("fork 进程 PID：", cmd.Process.Pid)
		// 在默认的 memory cgroup 下创建子目录，即创建一个子 cgroup
		err := os.Mkdir(filepath.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
		if err != nil {
			fmt.Println(err)
		}
		// 	将容器加入到这个 cgroup 中，即将进程PID加入到cgroup下的 cgroup.procs 文件中
		err = ioutil.WriteFile(filepath.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "cgroup.procs"),
			[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// 	限制进程的内存使用，往 memory.limit_in_bytes 文件中写入数据
		err = ioutil.WriteFile(filepath.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"),
			[]byte("100m"), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, _ = cmd.Process.Wait()
	}
}

// urfaveCli cli 包简单使用，具体可以参考官方文档
func urfaveCli() {
	app := cli.NewApp()

	// 指定全局参数
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "Language for the greeting",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}
	// 指定支持的命令列表
	app.Commands = []cli.Command{
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				log.Println("run command complete")
				for i, v := range c.Args() {
					log.Printf("args i:%v v:%v\n", i, v)
				}
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			// 每个命令下面还可以指定自己的参数
			Flags: []cli.Flag{cli.Int64Flag{
				Name:  "priority",
				Value: 1,
				Usage: "priority for the task",
			}},
			Usage: "add a task to the list",
			Action: func(c *cli.Context) error {
				log.Println("run command add")
				for i, v := range c.Args() {
					log.Printf("args i:%v v:%v\n", i, v)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
