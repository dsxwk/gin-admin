package permission

import (
	"context"
	"gin/app/service"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
)

type PermissionSync struct{ base.Command }

// NewPermissionSync 创建命令
func NewPermissionSync() *PermissionSync {
	return &PermissionSync{
		Command: base.NewCommand(base.Config{
			Name:        "permission:sync",
			Description: "全量同步所有用户权限到Redis",
			Options:     nil,
		}),
	}
}

func (s *PermissionSync) Handle(_ map[string]string) error {
	svc := service.RoleService{}
	tracker, progressDone := cli.StartProgress("同步用户权限", 0)

	err := svc.SyncWithProgress(context.Background(), false, func(current, total int) {
		if current == 0 {
			tracker.UpdateTotal(int64(total))
			return
		}
		tracker.Increment(1)
	})
	if err != nil {
		cli.StopProgress(tracker, progressDone, false)
		terminal.Errorf("权限同步失败: %v", err)
		return nil

	}

	cli.StopProgress(tracker, progressDone, true)
	terminal.Successf("权限同步完成")
	return nil
}
