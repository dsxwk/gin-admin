package queue

import "gin/config"

// ConsumerFactory 消费者工厂
type ConsumerFactory struct {
	connection string
	create     func() Consumer
}

// NewConsumerFactory 创建消费者工厂
func NewConsumerFactory[T Consumer](connection string, create func() T) ConsumerFactory {
	factory := ConsumerFactory{connection: connection}
	if create != nil {
		factory.create = func() Consumer {
			return create()
		}
	}
	return factory
}

// Create 创建消费者
func (f ConsumerFactory) Create(cfg *config.Config) Consumer {
	if f.create == nil || !connectionEnabled(cfg, f.connection) {
		return nil
	}
	return f.create()
}

// ProducerFactory 生产者工厂
type ProducerFactory struct {
	connection string
	create     func() Producer
}

// NewProducerFactory 创建生产者工厂
func NewProducerFactory[T Producer](connection string, create func() T) ProducerFactory {
	factory := ProducerFactory{connection: connection}
	if create != nil {
		factory.create = func() Producer {
			return create()
		}
	}
	return factory
}

// Create 创建生产者
func (f ProducerFactory) Create(cfg *config.Config) Producer {
	if f.create == nil || !connectionEnabled(cfg, f.connection) {
		return nil
	}
	return f.create()
}

// connectionEnabled 判断连接是否启用
func connectionEnabled(cfg *config.Config, connection string) bool {
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
