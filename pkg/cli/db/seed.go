package db

import (
	"gin/common/base"
	"gin/database/migrations"
	"gin/pkg/cli"
	"gin/pkg/db/connection"
	"github.com/fatih/color"
)

type Seed struct {
	base.BaseCommand
}

func (s *Seed) Name() string {
	return "db:seed"
}

func (s *Seed) Description() string {
	return "数据填充"
}

func (s *Seed) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{Short: "i", Long: "id"},
			"执行指定Seed ID, 如: 20251212_user_seed",
			false,
		},
	}
}

func (s *Seed) Execute(args []string) {
	values := s.ParseFlags(s.Name(), args, s.Help())
	color.Green("执行命令: %s %s", s.Name(), s.FormatArgs(values))
	color.Cyan("开始执行数据填充...")

	db := connection.Db{}.GetDB()
	id := values["id"]
	for _, seed := range migrations.AllSeeds() {
		if id != "" && seed.ID() != id {
			continue
		}

		if err := seed.Run(db); err != nil {
			color.Red("Seed %s 执行失败: %v", seed.ID(), err)
			return
		}
		color.Green("Seed %s 执行成功", seed.ID())
	}
}

func init() {
	cli.Register(&Seed{})
}
