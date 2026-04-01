package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/config"
	"gin/pkg"
	"gin/pkg/logger"
	"gin/pkg/queue"
)

// RabbitmqDemoConsumer RabbitMQ普通消费者
type RabbitmqDemoConsumer struct {
	*base.RabbitmqConsumer
}

// NewRabbitmqDemoConsumer 创建消费者实例
func NewRabbitmqDemoConsumer() *RabbitmqDemoConsumer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	// 创建RabbitMQ连接
	mq, err := base.NewRabbitMQ(cfg, log, bus)
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

// Name 消费者名称
func (c *RabbitmqDemoConsumer) Name() string {
	return "rabbitmq_demo"
}

// Start 启动消费者
func (c *RabbitmqDemoConsumer) Start(cfg *config.Config, log *logger.Logger) error {
	c.RabbitmqConsumer.Start(c)
	log.Info(pkg.Sprintf("RabbitMQ消费者启动成功: %s", c.Name()))
	return nil
}

// Stop 停止消费者
func (c *RabbitmqDemoConsumer) Stop() error {
	return c.RabbitmqConsumer.Stop()
}

// Enabled 是否启用
func (c *RabbitmqDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Rabbitmq.Enabled
}

// Status 消费者状态
func (c *RabbitmqDemoConsumer) Status() queue.ConsumerStatus {
	return c.RabbitmqConsumer.Status()
}

// Handle 处理消息的业务逻辑
func (c *RabbitmqDemoConsumer) Handle(msg string) error {
	facade.Log.Info(pkg.Sprintf("RabbitMq Received Msg: %s", msg))
	// todo 处理业务逻辑
	return nil
}

// init 注册消费者到注册表
func init() {
	queue.GetConsumerRegistry().Register(NewRabbitmqDemoConsumer())
}
