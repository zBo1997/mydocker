package main

import "github.com/urfave/cli"

// 这是一个提示
const usage = `mydocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

// 接下来开始创建一个clik的命令程序
func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{initCommand, runCommand}
}
