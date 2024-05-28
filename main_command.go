package main

import (
	"fmt"

	"mydocker/container"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:      "run",
	ShortName: "",
	Aliases:   []string{},
	Usage:     "Create a container with namespace and cgroups limit mydocker run -it [command]",
	Flags: []cli.Flag{cli.BoolFlag{
		Name:  "it",
		Usage: "enable try",
	}},
	/* 这里开始是命令执行的内容 */
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("it")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init ",
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
