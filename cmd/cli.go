//go:build cli

package main

import (
	"gin/app/facade"
	_ "gin/cmd/import"
	"gin/pkg/cli"
)

func main() {
	_ = facade.Config()

	cli.Execute()
}
