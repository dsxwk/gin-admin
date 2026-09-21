package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
)

// RedisDemoProducer Redis普通生产者
type RedisDemoProducer struct {
	*queue.RedisProducer
}

// NewRedisDemoProducer 创建生产者实例
func NewRedisDemoProducer() *RedisDemoProducer {
	p := &RedisDemoProducer{
		RedisProducer: &queue.RedisProducer{
			Queue:     "redis_demo",
			GetClient: facade.RedisClient,
		},
	}

	p.RedisProducer.Owner = p
	return p
}

func (p *RedisDemoProducer) Name() string {
	return "redis_demo"
}

func (p *RedisDemoProducer) Connection() string { return "redis" }

func (p *RedisDemoProducer) IsDelay() bool { return false }

func (p *RedisDemoProducer) DelayMs() int64 { return 0 }

func (p *RedisDemoProducer) Description() string {
	return "redis普通队列生产者"
}
