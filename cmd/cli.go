//go:build cli

package main

import (
	_ "gin/cmd/import"
	"gin/config"
	"gin/pkg/cli"
)

func main() {
	_ = config.NewConfig()

	cli.Execute()
}
