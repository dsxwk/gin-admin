package db

import (
	"context"

	"gin/app/facade"
	"gin/database/seeders"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"
)

type Seed struct {
	base.Command
}

// NewSeed 创建数据填充命令
func NewSeed() *Seed {
	return &Seed{
		Command: base.NewCommand(base.Config{
			Name:        "db:seed",
			Description: "执行数据填充",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "n", Long: "name"},
					Desc:     "执行指定Seeder,如:UserSeeder",
					Required: false,
				},
			},
		}),
	}
}

func (s *Seed) Handle(values map[string]string) error {
	results, err := migration.RunSeeders(
		context.Background(),
		facade.DB(),
		seeders.All(),
		values["name"],
	)
	for _, result := range results {
		terminal.Successf("Seeder %s 执行成功", result.Name)
	}
	return err
}
