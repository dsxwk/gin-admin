package producer

import (
	"gin/common/base"
)

type RabbitmqDemoProducer struct {
	*base.RabbitmqProducer
}

func NewRabbitMqDemoProducer() *RabbitmqDemoProducer {
	return &RabbitmqDemoProducer{
		&base.RabbitmqProducer{
			Mq:           base.NewRabbitMq(),
			Queue:        "rabbitmq_demo",
			Exchange:     "rabbitmq_demo_exchange",
			Routing:      "rabbitmq_demo",
			IsDelayQueue: false,
			DelayMs:      0,
			Headers:      nil,
		},
	}
}
