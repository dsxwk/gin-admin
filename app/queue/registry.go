package queue

import (
	"gin/app/queue/consumer"
	"gin/app/queue/producer"

	sq "gin/pkg/serviceprovider/queue"
)

// Consumers 消费者工厂列表
func Consumers() []sq.ConsumerFactory {
	return []sq.ConsumerFactory{
		sq.NewConsumerFactory("kafka", consumer.NewKafkaDemoConsumer),
		sq.NewConsumerFactory("kafka", consumer.NewKafkaDelayDemoConsumer),
		sq.NewConsumerFactory("rabbitmq", consumer.NewRabbitmqDemoConsumer),
		sq.NewConsumerFactory("rabbitmq", consumer.NewRabbitmqDelayDemoConsumer),
		sq.NewConsumerFactory("redis", consumer.NewRedisDemoConsumer),
		sq.NewConsumerFactory("redis", consumer.NewRedisDelayDemoConsumer),
	}
}

// Producers 生产者工厂列表
func Producers() []sq.ProducerFactory {
	return []sq.ProducerFactory{
		sq.NewProducerFactory("kafka", producer.NewKafkaDemoProducer),
		sq.NewProducerFactory("kafka", producer.NewKafkaDelayDemoProducer),
		sq.NewProducerFactory("rabbitmq", producer.NewRabbitmqDemoProducer),
		sq.NewProducerFactory("rabbitmq", producer.NewRabbitmqDelayDemoProducer),
		sq.NewProducerFactory("redis", producer.NewRedisDemoProducer),
		sq.NewProducerFactory("redis", producer.NewRedisDelayDemoProducer),
	}
}
