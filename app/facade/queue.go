package facade

import (
	"context"
	"fmt"
	"gin/common/flag"
	"gin/pkg/serviceprovider/queue"
	"sync"
)

var (
	queueOnce   sync.Once
	queueFacade *QueueFacade
)

func Queue() *QueueFacade {
	queueOnce.Do(func() {
		queueFacade = &QueueFacade{
			producers: make(map[string]queue.Producer),
		}
	})
	return queueFacade
}

type QueueFacade struct {
	mu        sync.RWMutex
	producers map[string]queue.Producer
}

func (q *QueueFacade) Producer(name string) queue.Producer {
	q.mu.Lock()
	defer q.mu.Unlock()

	if p, ok := q.producers[name]; ok {
		return p
	}

	registered := queue.GetProducerRegistry().Get(name)
	if registered == nil {
		flag.Errorf(fmt.Sprintf("queue producer [%s] not registered", name))
		return &nilProducer{name: name}
	}
	q.producers[name] = registered
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
	q.mu.RLock()
	defer q.mu.RUnlock()

	producers := make([]queue.Producer, 0, len(q.producers))
	for _, p := range q.producers {
		producers = append(producers, p)
	}
	return producers
}

func (q *QueueFacade) Consumer(name string) queue.Consumer {
	return queue.GetConsumerRegistry().Get(name)
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
