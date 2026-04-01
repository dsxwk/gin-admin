package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/config"
	"gin/pkg"
	"gin/pkg/logger"
	"gin/pkg/queue"
)

// RabbitmqDelayDemoConsumer RabbitMQ延迟消费者
type RabbitmqDelayDemoConsumer struct {
	*base.RabbitmqConsumer
}

// NewRabbitmqDelayDemoConsumer 创建延迟消费者实例
func NewRabbitmqDelayDemoConsumer() *RabbitmqDelayDemoConsumer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	mq, err := base.NewRabbitMQ(cfg, log, bus)
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
			IsDelayQueue: true, // 标记为延迟队列
			Retry:        3,
		},
	}
}

func (c *RabbitmqDelayDemoConsumer) Name() string {
	return "rabbitmq_delay_demo"
}

func (c *RabbitmqDelayDemoConsumer) Start(cfg *config.Config, log *logger.Logger) error {
	c.RabbitmqConsumer.Start(c)
	log.Info(pkg.Sprintf("RabbitMQ延迟消费者启动成功: %s", c.Name()))
	return nil
}

func (c *RabbitmqDelayDemoConsumer) Stop() error {
	return c.RabbitmqConsumer.Stop()
}

func (c *RabbitmqDelayDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Rabbitmq.Enabled
}

func (c *RabbitmqDelayDemoConsumer) Status() queue.ConsumerStatus {
	return c.RabbitmqConsumer.Status()
}

func (c *RabbitmqDelayDemoConsumer) Handle(msg string) error {
	facade.Log.Info(pkg.Sprintf("RabbitMq Delay Received Msg: %s", msg))
	// todo 处理业务逻辑
	return nil
}

func init() {
	queue.GetConsumerRegistry().Register(NewRabbitmqDelayDemoConsumer())
}
