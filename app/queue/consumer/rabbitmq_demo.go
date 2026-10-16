package consumer

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
)

// RabbitmqDemoPayload 示例消息体
type RabbitmqDemoPayload struct {
	Name string `json:"name"`
}

// RabbitmqDemoConsumer RabbitMQ普通消费者
type RabbitmqDemoConsumer struct{}

// NewRabbitmqDemoConsumer 创建RabbitMQ普通消费者
func NewRabbitmqDemoConsumer() *queue.ConsumerAdapter[RabbitmqDemoPayload] {
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(fmt.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "rabbitmq_demo",
			Description: "rabbitmq普通队列消费者",
			Connection:  "rabbitmq",
			Retry:       3,
		},
		&queue.RabbitmqConsumer{
			Mq:       mq,
			Queue:    "rabbitmq_demo",
			Exchange: "rabbitmq_demo_exchange",
			Routing:  "rabbitmq_demo",
		},
		(&RabbitmqDemoConsumer{}).Handle,
	)
}

// Handle 处理RabbitMQ普通消息
func (c *RabbitmqDemoConsumer) Handle(_ context.Context, data *RabbitmqDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("RabbitMq Received Msg: name=%s", data.Name))
	return nil
}
