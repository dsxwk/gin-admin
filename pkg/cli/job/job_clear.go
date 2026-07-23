package job

import (
	"context"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/fatih/color"
)

type JobClear struct{}

func (s *JobClear) Name() string {
	return "job:clear"
}

func (s *JobClear) Description() string {
	return "清除所有未消费的Job(仅支持Redis)"
}

func (s *JobClear) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *JobClear) Execute(values map[string]string) {
	if err := facade.Job().Clear(context.Background()); err != nil {
		flag.Errorf("清除Job失败: %v", err)
		return
	}
	color.Green("所有未消费Job已清除")
}

func init() {
	cli.Register(&JobClear{})
}
