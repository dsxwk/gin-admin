package consumer

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
)

// RedisDemoPayload 消息体
type RedisDemoPayload struct {
	Name string `json:"name"`
}

// RedisDemoConsumer Redis普通消费者
type RedisDemoConsumer struct{}

// NewRedisDemoConsumer 创建Redis普通消费者
func NewRedisDemoConsumer() *queue.ConsumerAdapter[RedisDemoPayload] {
	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "redis_demo",
			Description: "redis普通队列消费者",
			Connection:  "redis",
			Retry:       3,
		},
		&queue.RedisConsumer{
			Queue:     "redis_demo",
			GetClient: facade.RedisClient,
			Log:       facade.Log(),
		},
		(&RedisDemoConsumer{}).Handle,
	)
}

// Handle 处理Redis普通消息
func (c *RedisDemoConsumer) Handle(_ context.Context, data *RedisDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("Redis Received Msg: name=%s", data.Name))
	return nil
}
