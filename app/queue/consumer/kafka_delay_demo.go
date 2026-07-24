package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
	"github.com/segmentio/kafka-go"
	"time"
)

// KafkaDelayDemoConsumer Kafka延迟消费者
type KafkaDelayDemoConsumer struct {
	*base.KafkaConsumer
}

// KafkaDelayDemoPayload 延迟消息体
type KafkaDelayDemoPayload struct {
	Name string `json:"name"`
}

func (c *KafkaDelayDemoConsumer) NewPayload() any {
	return &KafkaDelayDemoPayload{}
}

func (c *KafkaDelayDemoConsumer) Connection() string { return "kafka" }

func (c *KafkaDelayDemoConsumer) Retry() int { return 3 }

func (c *KafkaDelayDemoConsumer) IsDelay() bool { return true }

func (c *KafkaDelayDemoConsumer) Handle(payload any) error {
	data := payload.(*KafkaDelayDemoPayload)
	facade.Log().Info(pkg.Sprintf("Kafka Delay Received Msg: name=%s", data.Name))
	return nil
}

func NewKafkaDelayDemoConsumer() *KafkaDelayDemoConsumer {
	cfg := facade.Config()
	kfk := base.NewKafka(cfg, facade.Log(), facade.Message())
	kfk.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Queue.Kafka.Brokers,
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
			Kafka: kfk,
			Topic: "kafka_delay_demo",
			Group: "kafka_delay_demo_group",
		},
	}
}

func (c *KafkaDelayDemoConsumer) Name() string {
	return "kafka_delay_demo"
}

func (c *KafkaDelayDemoConsumer) Description() string {
	return "kafka延迟队列消费者"
}

func (c *KafkaDelayDemoConsumer) Start() error {
	c.KafkaConsumer.Start(c)
	flag.Infof("Kafka延迟消费者启动成功: %s", c.Name())
	return nil
}

func (c *KafkaDelayDemoConsumer) Stop() error {
	return c.KafkaConsumer.Stop()
}

func (c *KafkaDelayDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Queue.Kafka.Enabled
}

func (c *KafkaDelayDemoConsumer) Status() queue.ConsumerStatus {
	return c.KafkaConsumer.Status()
}

func init() {
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Kafka.Enabled {
		if c := NewKafkaDelayDemoConsumer(); c != nil {
			queue.GetConsumerRegistry().Register(c)
		}
	}
}
