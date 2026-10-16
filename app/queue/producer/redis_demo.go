package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
)

// NewRedisDemoProducer 创建Redis普通生产者
func NewRedisDemoProducer() *queue.ProducerAdapter {
	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "redis_demo",
			Description: "redis普通队列生产者",
			Connection:  "redis",
		},
		&queue.RedisProducer{
			Queue:     "redis_demo",
			GetClient: facade.RedisClient,
		},
	)
}
