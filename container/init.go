package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// RunContainerInitProcess 启动容器的init进程
/*
这里的init函数是在容器内部执行的，也就是说，代码执行到这里后，容器所在的进程其实就已经创建出来了，
这是本容器执行的第一一个进程。
使用mount先去挂载proc文件系统，以便后面通过ps等系统命令去查看当前进程资源的情况。
*/
func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command:%s", command)

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示声明你要这个新的mount namespace独立。
	// 即 mount proc 之前先把所有挂载点的传播类型改为 private，避免本 namespace 中的挂载事件外泄。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	// 如果不先做 private mount，会导致挂载事件外泄，后续再执行 mydocker 命令时 /proc 文件系统异常
	// 可以执行 mount -t proc proc /proc 命令重新挂载来解决
	// ---分割线---
	/*
		这里 Mount 意思如下：
		NOEXEC 在本文件系统 许运行其 程序。
		MS_NOSUID 在本系统中运行程序的时候， 允许 set-user-ID set-group-ID
		MS_NOD 这个参数是自 Linux 2.4 ，所有 mount 的系统都会默认设定的参数。
	*/
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	// 本函数最后的syscall.Exec是最为重要的一句黑魔法，
	// 正是这个系统调用实现了完成初始化动作并将用户进程运行起来的操作。
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
