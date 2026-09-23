package queue

import (
	appproducer "gin/app/queue/producer"
	servicequeue "gin/pkg/serviceprovider/queue"
)

// Producers 获取生产者工厂列表
func Producers() []servicequeue.ProducerFactory {
	return []servicequeue.ProducerFactory{
		servicequeue.NewProducerFactory("kafka", appproducer.NewKafkaDemoProducer),
		servicequeue.NewProducerFactory("kafka", appproducer.NewKafkaDelayDemoProducer),
		servicequeue.NewProducerFactory("rabbitmq", appproducer.NewRabbitmqDemoProducer),
		servicequeue.NewProducerFactory("rabbitmq", appproducer.NewRabbitmqDelayDemoProducer),
		servicequeue.NewProducerFactory("redis", appproducer.NewRedisDemoProducer),
		servicequeue.NewProducerFactory("redis", appproducer.NewRedisDelayDemoProducer),
	}
}
