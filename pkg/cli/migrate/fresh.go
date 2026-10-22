package migrate

import (
	"bufio"
	"context"
	"fmt"
	"gin/app/facade"
	"gin/database/migrations"
	"gin/database/seeders"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"
	"os"
	"strings"

	"gorm.io/gorm"
)

type Fresh struct {
	base.Command
}

// NewFresh 创建数据库重建命令
func NewFresh() *Fresh {
	return &Fresh{
		Command: base.NewCommand(base.Config{
			Name:        "migrate:fresh",
			Description: "删除全部表后重新迁移",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Long: "seed", Default: "false"},
					Desc:     "迁移完成后执行全部Seeder",
					Required: false,
				},
				{
					Flag:     base.Flag{Long: "force", Default: "false"},
					Desc:     "跳过数据库名称确认",
					Required: false,
				},
			},
		}),
	}
}

func (s *Fresh) Handle(values map[string]string) error {
	db := facade.DB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	force := s.StringToBool(values["force"])
	databaseName := db.Migrator().CurrentDatabase()
	if !force && !confirmDatabase(databaseName) {
		return fmt.Errorf("数据库确认失败,操作已取消")
	}

	ctx := context.Background()
	runner, err := migration.NewRunner(db, migrations.All())
	if err != nil {
		return err
	}
	if err = FreshDatabase(ctx, db); err != nil {
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
		results, seedErr := migration.RunSeeders(ctx, db, seeders.All(), "")
		for _, result := range results {
			terminal.Successf("Seeder %s 执行成功", result.Name)
		}
		if seedErr != nil {
			return seedErr
		}
	}
	return nil
}

// FreshDatabase 删除全部数据表
func FreshDatabase(ctx context.Context, db *gorm.DB) error {
	tables, err := db.WithContext(ctx).Migrator().GetTables()
	if err != nil {
		return fmt.Errorf("读取数据表失败: %w", err)
	}
	if len(tables) == 0 {
		return nil
	}

	items := make([]any, 0, len(tables))
	for _, table := range tables {
		if strings.HasPrefix(strings.ToLower(table), "sqlite_") {
			continue
		}
		items = append(items, table)
	}
	if len(items) == 0 {
		return nil
	}

	return db.WithContext(ctx).Connection(func(conn *gorm.DB) error {
		if conn.Dialector.Name() == "mysql" {
			if err := conn.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
				return fmt.Errorf("关闭外键检查失败: %w", err)
			}
			defer func() {
				_ = conn.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
			}()
		}
		if err := conn.Migrator().DropTable(items...); err != nil {
			return fmt.Errorf("删除数据表失败: %w", err)
		}
		return nil
	})
}

// confirmDatabase 确认数据库名称
func confirmDatabase(name string) bool {
	fmt.Printf("输入数据库名%s确认重建: ", name)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input) == name
}
