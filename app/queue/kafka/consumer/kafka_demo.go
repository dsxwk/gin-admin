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

// KafkaDemoConsumer Kafka普通消费者
type KafkaDemoConsumer struct {
	*base.KafkaConsumer
}

// NewKafkaDemoConsumer 创建消费者实例
func NewKafkaDemoConsumer() *KafkaDemoConsumer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	kfk := base.NewKafka(cfg, log, bus)
	kfk.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          "kafka_demo",
		GroupID:        "kafka_demo_group",
		MinBytes:       1,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MaxWait:        5 * time.Second,
	})

	return &KafkaDemoConsumer{
		KafkaConsumer: &base.KafkaConsumer{
			Kafka: kfk,
			Topic: "kafka_demo",
			Group: "kafka_demo_group",
			Retry: 3,
		},
	}
}

func (c *KafkaDemoConsumer) Name() string {
	return "kafka_demo"
}

func (c *KafkaDemoConsumer) Description() string {
	return "kakfa普通队列消费者"
}

func (c *KafkaDemoConsumer) Start(cfg *config.Config, log *logger.Logger) error {
	c.KafkaConsumer.Start(c)
	log.Info(pkg.Sprintf("Kafka消费者启动成功: %s", c.Name()))
	return nil
}

func (c *KafkaDemoConsumer) Stop() error {
	return c.KafkaConsumer.Stop()
}

func (c *KafkaDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Kafka.Enabled
}

func (c *KafkaDemoConsumer) Status() queue.ConsumerStatus {
	return c.KafkaConsumer.Status()
}

func (c *KafkaDemoConsumer) Handle(msg string) error {
	facade.Log.Info(pkg.Sprintf("Kafka Received Msg: %s", msg))
	// todo 处理业务逻辑
	return nil
}

func init() {
	cfg := facade.Config.Get()
	if cfg != nil && cfg.Kafka.Enabled {
		queue.GetConsumerRegistry().Register(NewKafkaDemoConsumer())
	}
}
