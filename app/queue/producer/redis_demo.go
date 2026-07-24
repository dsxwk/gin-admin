package producer

import (
	"gin/common/base"
	"gin/pkg/serviceprovider/queue"
)

// RedisDemoProducer Redis普通生产者
type RedisDemoProducer struct {
	*base.RedisProducer
}

// NewRedisDemoProducer 创建生产者实例
func NewRedisDemoProducer() *RedisDemoProducer {
	p := &RedisDemoProducer{
		RedisProducer: &base.RedisProducer{
			Queue: "redis_demo",
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

func init() {
	queue.GetProducerRegistry().Register(NewRedisDemoProducer())
}
