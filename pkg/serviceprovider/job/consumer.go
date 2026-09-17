package job

import (
	"context"
	"fmt"
	"gin/config"
	"gin/pkg/serviceprovider/queue"
)

// Consumer 任务消费者适配器
type Consumer struct {
	name        string
	description string
	connection  string
	isDelay     bool
	redis       *queue.RedisConsumer
	kafka       *queue.KafkaConsumer
	rabbitmq    *queue.RabbitmqConsumer
}

// NewRedisConsumer 创建Redis任务消费者
func NewRedisConsumer(driver *queue.RedisConsumer) *Consumer {
	return &Consumer{
		name:        "job:redis",
		description: "Redis任务消费者",
		connection:  "redis",
		isDelay:     true,
		redis:       driver,
	}
}

// NewKafkaConsumer 创建Kafka任务消费者
func NewKafkaConsumer(driver *queue.KafkaConsumer) *Consumer {
	return &Consumer{
		name:        "job:kafka",
		description: "Kafka任务消费者",
		connection:  "kafka",
		kafka:       driver,
	}
}

// NewRabbitmqConsumer 创建RabbitMQ任务消费者
func NewRabbitmqConsumer(driver *queue.RabbitmqConsumer) *Consumer {
	return &Consumer{
		name:        "job:rabbitmq",
		description: "RabbitMQ任务消费者",
		connection:  "rabbitmq",
		rabbitmq:    driver,
	}
}

// Name 消费者名称
func (c *Consumer) Name() string {
	return c.name
}

// Description 消费者描述
func (c *Consumer) Description() string {
	return c.description
}

// Connection 消费者连接
func (c *Consumer) Connection() string {
	return c.connection
}

// Retry 队列层只处理一次,任务层负责自身重试
func (c *Consumer) Retry() int {
	return 1
}

// IsDelay 是否支持延迟消息
func (c *Consumer) IsDelay() bool {
	return c.isDelay
}

// Start 启动消费者
func (c *Consumer) Start() error {
	switch {
	case c.redis != nil:
		c.redis.Start(c)
	case c.kafka != nil:
		c.kafka.Start(c)
	case c.rabbitmq != nil:
		c.rabbitmq.Start(c)
	default:
		return fmt.Errorf("job consumer driver not initialized")
	}
	return nil
}

// Stop 停止消费者
func (c *Consumer) Stop() error {
	switch {
	case c.redis != nil:
		return c.redis.Stop()
	case c.kafka != nil:
		return c.kafka.Stop()
	case c.rabbitmq != nil:
		return c.rabbitmq.Stop()
	default:
		return nil
	}
}

// Enabled 是否启用消费者
func (c *Consumer) Enabled(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	switch c.connection {
	case "kafka":
		return cfg.Queue.Kafka.Enabled
	case "rabbitmq":
		return cfg.Queue.Rabbitmq.Enabled
	default:
		return true
	}
}

// Status 消费者状态
func (c *Consumer) Status() queue.ConsumerStatus {
	switch {
	case c.redis != nil:
		return c.redis.Status()
	case c.kafka != nil:
		return c.kafka.Status()
	case c.rabbitmq != nil:
		return c.rabbitmq.Status()
	default:
		return queue.ConsumerStatusStopped
	}
}

// NewPayload 创建任务消息
func (c *Consumer) NewPayload() any {
	return &Message{}
}

// Handle 执行任务消息
func (c *Consumer) Handle(payload any) error {
	return c.HandleContext(context.Background(), payload)
}

// HandleContext 执行任务消息,支持上下文取消
func (c *Consumer) HandleContext(ctx context.Context, payload any) error {
	message, ok := payload.(*Message)
	if !ok {
		return fmt.Errorf("job message type invalid: %T", payload)
	}
	return Execute(ctx, *message)
}
