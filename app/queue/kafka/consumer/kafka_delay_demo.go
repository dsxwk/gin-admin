package consumer

import (
	"fmt"
	"gin/common/base"
	"gin/config"
)

type KafkaDelayDemoConsumer struct {
	*base.KafkaConsumer
}

func NewKafkaDelayDemoConsumer() *KafkaDelayDemoConsumer {
	c := &KafkaDelayDemoConsumer{
		&base.KafkaConsumer{
			Reader:       base.NewReader(config.NewConfig().Kafka.Brokers, "kafka_delay_demo", "kafka_delay_demo_group"),
			Topic:        "kafka_delay_demo",
			Group:        "kafka_delay_demo_group",
			Retry:        3,
			IsDelayQueue: true,
		},
	}

	c.Start()

	return c
}

// Handle 处理业务
func (c *KafkaDelayDemoConsumer) Handle(msg string) error {
	fmt.Println("Kafka Received Delay Msg:", msg)
	return nil
}

// Start 启动消费者
func (c *KafkaDelayDemoConsumer) Start() {
	c.KafkaConsumer.Start(c)
}

func init() {
	if config.NewConfig().Kafka.Enabled {
		NewKafkaDelayDemoConsumer()
	}
}
