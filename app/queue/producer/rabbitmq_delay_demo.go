package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
)

// RabbitmqDelayDemoProducer RabbitMQ延迟生产者
type RabbitmqDelayDemoProducer struct {
	*base.RabbitmqProducer
}

// NewRabbitmqDelayDemoProducer 创建延迟生产者实例
func NewRabbitmqDelayDemoProducer() *RabbitmqDelayDemoProducer {
	log := facade.Log()
	mq, err := base.NewRabbitMQ(facade.Config(), log, facade.Message())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	p := &RabbitmqDelayDemoProducer{
		RabbitmqProducer: &base.RabbitmqProducer{
			Mq:       mq,
			Queue:    "rabbitmq_delay_demo",
			Exchange: "rabbitmq_delay_demo_exchange",
			Routing:  "rabbitmq_delay_demo",
		},
	}

	p.RabbitmqProducer.Owner = p
	return p
}

func (p *RabbitmqDelayDemoProducer) Name() string {
	return "rabbitmq_delay_demo"
}

func (p *RabbitmqDelayDemoProducer) Connection() string { return "rabbitmq" }

func (p *RabbitmqDelayDemoProducer) IsDelay() bool { return true }

func (p *RabbitmqDelayDemoProducer) DelayMs() int64 { return 10000 }

func (p *RabbitmqDelayDemoProducer) Description() string {
	return "rabbitmq延迟队列生产者"
}

func init() {
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Rabbitmq.Enabled {
		if p := NewRabbitmqDelayDemoProducer(); p != nil {
			queue.GetProducerRegistry().Register(p)
		}
	}
}
