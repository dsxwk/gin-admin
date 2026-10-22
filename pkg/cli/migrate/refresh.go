package migrate

import (
	"context"
	"gin/app/facade"
	"gin/database/migrations"
	"gin/database/seeders"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"
)

type Refresh struct {
	base.Command
}

// NewRefresh 创建迁移刷新命令
func NewRefresh() *Refresh {
	return &Refresh{
		Command: base.NewCommand(base.Config{
			Name:        "migrate:refresh",
			Description: "回滚全部迁移后重新执行",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Long: "seed", Default: "false"},
					Desc:     "迁移完成后执行全部Seeder",
					Required: false,
				},
			},
		}),
	}
}

func (s *Refresh) Handle(values map[string]string) error {
	ctx := context.Background()
	runner, err := migration.NewRunner(facade.DB(), migrations.All())
	if err != nil {
		return err
	}

	rolledBack, err := runner.RollbackAll(ctx)
	for _, name := range rolledBack {
		terminal.Successf("Migration %s 回滚成功", name)
	}
	if err != nil {
		return err
	}

	executed, err := runner.Migrate(ctx)
	for _, name := range executed {
		terminal.Successf("Migration %s 执行成功", name)
	}
	if err != nil {
		return err
	}

	if s.StringToBool(values["seed"]) {
		results, seedErr := migration.RunSeeders(ctx, facade.DB(), seeders.All(), "")
		for _, result := range results {
			terminal.Successf("Seeder %s 执行成功", result.Name)
		}
		if seedErr != nil {
			return seedErr
		}
	}
	return nil
}
