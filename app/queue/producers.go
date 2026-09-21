package queue

import (
	"gin/app/facade"
	appproducer "gin/app/queue/producer"
	servicequeue "gin/pkg/serviceprovider/queue"
)

// Producers 获取生产者工厂列表
func Producers() []func() servicequeue.Producer {
	return []func() servicequeue.Producer{
		ProducerFactory("kafka", appproducer.NewKafkaDemoProducer),
		ProducerFactory("kafka", appproducer.NewKafkaDelayDemoProducer),
		ProducerFactory("rabbitmq", appproducer.NewRabbitmqDemoProducer),
		ProducerFactory("rabbitmq", appproducer.NewRabbitmqDelayDemoProducer),
		ProducerFactory("redis", appproducer.NewRedisDemoProducer),
		ProducerFactory("redis", appproducer.NewRedisDelayDemoProducer),
	}
}

// ProducerFactory 创建带连接开关的生产者工厂
func ProducerFactory[T servicequeue.Producer](connection string, factory func() T) func() servicequeue.Producer {
	if factory == nil {
		return nil
	}

	return func() servicequeue.Producer {
		if !producerEnabled(connection) {
			return nil
		}
		return factory()
	}
}

// producerEnabled 判断生产者连接是否启用
func producerEnabled(connection string) bool {
	cfg := facade.Config()
	if cfg == nil {
		return false
	}

	switch connection {
	case "kafka":
		return cfg.Queue.Kafka.Enabled
	case "rabbitmq":
		return cfg.Queue.Rabbitmq.Enabled
	default:
		return true
	}
}
