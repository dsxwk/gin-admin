package db

import (
	"gin/common/base"
	"gin/database"
	"gin/database/migrations"
	"gin/pkg/cli"
	"gin/pkg/db/connection"
	"github.com/fatih/color"
)

type Migrate struct {
	base.BaseCommand
}

func (s *Migrate) Name() string {
	return "db:migrate"
}

func (s *Migrate) Description() string {
	return "数据迁移"
}

func (s *Migrate) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "i",
				Long:  "id",
			},
			"执行指定迁移ID, 如: create_user_table_20251212",
			false,
		},
	}
}

func (s *Migrate) Execute(args []string) {
	values := s.ParseFlags(s.Name(), args, s.Help())
	color.Green("执行命令: %s %s", s.Name(), s.FormatArgs(values))
	color.Cyan("开始执行数据迁移...")

	db := connection.Db{}.GetDB()
	db.Exec(`
CREATE TABLE IF NOT EXISTS migrations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    migration VARCHAR(191) NOT NULL UNIQUE,
    created_at DATETIME(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`)

	id := values["id"]
	for _, m := range migrations.AllMigrations() {
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
			color.Red("Migration %s 执行失败: %v", m.ID(), err)
			return
		}

		db.Create(&database.Migrations{Migration: m.ID()})
		color.Green("Migration %s 执行成功", m.ID())
	}
}

func init() {
	cli.Register(&Migrate{})
}
