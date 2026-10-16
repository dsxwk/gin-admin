package db

import (
	"gin/app/facade"
	"gin/database"
	"gin/database/migrations"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"

	"github.com/fatih/color"
)

type Rollback struct {
	base.Command
}

// NewRollback 创建命令
func NewRollback() *Rollback {
	return &Rollback{
		Command: base.NewCommand(base.Config{
			Name:        "db:rollback",
			Description: "数据回滚",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "i", Long: "id"},
					Desc:     "回滚指定迁移ID, 如: create_user_table_20251212",
					Required: false,
				},
			},
		}),
	}
}

func (s *Rollback) Handle(values map[string]string) error {
	terminal.Infof("开始执行数据回滚...")

	db := facade.DB()
	id := values["id"]
	for _, m := range migrations.Migrations() {
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
			terminal.Errorf("Migration %s 回滚失败: %v", m.ID(), err)
			return nil

		}

		db.Where("migration = ?", m.ID()).Delete(&database.Migrations{})
		terminal.Successf("Migration %s 回滚成功", m.ID())
	}
	return nil
}
