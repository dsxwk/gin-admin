package db

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/database/migrations"
	"gin/pkg/cli"
	"gin/pkg/orm"
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
	flag.Infof("开始执行数据填充...")

	db := orm.Connection()
	id := values["id"]
	for _, seed := range migrations.AllSeeds() {
		if id != "" && seed.ID() != id {
			continue
		}

		if err := seed.Run(db); err != nil {
			flag.Errorf("Seed %s 执行失败: %v", seed.ID(), err)
			return
		}
		flag.Successf("Seed %s 执行成功", seed.ID())
	}
}

func init() {
	cli.Register(&Seed{})
}
