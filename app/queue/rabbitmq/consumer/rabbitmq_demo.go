package consumer

import (
	"fmt"
	"gin/common/base"
	"gin/config"
)

type RabbitmqDemoConsumer struct {
	*base.RabbitmqConsumer
}

func NewRabbitMqDemoConsumer() *RabbitmqDemoConsumer {
	c := &RabbitmqDemoConsumer{
		&base.RabbitmqConsumer{
			Mq:           base.NewRabbitMq(),
			Queue:        "rabbitmq_demo",
			Exchange:     "rabbitmq_demo_exchange",
			Routing:      "rabbitmq_demo",
			Retry:        3,
			IsDelayQueue: false,
		},
	}

	c.Start()

	return c
}

// Start 启动消费者
func (c *RabbitmqDemoConsumer) Start() {
	c.RabbitmqConsumer.Start(c)
}

func (c *RabbitmqDemoConsumer) Handle(msg string) error {
	fmt.Println("RabbitMq Received Msg:", msg)
	return nil
}

func init() {
	if config.NewConfig().Rabbitmq.Enabled {
		NewRabbitMqDemoConsumer()
	}
}
