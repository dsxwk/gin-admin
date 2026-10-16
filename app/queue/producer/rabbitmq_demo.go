package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	"fmt"
)

// NewRabbitmqDemoProducer 创建RabbitMQ普通生产者
func NewRabbitmqDemoProducer() *queue.ProducerAdapter {
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(fmt.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "rabbitmq_demo",
			Description: "rabbitmq普通队列生产者",
			Connection:  "rabbitmq",
		},
		&queue.RabbitmqProducer{
			Mq:       mq,
			Queue:    "rabbitmq_demo",
			Exchange: "rabbitmq_demo_exchange",
			Routing:  "rabbitmq_demo",
		},
	)
}
