package consumer

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaDemoConsumer Kafka普通消费者
type KafkaDemoConsumer struct {
	*queue.KafkaConsumer
}

// KafkaDemoPayload 消息体
type KafkaDemoPayload struct {
	Name string `json:"name"`
}

func (c *KafkaDemoConsumer) NewPayload() any {
	return &KafkaDemoPayload{}
}

func (c *KafkaDemoConsumer) Connection() string { return "kafka" }

func (c *KafkaDemoConsumer) Retry() int { return 3 }

func (c *KafkaDemoConsumer) IsDelay() bool { return false }

func (c *KafkaDemoConsumer) Handle(payload any) error {
	data := payload.(*KafkaDemoPayload)
	facade.Log().Info(pkg.Sprintf("Kafka Received Msg: name=%s", data.Name))
	return nil
}

func NewKafkaDemoConsumer() *KafkaDemoConsumer {
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
	kfk.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Queue.Kafka.Brokers,
		Topic:          "kafka_demo",
		GroupID:        "kafka_demo_group",
		MinBytes:       1,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MaxWait:        5 * time.Second,
	})

	return &KafkaDemoConsumer{
		KafkaConsumer: &queue.KafkaConsumer{
			Kafka: kfk,
			Topic: "kafka_demo",
			Group: "kafka_demo_group",
		},
	}
}

func (c *KafkaDemoConsumer) Name() string {
	return "kafka_demo"
}

func (c *KafkaDemoConsumer) Description() string {
	return "kafka普通队列消费者"
}

func (c *KafkaDemoConsumer) Start() error {
	c.KafkaConsumer.Start(c)
	flag.Infof("Kafka消费者启动成功: %s", c.Name())
	return nil
}

func (c *KafkaDemoConsumer) Stop() error {
	return c.KafkaConsumer.Stop()
}

func (c *KafkaDemoConsumer) Enabled(cfg *config.Config) bool {
	return cfg.Queue.Kafka.Enabled
}
