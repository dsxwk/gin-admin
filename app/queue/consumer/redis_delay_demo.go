package consumer

import (
	"context"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	"fmt"
)

// RedisDelayDemoPayload 延迟消息体
type RedisDelayDemoPayload struct {
	Name string `json:"name"`
}

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
		handleRedisDelayDemo,
	)
}

// handleRedisDelayDemo 处理Redis延迟消息
func handleRedisDelayDemo(_ context.Context, data *RedisDelayDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("Redis Delay Received Msg: name=%s", data.Name))
	return nil
}
