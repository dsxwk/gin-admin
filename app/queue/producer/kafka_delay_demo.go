package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"

	"github.com/segmentio/kafka-go"
)

// NewKafkaDelayDemoProducer 创建Kafka延迟生产者
func NewKafkaDelayDemoProducer() *queue.ProducerAdapter {
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "kafka_delay_demo",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "kafka_delay_demo",
			Description: "kafka延迟队列生产者",
			Connection:  "kafka",
			IsDelay:     true,
			DelayMs:     10000,
		},
		&queue.KafkaProducer{
			Kafka: kfk,
			Topic: "kafka_delay_demo",
			Key:   "kafka_delay_demo_key",
		},
	)
}
