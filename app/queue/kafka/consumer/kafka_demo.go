package consumer

import (
	"fmt"
	"gin/common/base"
	"gin/config"
)

type KafkaDemoConsumer struct {
	*base.KafkaConsumer
}

func NewKafkaDemoConsumer() *KafkaDemoConsumer {
	c := &KafkaDemoConsumer{
		&base.KafkaConsumer{
			Reader:       base.NewReader(config.NewConfig().Kafka.Brokers, "kafka_demo", "kafka_demo_group"),
			Topic:        "kafka_demo",
			Group:        "kafka_demo_group",
			Retry:        3,
			IsDelayQueue: false,
		},
	}

	c.Start()

	return c
}

// Handle 业务逻辑
func (c *KafkaDemoConsumer) Handle(msg string) error {
	fmt.Println("Kafka Received Msg:", msg)
	return nil
}

// Start 启动消费者
func (c *KafkaDemoConsumer) Start() {
	c.KafkaConsumer.Start(c)
}

func init() {
	if config.NewConfig().Kafka.Enabled {
		NewKafkaDemoConsumer()
	}
}
