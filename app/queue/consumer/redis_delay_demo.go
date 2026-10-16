package consumer

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
)

// RedisDelayDemoPayload 延迟消息体
type RedisDelayDemoPayload struct {
	Name string `json:"name"`
}

// RedisDelayDemoConsumer Redis延迟消费者
type RedisDelayDemoConsumer struct{}

// NewRedisDelayDemoConsumer 创建Redis延迟消费者
func NewRedisDelayDemoConsumer() *queue.ConsumerAdapter[RedisDelayDemoPayload] {
	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "redis_delay_demo",
			Description: "redis延迟队列消费者",
			Connection:  "redis",
			Retry:       3,
			IsDelay:     true,
		},
		&queue.RedisConsumer{
			Queue:     "redis_delay_demo",
			GetClient: facade.RedisClient,
			Log:       facade.Log(),
		},
		(&RedisDelayDemoConsumer{}).Handle,
	)
}

// Handle 处理Redis延迟消息
func (c *RedisDelayDemoConsumer) Handle(_ context.Context, data *RedisDelayDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("Redis Delay Received Msg: name=%s", data.Name))
	return nil
}
