package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// NewParentProcess 构建 command 用于启动一个新进程
/*
这里是父进程，也就是当前进程执行的内容。
1.这里的/proc/se1f/exe调用中，/proc/self/ 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程进行初始化
2.后面的args是参数，其中init是传递给本进程的第一个参数，在本例中，其实就是会去调用initCommand去初始化进程的一些环境和资源
3.下面的clone参数就是去fork出来一个新进程，并且使用了namespace隔离新创建的进程和外部环境。
4.如果用户指定了-it参数，就需要把当前进程的输入输出导入到标准输入输出上
*/
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {

	//为了解决参数过多导致程序异常，所以采用匿名管道的方式获取用户参数
	//虽然匿名管道自带 4K 缓冲，但是如果写满之后就会阻塞，因此最好是等子进程启动后，再往里面写，尽量避免意外情况。
	readPipe, writePipe, err := os.Pipe()

	if err != nil {
		log.Errorf("New Pipe error %v", err)
	}
	//在新的父进程中执行init 命令
	args := []string{"init"}
	/*
		   "/proc/self/exe" 表示调用自身init命令，用来初始化容器;他是Linux中的一个符号连接地址，它指向当前进程可
		 	执行的文件，体而言，/proc/self 是一个指向当前进程自身的符号链接，而 exe 则是一个特殊的文件，通过这个文件
			可以访问当前进程的可执行文件。因此，/proc/self/exe 实际上是当前进程可执行文件的路径。
	*/
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//这里通过CloneFlags 在初始化容器的时候
		//fork 新进程时，通过指定 Cloneflags 会创建对应的 Namespace 以实现隔离，这里包括UTS（主机名）、PID（进程ID）、
		//挂载点、网络、IPC等方面的隔离。
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	/*
		当用户指定 -it 参数时，就将 cmd 的输入和输出连接到终端，以便我们可以与命令进行交互，并看到命令的输出。
	*/

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	//将 readPipe 作为 ExtraFiles，这样 cmd 执行时就会外带着这个文件句柄去创建子进程。
	//ExtraFiles 是 Go 语言 os/exec 包中 Cmd 结构体的一个字段。 使得可以把父进程的命令传递给子进程
	cmd.ExtraFiles = []*os.File{readPipe}
	//并且发写的管道返回出给父进程使用
	return cmd, writePipe
}
