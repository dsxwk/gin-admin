package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/serviceprovider/queue"
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

	p := &KafkaDemoProducer{
		KafkaProducer: &base.KafkaProducer{
			Kafka: kfk,
			Topic: "kafka_demo",
			Key:   "kafka_demo_key",
		},
	}

	p.KafkaProducer.Owner = p
	return p
}

func (p *KafkaDemoProducer) Name() string {
	return "kafka_demo"
}

func (p *KafkaDemoProducer) Connection() string { return "kafka" }

func (p *KafkaDemoProducer) IsDelay() bool { return false }

func (p *KafkaDemoProducer) DelayMs() int64 { return 0 }

func (p *KafkaDemoProducer) Description() string {
	return "kafka普通队列生产者"
}

func init() {
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Kafka.Enabled {
		if p := NewKafkaDemoProducer(); p != nil {
			queue.GetProducerRegistry().Register(p)
		}
	}
}
