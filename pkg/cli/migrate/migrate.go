package migrate

import (
	"context"
	"gin/app/facade"
	"gin/database/migrations"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"
)

type Migrate struct {
	base.Command
}

// NewMigrate 创建迁移命令
func NewMigrate() *Migrate {
	return &Migrate{
		Command: base.NewCommand(base.Config{
			Name:        "migrate",
			Description: "执行待处理迁移",
		}),
	}
}

func (s *Migrate) Handle(_ map[string]string) error {
	runner, err := migration.NewRunner(facade.DB(), migrations.All())
	if err != nil {
		return err
	}

	executed, err := runner.Migrate(context.Background())
	for _, name := range executed {
		terminal.Successf("Migration %s 执行成功", name)
	}
	if err != nil {
		return err
	}
	if len(executed) == 0 {
		terminal.Infof("没有待处理的迁移")
	}
	return nil
}
