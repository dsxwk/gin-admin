package facade

import (
	"context"
	"fmt"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/queue"
)

// Queue 队列门面
func Queue() *QueueFacade {
	return container.Default().Get[*QueueFacade](serviceprovider.ServiceQueue)
}

// NewQueueFacade 创建队列门面
func NewQueueFacade() *QueueFacade {
	return &QueueFacade{}
}

type QueueFacade struct{}

// Register 注册队列消费者或生产者
func (q *QueueFacade) Register[T queue.Named](item T) {
	switch value := any(item).(type) {
	case queue.Consumer:
		if err := queue.GetConsumerRegistry().Register(value); err != nil {
			flag.Errorf("%v", err)
		}
	case queue.Producer:
		if err := queue.GetProducerRegistry().Register(value); err != nil {
			flag.Errorf("%v", err)
		}
	default:
		flag.Errorf("queue register unsupported type: %T", item)
	}
}

func (q *QueueFacade) Producer(name string) queue.Producer {
	registered, ok := queue.GetProducerRegistry().Get(name)
	if !ok {
		flag.Errorf(fmt.Sprintf("queue producer [%s] not registered", name))
		return &nilProducer{name: name}
	}
	return registered
}

type nilProducer struct {
	name string
}

func (n *nilProducer) Name() string        { return n.name }
func (n *nilProducer) Description() string { return "not registered" }
func (n *nilProducer) Connection() string  { return "unknown" }
func (n *nilProducer) IsDelay() bool       { return false }
func (n *nilProducer) DelayMs() int64      { return 0 }
func (n *nilProducer) Publish(ctx context.Context, msg any) error {
	return fmt.Errorf("queue producer [%s] not registered", n.name)
}
func (n *nilProducer) Close() error { return nil }

func (q *QueueFacade) GetAllProducers() []queue.Producer {
	return queue.GetProducerRegistry().GetAll()
}

func (q *QueueFacade) Consumer(name string) queue.Consumer {
	consumer, _ := queue.GetConsumerRegistry().Get(name)
	return consumer
}

func (q *QueueFacade) GetAllConsumers() []queue.Consumer {
	return queue.GetConsumerRegistry().GetAll()
}

func (q *QueueFacade) GetAllConsumerNames() []string {
	return queue.GetConsumerRegistry().GetNames()
}

func (q *QueueFacade) GetRunningConsumers() []queue.Consumer {
	consumers := queue.GetConsumerRegistry().GetAll()
	running := make([]queue.Consumer, 0)
	for _, c := range consumers {
		if c.Status() == queue.ConsumerStatusRunning {
			running = append(running, c)
		}
	}
	return running
}

func (q *QueueFacade) GetStoppedConsumers() []queue.Consumer {
	consumers := queue.GetConsumerRegistry().GetAll()
	stopped := make([]queue.Consumer, 0)
	for _, c := range consumers {
		if c.Status() == queue.ConsumerStatusStopped {
			stopped = append(stopped, c)
		}
	}
	return stopped
}

type ConsumerStats struct {
	Name    string               `json:"name"`
	Status  queue.ConsumerStatus `json:"status"`
	Enabled bool                 `json:"enabled"`
}

func (q *QueueFacade) GetAllConsumerStats() []ConsumerStats {
	consumers := queue.GetConsumerRegistry().GetAll()
	stats := make([]ConsumerStats, 0, len(consumers))
	for _, c := range consumers {
		stats = append(stats, ConsumerStats{
			Name:    c.Name(),
			Status:  c.Status(),
			Enabled: c.Enabled(Config()),
		})
	}
	return stats
}

type ProducerStats struct {
	Name string `json:"name"`
}

func (q *QueueFacade) GetAllProducerStats() []ProducerStats {
	producers := queue.GetProducerRegistry().GetAll()
	stats := make([]ProducerStats, 0, len(producers))
	for _, p := range producers {
		stats = append(stats, ProducerStats{
			Name: p.Name(),
		})
	}
	return stats
}
