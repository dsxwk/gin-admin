package job

import (
	"context"
	"gin/app/facade"
	servicejob "gin/pkg/serviceprovider/job"
	"fmt"
)

// SendEmail 邮件参数
type SendEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

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
		handleSendEmail,
	)
}

// handleSendEmail 处理发送邮件任务
func handleSendEmail(_ context.Context, data *SendEmail) error {
	facade.Log().Info(fmt.Sprintf("Job [send_email] 发送邮件: to=%s subject=%s", data.To, data.Subject))
	return nil
}
