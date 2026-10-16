package job

import (
	"context"
	"gin/pkg/serviceprovider/queue"
)

// Job 任务接口
type Job interface {
	queue.Named
	queue.PayloadHandler
	// Description 任务描述
	Description() string
	// Connection 任务连接, 默认 "redis",可选 "sync","kafka","rabbitmq"
	Connection() string
	// Retry 重试次数, 0表示不重试
	Retry() int
	// Delay 首次延迟及重试间隔(毫秒), 0表示不延迟
	Delay() int64
}

// ContextJob 支持上下文的任务接口
type ContextJob interface {
	Job
	HandleContext(ctx context.Context, payload any) error
}

// Config 任务配置
type Config struct {
	Name        string
	Description string
	Connection  string
	Retry       int
	Delay       int64
}
