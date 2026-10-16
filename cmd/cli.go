//go:build cli

package main

import (
	"gin/app/console"
	"gin/bootstrap"
	"gin/pkg/terminal"
	"gin/pkg/cli"
	"gin/pkg/cli/registry"
	"os"
)

func main() {
	if err := bootstrap.SetupCLI(); err != nil {
		terminal.Errorf("初始化CLI失败: %v", err)
		os.Exit(1)
	}

	cli.Execute(append(console.All(), registry.All()...)...)
}
