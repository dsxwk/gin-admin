package job

import (
	"context"
	"fmt"
	"gin/app/facade"
	servicejob "gin/pkg/serviceprovider/job"
)

// SendEmail 邮件参数
type SendEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// SendEmailJob 发送邮件任务
type SendEmailJob struct{}

// NewSendEmailJob 创建发送邮件任务
func NewSendEmailJob() *servicejob.Adapter[SendEmail] {
	return servicejob.NewAdapter(
		servicejob.Config{
			Name:        "send_email",
			Description: "发送邮件任务",
			Connection:  "redis",
			Retry:       3,
			Delay:       3000,
		},
		(&SendEmailJob{}).Handle,
	)
}

// Handle 处理发送邮件任务
func (j *SendEmailJob) Handle(_ context.Context, data *SendEmail) error {
	facade.Log().Info(fmt.Sprintf("Job [send_email] 发送邮件: to=%s subject=%s", data.To, data.Subject))
	return nil
}
