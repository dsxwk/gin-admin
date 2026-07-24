package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
)

// RabbitmqDelayDemoConsumer RabbitMQ延迟消费者
type RabbitmqDelayDemoConsumer struct {
	*base.RabbitmqConsumer
}

// NewRabbitmqDelayDemoConsumer 创建延迟消费者实例
func NewRabbitmqDelayDemoConsumer() *RabbitmqDelayDemoConsumer {
	log := facade.Log()
	mq, err := base.NewRabbitMQ(facade.Config(), log, facade.Message())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return &RabbitmqDelayDemoConsumer{
		RabbitmqConsumer: &base.RabbitmqConsumer{
			Mq:           mq,
			Queue:        "rabbitmq_delay_demo",
			Exchange:     "rabbitmq_delay_demo_exchange",
			Routing:      "rabbitmq_delay_demo",
			IsDelayQueue: true,
			Retry:        3,
		},
	}
}

func (c *RabbitmqDelayDemoConsumer) Name() string {
	return "rabbitmq_delay_demo"
}

func (c *RabbitmqDelayDemoConsumer) Description() string {
	return "rabbitmq延迟队列消费者"
}

func (c *RabbitmqDelayDemoConsumer) Start() error {
	c.RabbitmqConsumer.Start(c)
	flag.Infof("RabbitMQ延迟消费者启动成功: %s", c.Name())
	return nil
}

func (c *RabbitmqDelayDemoConsumer) Stop() error {
	return c.RabbitmqConsumer.Stop()
}

func (c *RabbitmqDelayDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Queue.Rabbitmq.Enabled
}

func (c *RabbitmqDelayDemoConsumer) Status() queue.ConsumerStatus {
	return c.RabbitmqConsumer.Status()
}

func (c *RabbitmqDelayDemoConsumer) Handle(msg string) error {
	facade.Log().Info(pkg.Sprintf("RabbitMq Delay Received Msg: %s", msg))
	return nil
}

func init() {
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Rabbitmq.Enabled {
		if c := NewRabbitmqDelayDemoConsumer(); c != nil {
			queue.GetConsumerRegistry().Register(c)
		}
	}
}
