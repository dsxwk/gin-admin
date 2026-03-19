package producer

import (
	"gin/common/base"
	"gin/config"
)

type KafkaDelayDemoProducer struct {
	*base.KafkaProducer
}

func NewKafkaDelayDemoProducer() *KafkaDelayDemoProducer {
	return &KafkaDelayDemoProducer{
		&base.KafkaProducer{
			Writer:       base.NewWriter(config.NewConfig().Kafka.Brokers, "kafka_delay_demo"),
			Topic:        "kafka_delay_demo",
			Key:          "kafka_delay_demo_key",
			IsDelayQueue: true,
			DelayMs:      20000,
		},
	}
}
