package job

import (
	"gin/app/facade"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
)

// SendEmailJob 发送邮件任务
type SendEmailJob struct{}

// SendEmail 邮件参数
type SendEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (j *SendEmailJob) Name() string        { return "send_email" }
func (j *SendEmailJob) Description() string { return "发送邮件任务" }
func (j *SendEmailJob) Connection() string  { return "redis" }
func (j *SendEmailJob) Retry() int          { return 3 }
func (j *SendEmailJob) Delay() int64        { return 3000 }
func (j *SendEmailJob) NewPayload() any     { return &SendEmail{} }

func (j *SendEmailJob) Handle(payload any) error {
	data := payload.(*SendEmail)
	facade.Log().Info(pkg.Sprintf("Job [send_email] 发送邮件: to=%s subject=%s", data.To, data.Subject))
	return nil
}

func init() {
	job.Register(&SendEmailJob{})
}
