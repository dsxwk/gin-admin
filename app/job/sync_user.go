package job

import (
	"gin/app/facade"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
)

// SyncUserJob 同步用户任务
type SyncUserJob struct{}

// SyncUser 同步用户参数
type SyncUser struct {
	UserID int64  `json:"user_id"`
	Action string `json:"action"`
}

func (j *SyncUserJob) Name() string        { return "sync_user" }
func (j *SyncUserJob) Description() string { return "同步用户任务" }
func (j *SyncUserJob) Connection() string  { return "sync" }
func (j *SyncUserJob) Retry() int          { return 1 }
func (j *SyncUserJob) Delay() int64        { return 0 }
func (j *SyncUserJob) NewPayload() any     { return &SyncUser{} }

func (j *SyncUserJob) Handle(payload any) error {
	data := payload.(*SyncUser)
	facade.Log().Info(pkg.Sprintf("Job [sync_user] 同步用户: user=%d action=%s", data.UserID, data.Action))
	return nil
}

func init() {
	job.Register(&SyncUserJob{})
}
