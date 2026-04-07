package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg"
	"gin/pkg/queue"
)

// RabbitmqDemoProducer RabbitMQ普通生产者
type RabbitmqDemoProducer struct {
	*base.RabbitmqProducer
}

// NewRabbitmqDemoProducer 创建生产者实例
func NewRabbitmqDemoProducer() *RabbitmqDemoProducer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	mq, err := base.NewRabbitMQ(cfg, log, bus)
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return &RabbitmqDemoProducer{
		RabbitmqProducer: &base.RabbitmqProducer{
			Mq:           mq,
			Queue:        "rabbitmq_demo",
			Exchange:     "rabbitmq_demo_exchange",
			Routing:      "rabbitmq_demo",
			IsDelayQueue: false,
		},
	}
}

func (p *RabbitmqDemoProducer) Name() string {
	return "rabbitmq_demo"
}

func (p *RabbitmqDemoProducer) Description() string {
	return "rabbitmq普通队列生产者"
}

func init() {
	cfg := facade.Config.Get()
	if cfg != nil && cfg.Rabbitmq.Enabled {
		queue.GetProducerRegistry().Register(NewRabbitmqDemoProducer())
	}
}
