package consumer

import (
	"gin/config"
	"gin/pkg/serviceprovider/queue"
)

// Register 注册所有 RabbitMQ 消费者(在 config 加载后调用)
func Register(cfg *config.Config) {
	if !cfg.Queue.Rabbitmq.Enabled {
		return
	}
	if c := NewRabbitmqDemoConsumer(); c != nil {
		queue.GetConsumerRegistry().Register(c)
	}
	if c := NewRabbitmqDelayDemoConsumer(); c != nil {
		queue.GetConsumerRegistry().Register(c)
	}
}
