package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/config"
	"gin/pkg"
	"gin/pkg/logger"
	"gin/pkg/queue"
	"github.com/segmentio/kafka-go"
	"time"
)

// KafkaDelayDemoConsumer Kafka延迟消费者
type KafkaDelayDemoConsumer struct {
	*base.KafkaConsumer
}

// NewKafkaDelayDemoConsumer 创建延迟消费者实例
func NewKafkaDelayDemoConsumer() *KafkaDelayDemoConsumer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	kfk := base.NewKafka(cfg, log, bus)
	kfk.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          "kafka_delay_demo",
		GroupID:        "kafka_delay_demo_group",
		MinBytes:       1,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MaxWait:        5 * time.Second,
	})

	return &KafkaDelayDemoConsumer{
		KafkaConsumer: &base.KafkaConsumer{
			Kafka:        kfk,
			Topic:        "kafka_delay_demo",
			Group:        "kafka_delay_demo_group",
			Retry:        3,
			IsDelayQueue: true,
		},
	}
}

func (c *KafkaDelayDemoConsumer) Name() string {
	return "kafka_delay_demo"
}

func (c *KafkaDelayDemoConsumer) Start(cfg *config.Config, log *logger.Logger) error {
	c.KafkaConsumer.Start(c)
	log.Info(pkg.Sprintf("Kafka延迟消费者启动成功: %s", c.Name()))
	return nil
}

func (c *KafkaDelayDemoConsumer) Stop() error {
	return c.KafkaConsumer.Stop()
}

func (c *KafkaDelayDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Kafka.Enabled
}

func (c *KafkaDelayDemoConsumer) Status() queue.ConsumerStatus {
	return c.KafkaConsumer.Status()
}

func (c *KafkaDelayDemoConsumer) Handle(msg string) error {
	facade.Log.Info(pkg.Sprintf("Kafka Delay Received Msg: %s", msg))
	// todo 处理业务逻辑
	return nil
}

func init() {
	queue.GetConsumerRegistry().Register(NewKafkaDelayDemoConsumer())
}
