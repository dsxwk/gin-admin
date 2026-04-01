package facade

import (
	"gin/config"
	"gin/pkg"
	"gin/pkg/logger"
	"gin/pkg/queue"
	"sync"
)

// Queue 队列门面
// 生产者和消费者统一入口
var Queue = &queueFacade{
	producers: make(map[string]queue.Producer),
}

type queueFacade struct {
	mu        sync.RWMutex
	producers map[string]queue.Producer
	cfg       *config.Config
	log       *logger.Logger
	initOnce  sync.Once
}

// init 初始化门面
func (q *queueFacade) init() {
	q.initOnce.Do(func() {
		q.cfg = Config.Get()
		q.log = Log.Logger()
	})
}

// Producer 获取指定名称的生产者
// 使用示例: producer := facade.Queue.Producer("kafka_demo")
func (q *queueFacade) Producer(name string) queue.Producer {
	q.init()

	q.mu.RLock()
	producer, ok := q.producers[name]
	q.mu.RUnlock()

	if ok {
		return producer
	}

	registered := queue.GetProducerRegistry().Get(name)
	if registered == nil {
		q.log.Error(pkg.Sprintf("生产者 %s 未注册", name))
		return nil
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	// 双重检查
	if p, ok := q.producers[name]; ok {
		return p
	}
	q.producers[name] = registered
	return registered
}

// GetAllProducers 获取所有生产者
func (q *queueFacade) GetAllProducers() []queue.Producer {
	q.init()
	q.mu.RLock()
	defer q.mu.RUnlock()

	producers := make([]queue.Producer, 0, len(q.producers))
	for _, p := range q.producers {
		producers = append(producers, p)
	}
	return producers
}

// Consumer 获取指定名称的消费者
func (q *queueFacade) Consumer(name string) queue.Consumer {
	q.init()
	return queue.GetConsumerRegistry().Get(name)
}

// GetAllConsumers 获取所有消费者
func (q *queueFacade) GetAllConsumers() []queue.Consumer {
	q.init()
	return queue.GetConsumerRegistry().GetAll()
}

// GetAllConsumerNames 获取所有消费者名称
func (q *queueFacade) GetAllConsumerNames() []string {
	q.init()
	return queue.GetConsumerRegistry().GetNames()
}

// GetRunningConsumers 获取所有运行中的消费者
func (q *queueFacade) GetRunningConsumers() []queue.Consumer {
	q.init()
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
func (q *queueFacade) GetStoppedConsumers() []queue.Consumer {
	q.init()
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
func (q *queueFacade) GetAllConsumerStats() []ConsumerStats {
	q.init()
	consumers := queue.GetConsumerRegistry().GetAll()
	stats := make([]ConsumerStats, 0, len(consumers))
	for _, c := range consumers {
		stats = append(stats, ConsumerStats{
			Name:    c.Name(),
			Status:  c.Status(),
			Enabled: c.Enabled(q.cfg),
		})
	}
	return stats
}

// ProducerStats 生产者统计信息
type ProducerStats struct {
	Name string `json:"name"`
}

// GetAllProducerStats 获取所有生产者统计信息
func (q *queueFacade) GetAllProducerStats() []ProducerStats {
	q.init()
	producers := queue.GetProducerRegistry().GetAll()
	stats := make([]ProducerStats, 0, len(producers))
	for _, p := range producers {
		stats = append(stats, ProducerStats{
			Name: p.Name(),
		})
	}
	return stats
}
