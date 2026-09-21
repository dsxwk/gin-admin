package queue

import (
	"gin/app/facade"
	appconsumer "gin/app/queue/consumer"
	"gin/config"
	servicequeue "gin/pkg/serviceprovider/queue"
)

// Consumers 获取消费者工厂列表
func Consumers() []func() servicequeue.Consumer {
	return []func() servicequeue.Consumer{
		ConsumerFactory("kafka", appconsumer.NewKafkaDemoConsumer),
		ConsumerFactory("kafka", appconsumer.NewKafkaDelayDemoConsumer),
		ConsumerFactory("rabbitmq", appconsumer.NewRabbitmqDemoConsumer),
		ConsumerFactory("rabbitmq", appconsumer.NewRabbitmqDelayDemoConsumer),
		ConsumerFactory("redis", appconsumer.NewRedisDemoConsumer),
		ConsumerFactory("redis", appconsumer.NewRedisDelayDemoConsumer),
	}
}

// ConsumerFactory 创建带连接开关的消费者工厂
func ConsumerFactory[T servicequeue.Consumer](connection string, factory func() T) func() servicequeue.Consumer {
	if factory == nil {
		return nil
	}

	return func() servicequeue.Consumer {
		if !consumerEnabled(facade.Config(), connection) {
			return nil
		}
		return factory()
	}
}

// consumerEnabled 判断消费者连接是否启用
func consumerEnabled(cfg *config.Config, connection string) bool {
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
