package job

import (
	"context"
	"gin/app/facade"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"

	"github.com/fatih/color"
)

type JobClear struct{ base.Command }

// NewJobClear 创建命令
func NewJobClear() *JobClear {
	return &JobClear{
		Command: base.NewCommand(base.Config{
			Name:        "job:clear",
			Description: "清除所有未消费的Job(仅支持Redis)",
		}),
	}
}

func (s *JobClear) Handle(_ map[string]string) error {
	if err := facade.Job().Clear(context.Background()); err != nil {
		terminal.Errorf("清除Job失败: %v", err)
		return nil

	}
	color.Green("所有未消费Job已清除")
	return nil
}
