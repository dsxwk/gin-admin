package consumer

import (
	"context"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	"time"

	"github.com/segmentio/kafka-go"
	"fmt"
)

// KafkaDelayDemoPayload 延迟消息体
type KafkaDelayDemoPayload struct {
	Name string `json:"name"`
}

// NewKafkaDelayDemoConsumer 创建Kafka延迟消费者
func NewKafkaDelayDemoConsumer() *queue.ConsumerAdapter[KafkaDelayDemoPayload] {
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
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

	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "kafka_delay_demo",
			Description: "kafka延迟队列消费者",
			Connection:  "kafka",
			Retry:       3,
			IsDelay:     true,
		},
		&queue.KafkaConsumer{
			Kafka: kfk,
			Topic: "kafka_delay_demo",
			Group: "kafka_delay_demo_group",
		},
		handleKafkaDelayDemo,
	)
}

// handleKafkaDelayDemo 处理Kafka延迟消息
func handleKafkaDelayDemo(_ context.Context, data *KafkaDelayDemoPayload) error {
	facade.Log().Info(fmt.Sprintf("Kafka Delay Received Msg: name=%s", data.Name))
	return nil
}
