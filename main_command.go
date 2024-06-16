package main

import (
	"fmt"

	"mydocker/container"
	"mydocker/subsystems"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:      "run",
	ShortName: "",
	Aliases:   []string{},
	Usage:     "Create a container with namespace and cgroups limit mydocker run -it [command]",
	Flags: []cli.Flag{
		// -i：以交互模式运行容器
		// -t：为容器重新分配一个伪输入终端
		// -it中 i 是主体，输入 i 可以执行命令，t 更像是起到一种美化的功能
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable try",
		},
		//限制容器可以使用的内存大小 意兆位单位
		cli.StringFlag{
			Name:  "mem",
			Usage: "memory limit,e.g: -mem 100m ",
		},
		//限制容器的cpu使用率 按照百分比使设置
		cli.StringFlag{
			Name:  "cpu",
			Usage: "cpu quota,e.g: -cpu 100",
		},
		//限制cpu 限制的进程使用率
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit,e.g: -cputset 2,4",
		},
	},

	/* 这里开始是命令执行的内容 */
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		//创建一个参数的便利 数组sring
		var cmdArray []string
		//添加传递的参数
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("it")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("mem"),
			CpuSet:      context.String("cpuset"),
			CpuCfsQuota: context.Int("cpu"),
		}
		Run(tty, cmdArray, resConf)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	/*
		1.获取传递过来的 command 参数
		2.执行容器初始化操作
	*/
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command: %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
