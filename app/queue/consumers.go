package queue

import (
	appconsumer "gin/app/queue/consumer"
	servicequeue "gin/pkg/serviceprovider/queue"
)

// Consumers 获取消费者工厂列表
func Consumers() []servicequeue.ConsumerFactory {
	return []servicequeue.ConsumerFactory{
		servicequeue.NewConsumerFactory("kafka", appconsumer.NewKafkaDemoConsumer),
		servicequeue.NewConsumerFactory("kafka", appconsumer.NewKafkaDelayDemoConsumer),
		servicequeue.NewConsumerFactory("rabbitmq", appconsumer.NewRabbitmqDemoConsumer),
		servicequeue.NewConsumerFactory("rabbitmq", appconsumer.NewRabbitmqDelayDemoConsumer),
		servicequeue.NewConsumerFactory("redis", appconsumer.NewRedisDemoConsumer),
		servicequeue.NewConsumerFactory("redis", appconsumer.NewRedisDelayDemoConsumer),
	}
}
