package producer

import (
	"gin/config"
	"gin/pkg/serviceprovider/queue"
)

// Register 注册所有 RabbitMQ 生产者(在 config 加载后调用)
func Register(cfg *config.Config) {
	if !cfg.Queue.Rabbitmq.Enabled {
		return
	}
	if p := NewRabbitmqDemoProducer(); p != nil {
		queue.GetProducerRegistry().Register(p)
	}
	if p := NewRabbitmqDelayDemoProducer(); p != nil {
		queue.GetProducerRegistry().Register(p)
	}
}
