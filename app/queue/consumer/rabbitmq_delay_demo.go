package consumer

import (
	"context"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	"fmt"
)

// RabbitmqDelayDemoPayload 延迟消息体
type RabbitmqDelayDemoPayload struct {
	Name string `json:"name"`
}

// NewRabbitmqDelayDemoConsumer 创建RabbitMQ延迟消费者
func NewRabbitmqDelayDemoConsumer() *queue.ConsumerAdapter[RabbitmqDelayDemoPayload] {
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(fmt.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "rabbitmq_delay_demo",
			Description: "rabbitmq延迟队列消费者",
			Connection:  "rabbitmq",
			Retry:       3,
			IsDelay:     true,
		},
		&queue.RabbitmqConsumer{
			Mq:       mq,
			Queue:    "rabbitmq_delay_demo",
			Exchange: "rabbitmq_delay_demo_exchange",
			Routing:  "rabbitmq_delay_demo",
		},
		handleRabbitmqDelayDemo,
	)
}

// handleRabbitmqDelayDemo 处理RabbitMQ延迟消息
func handleRabbitmqDelayDemo(_ context.Context, data *RabbitmqDelayDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("RabbitMq Delay Received Msg: name=%s", data.Name))
	return nil
}
