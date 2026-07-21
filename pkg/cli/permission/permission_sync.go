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
	svc.Set(context.Background())

	if err := svc.SyncAllUserPermissions(); err != nil {
		flag.Errorf("权限同步失败: %v", err)
		return
	}

	flag.Successf("权限同步完成")
}

func init() {
	cli.Register(&PermissionSync{})
}
