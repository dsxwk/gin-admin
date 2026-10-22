package migrate

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/database/migrations"
	"gin/pkg/cli/base"
	"gin/pkg/migration"
	"gin/pkg/terminal"
	"strconv"
	"strings"
)

type Rollback struct {
	base.Command
}

// NewRollback 创建迁移回滚命令
func NewRollback() *Rollback {
	return &Rollback{
		Command: base.NewCommand(base.Config{
			Name:        "migrate:rollback",
			Description: "回滚最近迁移批次",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "s", Long: "step", Default: "1"},
					Desc:     "回滚最近N个批次,默认为1",
					Required: false,
				},
			},
		}),
	}
}

func (s *Rollback) Handle(values map[string]string) error {
	steps, err := parseSteps(values["step"])
	if err != nil {
		return err
	}

	runner, err := migration.NewRunner(facade.DB(), migrations.All())
	if err != nil {
		return err
	}

	executed, err := runner.Rollback(context.Background(), steps)
	for _, name := range executed {
		terminal.Successf("Migration %s 回滚成功", name)
	}
	if err != nil {
		return err
	}
	if len(executed) == 0 {
		terminal.Infof("没有可回滚的迁移")
	}
	return nil
}

// parseSteps 解析回滚批次
func parseSteps(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 1, nil
	}
	steps, err := strconv.Atoi(value)
	if err != nil || steps <= 0 {
		return 0, fmt.Errorf("回滚批次必须为大于0的整数")
	}
	return steps, nil
}
