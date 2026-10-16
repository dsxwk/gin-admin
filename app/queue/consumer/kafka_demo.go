package consumer

import (
	"context"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	"time"

	"github.com/segmentio/kafka-go"
	"fmt"
)

// KafkaDemoPayload 消息体
type KafkaDemoPayload struct {
	Name string `json:"name"`
}

// NewKafkaDemoConsumer 创建Kafka普通消费者
func NewKafkaDemoConsumer() *queue.ConsumerAdapter[KafkaDemoPayload] {
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

	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "kafka_demo",
			Description: "kafka普通队列消费者",
			Connection:  "kafka",
			Retry:       3,
		},
		&queue.KafkaConsumer{
			Kafka: kfk,
			Topic: "kafka_demo",
			Group: "kafka_demo_group",
		},
		handleKafkaDemo,
	)
}

// handleKafkaDemo 处理Kafka普通消息
func handleKafkaDemo(_ context.Context, data *KafkaDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("Kafka Received Msg: name=%s", data.Name))
	return nil
}
