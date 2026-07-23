package consumer

import (
	"gin/config"
	"gin/pkg/serviceprovider/queue"
)

// Register 注册所有 Kafka 消费者(在 config 加载后调用)
func Register(cfg *config.Config) {
	if !cfg.Queue.Kafka.Enabled {
		return
	}
	if c := NewKafkaDemoConsumer(); c != nil {
		queue.GetConsumerRegistry().Register(c)
	}
	if c := NewKafkaDelayDemoConsumer(); c != nil {
		queue.GetConsumerRegistry().Register(c)
	}
}
