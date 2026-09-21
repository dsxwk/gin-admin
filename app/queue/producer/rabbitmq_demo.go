package producer

import (
	"gin/app/facade"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
)

// RabbitmqDemoProducer RabbitMQ普通生产者
type RabbitmqDemoProducer struct {
	*queue.RabbitmqProducer
}

// NewRabbitmqDemoProducer 创建生产者实例
func NewRabbitmqDemoProducer() *RabbitmqDemoProducer {
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	p := &RabbitmqDemoProducer{
		RabbitmqProducer: &queue.RabbitmqProducer{
			Mq:       mq,
			Queue:    "rabbitmq_demo",
			Exchange: "rabbitmq_demo_exchange",
			Routing:  "rabbitmq_demo",
		},
	}

	p.RabbitmqProducer.Owner = p
	return p
}

func (p *RabbitmqDemoProducer) Name() string {
	return "rabbitmq_demo"
}

func (p *RabbitmqDemoProducer) Connection() string { return "rabbitmq" }

func (p *RabbitmqDemoProducer) IsDelay() bool { return false }

func (p *RabbitmqDemoProducer) DelayMs() int64 { return 0 }

func (p *RabbitmqDemoProducer) Description() string {
	return "rabbitmq普通队列生产者"
}
