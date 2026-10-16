package facade

import (
	"gin/pkg/serviceprovider/queue"
)

// Queue 队列门面
func Queue() *QueueFacade {
	return &QueueFacade{manager: App().Queue()}
}

// QueueFacade 队列门面
type QueueFacade struct {
	manager *queue.Manager
}

// RegisterConsumer 注册队列消费者
func (q *QueueFacade) RegisterConsumer(consumer queue.Consumer) {
	if q == nil {
		return
	}
	q.manager.RegisterConsumer(consumer)
}

// RegisterProducer 注册队列生产者
func (q *QueueFacade) RegisterProducer(producer queue.Producer) {
	if q == nil {
		return
	}
	q.manager.RegisterProducer(producer)
}

func (q *QueueFacade) Producer(name string) queue.Producer {
	if q == nil {
		return nil
	}
	return q.manager.Producer(name)
}

func (q *QueueFacade) Producers() []queue.Producer {
	if q == nil {
		return nil
	}
	return q.manager.Producers()
}

func (q *QueueFacade) Consumer(name string) queue.Consumer {
	if q == nil {
		return nil
	}
	return q.manager.Consumer(name)
}

func (q *QueueFacade) Consumers() []queue.Consumer {
	if q == nil {
		return nil
	}
	return q.manager.Consumers()
}

func (q *QueueFacade) ConsumerNames() []string {
	if q == nil {
		return nil
	}
	return q.manager.ConsumerNames()
}

func (q *QueueFacade) RunningConsumers() []queue.Consumer {
	if q == nil {
		return nil
	}
	return q.manager.RunningConsumers()
}

func (q *QueueFacade) StoppedConsumers() []queue.Consumer {
	if q == nil {
		return nil
	}
	return q.manager.StoppedConsumers()
}

type ConsumerStatus = queue.ConsumerStatus

func (q *QueueFacade) ConsumerStatus() []ConsumerStatus {
	if q == nil {
		return nil
	}
	return q.manager.ConsumerStatus()
}

type ProducerStatus = queue.ProducerStatus

func (q *QueueFacade) ProducerStatus() []ProducerStatus {
	if q == nil {
		return nil
	}
	return q.manager.ProducerStatus()
}
