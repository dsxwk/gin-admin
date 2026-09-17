//go:build cli

package main

import (
	"gin/bootstrap"
	"gin/common/flag"
	_ "gin/common/imports"
	"gin/pkg/cli"
	"os"
)

func main() {
	if err := bootstrap.InitCLI(); err != nil {
		flag.Errorf("初始化CLI失败: %v", err)
		os.Exit(1)
	}

	cli.Execute()
}
