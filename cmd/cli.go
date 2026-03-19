//go:build cli

package main

import (
	_ "gin/app/command"
	_ "gin/app/listener"
	"gin/config"
	"gin/pkg/cli"
	_ "gin/pkg/cli/db"
	_ "gin/pkg/cli/event"
	_ "gin/pkg/cli/make"
	_ "gin/pkg/cli/route"
)

func main() {
	config.NewConfig()
	cli.Execute()
}
