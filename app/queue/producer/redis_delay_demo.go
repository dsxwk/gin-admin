package producer

import (
	"gin/common/base"
	"gin/pkg/serviceprovider/queue"
)

// RedisDelayDemoProducer Redis延迟生产者
type RedisDelayDemoProducer struct {
	*base.RedisProducer
}

// NewRedisDelayDemoProducer 创建延迟生产者实例
func NewRedisDelayDemoProducer() *RedisDelayDemoProducer {
	p := &RedisDelayDemoProducer{
		RedisProducer: &base.RedisProducer{
			Queue: "redis_delay_demo",
		},
	}

	p.RedisProducer.Owner = p
	return p
}

func (p *RedisDelayDemoProducer) Name() string {
	return "redis_delay_demo"
}

func (p *RedisDelayDemoProducer) Connection() string { return "redis" }

func (p *RedisDelayDemoProducer) IsDelay() bool { return true }

func (p *RedisDelayDemoProducer) DelayMs() int64 { return 10000 }

func (p *RedisDelayDemoProducer) Description() string {
	return "redis延迟队列生产者"
}

func init() {
	queue.GetProducerRegistry().Register(NewRedisDelayDemoProducer())
}
