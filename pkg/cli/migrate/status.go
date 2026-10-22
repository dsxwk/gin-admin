package migrate

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/database/migrations"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Status struct {
	base.Command
}

// NewStatus 创建迁移状态命令
func NewStatus() *Status {
	return &Status{
		Command: base.NewCommand(base.Config{
			Name:        "migrate:status",
			Description: "查看迁移状态",
		}),
	}
}

func (s *Status) Handle(_ map[string]string) error {
	db := facade.DB().Session(&gorm.Session{Logger: logger.Discard})
	runner, err := migration.NewRunner(db, migrations.All())
	if err != nil {
		return err
	}
	items, err := runner.Status(context.Background())
	if err != nil {
		return err
	}

	writer := cli.NewTable()
	writer.AppendHeader(table.Row{"Ran?", "Migration", "Batch", "CreatedAt"})
	for _, item := range items {
		ran := "N"
		if item.Ran {
			ran = "Y"
		}
		batch := ""
		if item.Batch > 0 {
			batch = strconv.Itoa(item.Batch)
		}
		createdAt := ""
		if item.CreatedAt != nil {
			createdAt = item.CreatedAt.Format("2006-01-02 15:04:05")
		}
		writer.AppendRow(table.Row{ran, item.Name, batch, createdAt})
	}
	fmt.Println(writer.Render())
	return nil
}
