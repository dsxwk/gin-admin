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

// Queue 队列门面实例(单例)
// 使用示例:
//
//	producer := facade.Queue().Producer("kafka_demo")
//	consumers := facade.Queue().GetRunningConsumers()
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

// Producer 获取指定名称的生产者
// 使用示例: producer := facade.Queue().Producer("kafka_demo")
func (q *QueueFacade) Producer(name string) queue.Producer {
	q.mu.Lock()
	defer q.mu.Unlock()

	if p, ok := q.producers[name]; ok {
		return p
	}

	registered := queue.GetProducerRegistry().Get(name)
	if registered == nil {
		flag.Errorf(fmt.Sprintf("队列生产者 [%s] 未注册, 请检查配置是否启用或生产者是否存在", name))
		return &nilProducer{name: name}
	}
	q.producers[name] = registered
	return registered
}

// nilProducer 未找到生产者时的安全桩
type nilProducer struct {
	name string
}

func (n *nilProducer) Name() string        { return n.name }
func (n *nilProducer) Description() string { return "未注册的生产者" }
func (n *nilProducer) Publish(ctx context.Context, msg []byte) error {
	return fmt.Errorf("队列生产者 [%s] 未注册", n.name)
}
func (n *nilProducer) Close() error { return nil }

// GetAllProducers 获取所有生产者
func (q *QueueFacade) GetAllProducers() []queue.Producer {
	q.mu.RLock()
	defer q.mu.RUnlock()

	producers := make([]queue.Producer, 0, len(q.producers))
	for _, p := range q.producers {
		producers = append(producers, p)
	}
	return producers
}

// Consumer 获取指定名称的消费者
func (q *QueueFacade) Consumer(name string) queue.Consumer {
	return queue.GetConsumerRegistry().Get(name)
}

// GetAllConsumers 获取所有消费者
func (q *QueueFacade) GetAllConsumers() []queue.Consumer {
	return queue.GetConsumerRegistry().GetAll()
}

// GetAllConsumerNames 获取所有消费者名称列表
func (q *QueueFacade) GetAllConsumerNames() []string {
	return queue.GetConsumerRegistry().GetNames()
}

// GetRunningConsumers 获取所有运行中的消费者
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

// GetStoppedConsumers 获取所有已停止的消费者
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

// ConsumerStats 消费者统计信息
type ConsumerStats struct {
	Name    string               `json:"name"`
	Status  queue.ConsumerStatus `json:"status"`
	Enabled bool                 `json:"enabled"`
}

// GetAllConsumerStats 获取所有消费者统计信息
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

// ProducerStats 生产者统计信息
type ProducerStats struct {
	Name string `json:"name"`
}

// GetAllProducerStats 获取所有生产者统计信息
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
