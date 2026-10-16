package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
)

// NewRedisDelayDemoProducer 创建Redis延迟生产者
func NewRedisDelayDemoProducer() *queue.ProducerAdapter {
	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "redis_delay_demo",
			Description: "redis延迟队列生产者",
			Connection:  "redis",
			IsDelay:     true,
			DelayMs:     10000,
		},
		&queue.RedisProducer{
			Queue:     "redis_delay_demo",
			GetClient: facade.RedisClient,
		},
	)
}
