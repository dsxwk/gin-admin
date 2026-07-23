package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"github.com/segmentio/kafka-go"
)

// KafkaDemoProducer Kafka普通生产者
type KafkaDemoProducer struct {
	*base.KafkaProducer
}

// NewKafkaDemoProducer 创建生产者实例
func NewKafkaDemoProducer() *KafkaDemoProducer {
	cfg := facade.Config()
	kfk := base.NewKafka(cfg, facade.Log(), facade.Message())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "kafka_demo",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &KafkaDemoProducer{
		KafkaProducer: &base.KafkaProducer{
			Kafka: kfk,
			Topic: "kafka_demo",
			Key:   "kafka_demo_key",
		},
	}
}

func (p *KafkaDemoProducer) Name() string {
	return "kafka_demo"
}

func (p *KafkaDemoProducer) Description() string {
	return "kafka普通队列生产者"
}
