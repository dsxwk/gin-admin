package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg"
	"gin/pkg/queue"
)

// RabbitmqDelayDemoProducer RabbitMQ延迟生产者
type RabbitmqDelayDemoProducer struct {
	*base.RabbitmqProducer
}

// NewRabbitmqDelayDemoProducer 创建延迟生产者实例
func NewRabbitmqDelayDemoProducer() *RabbitmqDelayDemoProducer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	mq, err := base.NewRabbitMQ(cfg, log, bus)
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return &RabbitmqDelayDemoProducer{
		RabbitmqProducer: &base.RabbitmqProducer{
			Mq:           mq,
			Queue:        "rabbitmq_delay_demo",
			Exchange:     "rabbitmq_delay_demo_exchange",
			Routing:      "rabbitmq_delay_demo",
			IsDelayQueue: true,
			DelayMs:      10000, // 10秒延迟
		},
	}
}

func (p *RabbitmqDelayDemoProducer) Name() string {
	return "rabbitmq_delay_demo"
}

func init() {
	queue.GetProducerRegistry().Register(NewRabbitmqDelayDemoProducer())
}
