package permission

import (
	"context"
	"gin/app/service"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
)

type PermissionSync struct{}

func (s *PermissionSync) Name() string {
	return "permission:sync"
}

func (s *PermissionSync) Description() string {
	return "全量同步所有用户权限到Redis"
}

func (s *PermissionSync) Help() []base.CommandOption {
	return nil
}

func (s *PermissionSync) Execute(values map[string]string) {
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
		flag.Errorf("权限同步失败: %v", err)
		return
	}

	cli.StopProgress(tracker, progressDone, true)
	flag.Successf("权限同步完成")
}

func init() {
	cli.Register(&PermissionSync{})
}
