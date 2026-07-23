package job

import (
	"gin/app/facade"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
)

// SendEmail1Job 发送邮件
type SendEmail1Job struct{}

// SendEmail1 发送邮件参数
type SendEmail1 struct {
}

func (j *SendEmail1Job) Name() string        { return "send_email1" }
func (j *SendEmail1Job) Description() string { return "发送邮件" }
func (j *SendEmail1Job) Connection() string  { return "redis" }
func (j *SendEmail1Job) Retry() int          { return 3 }
func (j *SendEmail1Job) Delay() int64        { return 1000 }
func (j *SendEmail1Job) NewPayload() any     { return &SendEmail1{} }

func (j *SendEmail1Job) Handle(payload any) error {
	data := payload.(*SendEmail1)
	facade.Log().Info(pkg.Sprintf("Job [send_email1] 处理: %v", data))
	return nil
}

func init() {
	job.Register(&SendEmail1Job{})
}
