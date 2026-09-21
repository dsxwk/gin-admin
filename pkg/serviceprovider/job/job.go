package job

import "gin/pkg/serviceprovider/queue"

// Job 任务接口
type Job interface {
	queue.Named
	queue.PayloadHandler
	// Description 任务描述
	Description() string
	// Connection 任务连接, 默认 "redis",可选 "sync","kafka","rabbitmq"
	Connection() string
	// Retry 重试次数, 默认3
	Retry() int
	// Delay 重试间隔时间(毫秒), 默认1000
	Delay() int64
}
