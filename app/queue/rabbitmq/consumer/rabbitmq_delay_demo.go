package consumer

import (
	"fmt"
	"gin/common/base"
	"gin/config"
)

type RabbitmqDelayDemoConsumer struct {
	*base.RabbitmqConsumer
}

func NewRabbitMqDelayDemoConsumer() *RabbitmqDelayDemoConsumer {
	c := &RabbitmqDelayDemoConsumer{
		&base.RabbitmqConsumer{
			Mq:           base.NewRabbitMq(),
			Queue:        "rabbitmq_delay_demo",
			Exchange:     "rabbitmq_delay_demo_exchange",
			Routing:      "rabbitmq_delay_demo",
			Retry:        3,
			IsDelayQueue: true,
		},
	}

	c.Start()

	return c
}

// Start 启动消费者
func (c *RabbitmqDelayDemoConsumer) Start() {
	c.RabbitmqConsumer.Start(c)
}

func (c *RabbitmqDelayDemoConsumer) Handle(msg string) error {
	fmt.Println("RabbitMq Received Delay Msg:", msg)
	return nil
}

func init() {
	if config.NewConfig().Rabbitmq.Enabled {
		NewRabbitMqDelayDemoConsumer()
	}
}
