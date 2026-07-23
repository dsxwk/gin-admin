package producer

import (
	"gin/config"
	"gin/pkg/serviceprovider/queue"
)

// Register 注册所有 Kafka 生产者(在 config 加载后调用)
func Register(cfg *config.Config) {
	if !cfg.Queue.Kafka.Enabled {
		return
	}
	if p := NewKafkaDemoProducer(); p != nil {
		queue.GetProducerRegistry().Register(p)
	}
	if p := NewKafkaDelayDemoProducer(); p != nil {
		queue.GetProducerRegistry().Register(p)
	}
}
