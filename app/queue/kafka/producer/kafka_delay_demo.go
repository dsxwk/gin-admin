package producer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/queue"
	"github.com/segmentio/kafka-go"
)

// KafkaDelayDemoProducer Kafka延迟生产者
type KafkaDelayDemoProducer struct {
	*base.KafkaProducer
}

// NewKafkaDelayDemoProducer 创建延迟生产者实例
func NewKafkaDelayDemoProducer() *KafkaDelayDemoProducer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	kfk := base.NewKafka(cfg, log, bus)
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        "kafka_delay_demo",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &KafkaDelayDemoProducer{
		KafkaProducer: &base.KafkaProducer{
			Kafka:        kfk,
			Topic:        "kafka_delay_demo",
			Key:          "kafka_delay_demo_key",
			IsDelayQueue: true,
			DelayMs:      20000, // 20秒延迟
		},
	}
}

func (p *KafkaDelayDemoProducer) Name() string {
	return "kafka_delay_demo"
}

func init() {
	queue.GetProducerRegistry().Register(NewKafkaDelayDemoProducer())
}
