package job

import (
	"context"
	"fmt"
	"gin/app/facade"
	servicejob "gin/pkg/serviceprovider/job"
)

// SyncUser 同步用户参数
type SyncUser struct {
	UserID int64  `json:"userId"`
	Action string `json:"action"`
}

// SyncUserJob 同步用户任务
type SyncUserJob struct{}

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
		(&SyncUserJob{}).Handle,
	)
}

// Handle 处理同步用户任务
func (j *SyncUserJob) Handle(_ context.Context, data *SyncUser) error {
	facade.Log().Info(fmt.Sprintf("Job [sync_user] 同步用户: user=%d action=%s", data.UserID, data.Action))
	return nil
}
