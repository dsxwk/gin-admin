package producer

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"

	"github.com/segmentio/kafka-go"
)

// NewKafkaDemoProducer 创建Kafka普通生产者
func NewKafkaDemoProducer() *queue.ProducerAdapter {
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "kafka_demo",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "kafka_demo",
			Description: "kafka普通队列生产者",
			Connection:  "kafka",
		},
		&queue.KafkaProducer{
			Kafka: kfk,
			Topic: "kafka_demo",
			Key:   "kafka_demo_key",
		},
	)
}
