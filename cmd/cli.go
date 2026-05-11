//go:build cli

package main

import (
	"gin/app/facade"
	_ "gin/cmd/import"
	"gin/pkg/cli"
	"gin/pkg/provider/orm"
)

func main() {
	conf := facade.Config()
	orm.SetConfig(conf)

	cli.Execute()
}
