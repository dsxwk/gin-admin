package queue

import "context"

// Consumer 队列消费者接口
type Consumer interface {
	Name() string
	Description() string
	Connection() string
	Retry() int
	IsDelay() bool
	Start() error
	Stop() error
	Status() string
}

const (
	// ConsumerStatusStopped 已停止
	ConsumerStatusStopped string = "stopped"
	// ConsumerStatusStarting 启动中
	ConsumerStatusStarting string = "starting"
	// ConsumerStatusRunning 运行中
	ConsumerStatusRunning string = "running"
	// ConsumerStatusError 异常
	ConsumerStatusError string = "error"
)

// PayloadHandler 消息负载处理接口
type PayloadHandler interface {
	NewPayload() any
	Handle(payload any) error
}

// ContextPayloadHandler 支持上下文的消息处理接口
type ContextPayloadHandler interface {
	PayloadHandler
	HandleContext(ctx context.Context, payload any) error
}

// Producer 队列生产者接口
type Producer interface {
	Name() string
	Description() string
	Connection() string
	IsDelay() bool
	DelayMs() int64
	Publish(ctx context.Context, msg any) error
	Close() error
}

// Named 名称接口
type Named interface {
	Name() string
}

// ConsumerHandler 消费者处理接口
type ConsumerHandler interface {
	Consumer
	PayloadHandler
}

// ConsumerDriver 消费者驱动接口
type ConsumerDriver interface {
	Start(handler ConsumerHandler)
	Stop() error
	Status() string
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	Name        string
	Description string
	Connection  string
	Retry       int
	IsDelay     bool
}

// ProducerDriver 生产者驱动接口
type ProducerDriver interface {
	Publish(ctx context.Context, msg any) error
	Close() error
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	Name        string
	Description string
	Connection  string
	IsDelay     bool
	DelayMs     int64
}
