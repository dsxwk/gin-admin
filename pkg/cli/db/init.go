package db

import (
	"context"
	"gin/app/facade"
	"gin/database/migrations"
	"gin/database/seeders"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/pathx"
	"gin/pkg/terminal"
	"path/filepath"
)

type Init struct {
	base.Command
}

// NewInit 创建初始化数据库命令
func NewInit() *Init {
	return &Init{
		Command: base.NewCommand(base.Config{
			Name:        "db:init",
			Description: "初始化数据库与数据",
		}),
	}
}

func (s *Init) Handle(_ map[string]string) error {
	terminal.Warningf("初始化会替换现有表和数据")
	db := facade.DB()
	file := filepath.Join(pathx.RootPath(), "database", "gin.sql")
	ctx := context.Background()
	if err := migration.ExecuteSQLFile(ctx, db, file); err != nil {
		return err
	}

	runner, err := migration.NewRunner(db, migrations.All())
	if err != nil {
		return err
	}
	marked, err := runner.MarkMigrated(ctx)
	for _, name := range marked {
		terminal.Successf("Migration %s 已标记执行", name)
	}
	if err != nil {
		return err
	}

	results, err := migration.RunSeeders(ctx, db, seeders.All(), "")
	for _, result := range results {
		terminal.Successf("Seeder %s 执行成功", result.Name)
	}
	return err
}
