package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
)

// RabbitmqDemoConsumer RabbitMQ普通消费者
type RabbitmqDemoConsumer struct {
	*base.RabbitmqConsumer
}

// NewRabbitmqDemoConsumer 创建消费者实例
func NewRabbitmqDemoConsumer() *RabbitmqDemoConsumer {
	log := facade.Log()
	mq, err := base.NewRabbitMQ(facade.Config(), log, facade.Message())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return &RabbitmqDemoConsumer{
		RabbitmqConsumer: &base.RabbitmqConsumer{
			Mq:           mq,
			Queue:        "rabbitmq_demo",
			Exchange:     "rabbitmq_demo_exchange",
			Routing:      "rabbitmq_demo",
			IsDelayQueue: false,
			Retry:        3,
		},
	}
}

func (c *RabbitmqDemoConsumer) Name() string {
	return "rabbitmq_demo"
}

func (c *RabbitmqDemoConsumer) Description() string {
	return "rabbitmq普通队列消费者"
}

func (c *RabbitmqDemoConsumer) Start() error {
	c.RabbitmqConsumer.Start(c)
	flag.Infof("RabbitMQ消费者启动成功: %s", c.Name())
	return nil
}

func (c *RabbitmqDemoConsumer) Stop() error {
	return c.RabbitmqConsumer.Stop()
}

func (c *RabbitmqDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Queue.Rabbitmq.Enabled
}

func (c *RabbitmqDemoConsumer) Status() queue.ConsumerStatus {
	return c.RabbitmqConsumer.Status()
}

func (c *RabbitmqDemoConsumer) Handle(msg string) error {
	facade.Log().Info(pkg.Sprintf("RabbitMq Received Msg: %s", msg))
	return nil
}
