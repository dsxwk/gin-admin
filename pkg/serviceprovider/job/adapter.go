package job

import (
	"context"
	"fmt"
)

// Adapter 泛型任务适配器
type Adapter[T any] struct {
	config Config
	handle func(ctx context.Context, payload *T) error
}

// NewAdapter 创建泛型任务适配器
func NewAdapter[T any](
	config Config,
	handle func(ctx context.Context, payload *T) error,
) *Adapter[T] {
	if config.Description == "" {
		config.Description = config.Name
	}

	return &Adapter[T]{
		config: config,
		handle: handle,
	}
}

// Name 获取任务名称
func (j *Adapter[T]) Name() string {
	return j.config.Name
}

// Description 获取任务描述
func (j *Adapter[T]) Description() string {
	return j.config.Description
}

// Connection 获取任务连接
func (j *Adapter[T]) Connection() string {
	return j.config.Connection
}

// Retry 获取重试次数
func (j *Adapter[T]) Retry() int {
	return j.config.Retry
}

// Delay 获取任务延迟
func (j *Adapter[T]) Delay() int64 {
	return j.config.Delay
}

// NewPayload 创建任务负载
func (j *Adapter[T]) NewPayload() any {
	return new(T)
}

// Handle 处理任务负载
func (j *Adapter[T]) Handle(payload any) error {
	return j.HandleContext(context.Background(), payload)
}

// HandleContext 处理任务负载,支持上下文
func (j *Adapter[T]) HandleContext(ctx context.Context, payload any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	data, ok := payload.(*T)
	if !ok {
		return fmt.Errorf("job [%s] payload type invalid: %T", j.config.Name, payload)
	}
	if j.handle == nil {
		return fmt.Errorf("job [%s] handler not initialized", j.config.Name)
	}
	return j.handle(ctx, data)
}
