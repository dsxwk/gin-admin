package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	"fmt"
)

// NewRabbitmqDelayDemoProducer 创建RabbitMQ延迟生产者
func NewRabbitmqDelayDemoProducer() *queue.ProducerAdapter {
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(fmt.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "rabbitmq_delay_demo",
			Description: "rabbitmq延迟队列生产者",
			Connection:  "rabbitmq",
			IsDelay:     true,
			DelayMs:     10000,
		},
		&queue.RabbitmqProducer{
			Mq:       mq,
			Queue:    "rabbitmq_delay_demo",
			Exchange: "rabbitmq_delay_demo_exchange",
			Routing:  "rabbitmq_delay_demo",
		},
	)
}
