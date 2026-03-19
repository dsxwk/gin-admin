package producer

import (
	"gin/common/base"
)

type RabbitmqDelayDemoProducer struct {
	*base.RabbitmqProducer
}

func NewRabbitMqDelayDemoProducer() *RabbitmqDelayDemoProducer {
	return &RabbitmqDelayDemoProducer{
		&base.RabbitmqProducer{
			Mq:           base.NewRabbitMq(),
			Queue:        "rabbitmq_delay_demo",
			Exchange:     "rabbitmq_delay_demo_exchange",
			Routing:      "rabbitmq_delay_demo",
			IsDelayQueue: true,
			DelayMs:      10000,
			Headers:      nil,
		},
	}
}
