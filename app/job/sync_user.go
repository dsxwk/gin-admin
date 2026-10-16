package job

import (
	"context"
	"gin/app/facade"
	servicejob "gin/pkg/serviceprovider/job"
	"fmt"
)

// SyncUser 同步用户参数
type SyncUser struct {
	UserID int64  `json:"userId"`
	Action string `json:"action"`
}

// NewSyncUserJob 创建同步用户任务
func NewSyncUserJob() *servicejob.Adapter[SyncUser] {
	return servicejob.NewAdapter(
		servicejob.Config{
			Name:        "sync_user",
			Description: "同步用户任务",
			Connection:  "sync",
			Retry:       1,
			Delay:       0,
		},
		handleSyncUser,
	)
}

// handleSyncUser 处理同步用户任务
func handleSyncUser(_ context.Context, data *SyncUser) error {
	facade.Log().Info(fmt.Sprintf("Job [sync_user] 同步用户: user=%d action=%s", data.UserID, data.Action))
	return nil
}
