package db

import (
	"gin/app/facade"
	"gin/database"
	"gin/database/migrations"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"

	"github.com/fatih/color"
)

type Migrate struct {
	base.Command
}

// NewMigrate 创建命令
func NewMigrate() *Migrate {
	return &Migrate{
		Command: base.NewCommand(base.Config{
			Name:        "db:migrate",
			Description: "数据迁移",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "i", Long: "id"},
					Desc:     "执行指定迁移ID, 如: create_user_table_20251212",
					Required: false,
				},
			},
		}),
	}
}

func (s *Migrate) Handle(values map[string]string) error {
	terminal.Infof("开始执行数据迁移...")

	db := facade.DB()
	db.Exec(`
CREATE TABLE IF NOT EXISTS migrations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    migration VARCHAR(191) NOT NULL UNIQUE,
    created_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`)

	id := values["id"]
	for _, m := range migrations.Migrations() {
		if id != "" && m.ID() != id {
			continue
		}

		var count int64
		db.Model(&database.Migrations{}).Where("migration = ?", m.ID()).Count(&count)
		if count > 0 {
			color.Yellow("Migration %s 已执行,跳过", m.ID())
			continue
		}

		if err := m.Migrate(db); err != nil {
			terminal.Errorf("Migration %s 执行失败: %v", m.ID(), err)
			return nil

		}

		db.Create(&database.Migrations{Migration: m.ID()})
		terminal.Successf("Migration %s 执行成功", m.ID())
	}
	return nil
}
