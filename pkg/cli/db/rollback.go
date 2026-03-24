package db

import (
	"gin/common/base"
	"gin/database"
	"gin/database/migrations"
	"gin/pkg/cli"
	"gin/pkg/orm"
	"github.com/fatih/color"
)

type Rollback struct {
	base.BaseCommand
}

func (s *Rollback) Name() string {
	return "db:rollback"
}

func (s *Rollback) Description() string {
	return "数据回滚"
}

func (s *Rollback) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "i",
				Long:  "id",
			},
			"回滚指定迁移ID, 如: create_user_table_20251212",
			false,
		},
	}
}

func (s *Rollback) Execute(args []string) {
	values := s.ParseFlags(s.Name(), args, s.Help())
	color.Green("执行命令: %s %s", s.Name(), s.FormatArgs(values))
	color.Cyan("开始执行数据回滚...")

	db := orm.Connection()
	id := values["id"]
	for _, m := range migrations.AllMigrations() {
		if id != "" && m.ID() != id {
			continue
		}

		var count int64
		db.Model(&database.Migrations{}).Where("migration = ?", m.ID()).Count(&count)
		if count == 0 {
			color.Yellow("Migration %s 未执行,跳过", m.ID())
			continue
		}

		if err := m.Rollback(db); err != nil {
			color.Red("Migration %s 回滚失败: %v", m.ID(), err)
			return
		}

		db.Where("migration = ?", m.ID()).Delete(&database.Migrations{})
		color.Green("Migration %s 回滚成功", m.ID())
	}
}

func init() {
	cli.Register(&Rollback{})
}
