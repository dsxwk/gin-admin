package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/serviceprovider/queue"
	"github.com/segmentio/kafka-go"
)

// KafkaDelayDemoProducer Kafka延迟生产者
type KafkaDelayDemoProducer struct {
	*base.KafkaProducer
}

// NewKafkaDelayDemoProducer 创建延迟生产者实例
func NewKafkaDelayDemoProducer() *KafkaDelayDemoProducer {
	cfg := facade.Config()
	kfk := base.NewKafka(cfg, facade.Log(), facade.Message())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "kafka_delay_demo",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	p := &KafkaDelayDemoProducer{
		KafkaProducer: &base.KafkaProducer{
			Kafka: kfk,
			Topic: "kafka_delay_demo",
			Key:   "kafka_delay_demo_key",
		},
	}

	p.KafkaProducer.Owner = p
	return p
}

func (p *KafkaDelayDemoProducer) Name() string {
	return "kafka_delay_demo"
}

func (p *KafkaDelayDemoProducer) Connection() string { return "kafka" }

func (p *KafkaDelayDemoProducer) IsDelay() bool { return true }

func (p *KafkaDelayDemoProducer) DelayMs() int64 { return 10000 }

func (p *KafkaDelayDemoProducer) Description() string {
	return "kafka延迟队列生产者"
}

func init() {
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Kafka.Enabled {
		if p := NewKafkaDelayDemoProducer(); p != nil {
			queue.GetProducerRegistry().Register(p)
		}
	}
}
