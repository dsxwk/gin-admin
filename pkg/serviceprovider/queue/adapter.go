package queue

import (
	"context"
	"fmt"
)

// ConsumerAdapter 泛型消费者适配器
type ConsumerAdapter[T any] struct {
	config ConsumerConfig
	driver ConsumerDriver
	handle func(ctx context.Context, payload *T) error
}

// NewConsumerAdapter 创建泛型消费者适配器
func NewConsumerAdapter[T any](config ConsumerConfig, driver ConsumerDriver, handle func(ctx context.Context, payload *T) error) *ConsumerAdapter[T] {
	if config.Description == "" {
		config.Description = config.Name
	}
	if config.Retry < 1 {
		config.Retry = 1
	}

	return &ConsumerAdapter[T]{
		config: config,
		driver: driver,
		handle: handle,
	}
}

// Name 消费者名称
func (c *ConsumerAdapter[T]) Name() string {
	return c.config.Name
}

// Description 消费者描述
func (c *ConsumerAdapter[T]) Description() string {
	return c.config.Description
}

// Connection 消费者连接
func (c *ConsumerAdapter[T]) Connection() string {
	return c.config.Connection
}

// Retry 重试次数
func (c *ConsumerAdapter[T]) Retry() int {
	return c.config.Retry
}

// IsDelay 是否延迟
func (c *ConsumerAdapter[T]) IsDelay() bool {
	return c.config.IsDelay
}

// Start 启动消费者
func (c *ConsumerAdapter[T]) Start() error {
	if c.driver == nil {
		return fmt.Errorf("queue consumer [%s] driver not initialized", c.config.Name)
	}
	c.driver.Start(c)
	return nil
}

// Stop 停止消费者
func (c *ConsumerAdapter[T]) Stop() error {
	if c.driver == nil {
		return nil
	}
	return c.driver.Stop()
}

// Status 消费者状态
func (c *ConsumerAdapter[T]) Status() string {
	if c.driver == nil {
		return ConsumerStatusStopped
	}
	return c.driver.Status()
}

// NewPayload 创建消息负载
func (c *ConsumerAdapter[T]) NewPayload() any {
	return new(T)
}

// Handle 处理消息负载
func (c *ConsumerAdapter[T]) Handle(payload any) error {
	return c.HandleContext(context.Background(), payload)
}

// HandleContext 处理消息负载,支持上下文
func (c *ConsumerAdapter[T]) HandleContext(ctx context.Context, payload any) error {
	data, ok := payload.(*T)
	if !ok {
		return fmt.Errorf("queue consumer [%s] payload type invalid: %T", c.config.Name, payload)
	}
	if c.handle == nil {
		return fmt.Errorf("queue consumer [%s] handler not initialized", c.config.Name)
	}
	return c.handle(ctx, data)
}

// producerOwnerSetter 生产者归属设置接口
type producerOwnerSetter interface {
	SetOwner(owner Producer)
}

// ProducerAdapter 生产者适配器
type ProducerAdapter struct {
	config ProducerConfig
	driver ProducerDriver
}

// NewProducerAdapter 创建生产者适配器
func NewProducerAdapter(config ProducerConfig, driver ProducerDriver) *ProducerAdapter {
	if config.Description == "" {
		config.Description = config.Name
	}

	adapter := &ProducerAdapter{
		config: config,
		driver: driver,
	}
	if setter, ok := driver.(producerOwnerSetter); ok {
		setter.SetOwner(adapter)
	}
	return adapter
}

// Name 生产者名称
func (p *ProducerAdapter) Name() string {
	return p.config.Name
}

// Description 生产者描述
func (p *ProducerAdapter) Description() string {
	return p.config.Description
}

// Connection 生产者连接
func (p *ProducerAdapter) Connection() string {
	return p.config.Connection
}

// IsDelay 是否延迟
func (p *ProducerAdapter) IsDelay() bool {
	return p.config.IsDelay
}

// DelayMs 延迟毫秒
func (p *ProducerAdapter) DelayMs() int64 {
	return p.config.DelayMs
}

// Publish 发布消息
func (p *ProducerAdapter) Publish(ctx context.Context, msg any) error {
	if p.driver == nil {
		return fmt.Errorf("queue producer [%s] driver not initialized", p.config.Name)
	}
	return p.driver.Publish(ctx, msg)
}

// Close 关闭生产者
func (p *ProducerAdapter) Close() error {
	if p.driver == nil {
		return nil
	}
	return p.driver.Close()
}
