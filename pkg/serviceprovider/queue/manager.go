package queue

import (
	"context"
	"fmt"
	"gin/common/flag"
	"gin/config"
)

// Manager 队列管理器
type Manager struct {
	config    func() *config.Config
	consumers *Registry[Consumer]
	producers *Registry[Producer]
}

// NewManager 创建队列管理器
func NewManager(
	config func() *config.Config,
	consumerFactories []func() Consumer,
	producerFactories []func() Producer,
) (*Manager, error) {
	manager := &Manager{
		config:    config,
		consumers: NewRegistry[Consumer](),
		producers: NewRegistry[Producer](),
	}

	for _, factory := range consumerFactories {
		manager.consumers.Factory(factory)
	}
	for _, factory := range producerFactories {
		manager.producers.Factory(factory)
	}

	if err := manager.consumers.RegisterFactories(); err != nil {
		return nil, fmt.Errorf("消费者注册失败: %w", err)
	}
	if err := manager.producers.RegisterFactories(); err != nil {
		return nil, fmt.Errorf("生产者注册失败: %w", err)
	}

	return manager, nil
}

// Config 获取配置
func (m *Manager) Config() *config.Config {
	if m == nil || m.config == nil {
		return nil
	}
	return m.config()
}

// Register 注册队列消费者或生产者
func (m *Manager) Register(item any) {
	if m == nil {
		return
	}

	switch value := item.(type) {
	case Consumer:
		if err := m.consumers.Register(value); err != nil {
			flag.Errorf("%v", err)
		}
	case Producer:
		if err := m.producers.Register(value); err != nil {
			flag.Errorf("%v", err)
		}
	default:
		flag.Errorf("queue register unsupported type: %T", item)
	}
}

// Producer 获取生产者
func (m *Manager) Producer(name string) Producer {
	if m == nil {
		flag.Errorf("queue producer [%s] not registered", name)
		return &nilProducer{name: name}
	}

	registered, ok := m.producers.Get(name)
	if !ok {
		flag.Errorf("queue producer [%s] not registered", name)
		return &nilProducer{name: name}
	}
	return registered
}

// Consumer 获取消费者
func (m *Manager) Consumer(name string) Consumer {
	if m == nil {
		return nil
	}

	consumer, _ := m.consumers.Get(name)
	return consumer
}

// Producers 获取所有生产者
func (m *Manager) Producers() []Producer {
	if m == nil {
		return nil
	}
	return m.producers.GetAll()
}

// Consumers 获取所有消费者
func (m *Manager) Consumers() []Consumer {
	if m == nil {
		return nil
	}
	return m.consumers.GetAll()
}

// ConsumerNames 获取所有消费者名称
func (m *Manager) ConsumerNames() []string {
	if m == nil {
		return nil
	}
	return m.consumers.GetNames()
}

// RunningConsumers 获取运行中的消费者
func (m *Manager) RunningConsumers() []Consumer {
	consumers := m.Consumers()
	running := make([]Consumer, 0)
	for _, consumer := range consumers {
		if consumer.Status() == ConsumerStatusRunning {
			running = append(running, consumer)
		}
	}
	return running
}

// StoppedConsumers 获取已停止的消费者
func (m *Manager) StoppedConsumers() []Consumer {
	consumers := m.Consumers()
	stopped := make([]Consumer, 0)
	for _, consumer := range consumers {
		if consumer.Status() == ConsumerStatusStopped {
			stopped = append(stopped, consumer)
		}
	}
	return stopped
}

// ConsumerStatus 消费者状态
type ConsumerStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Enabled bool   `json:"enabled"`
}

// ConsumerStatus 获取消费者状态列表
func (m *Manager) ConsumerStatus() []ConsumerStatus {
	consumers := m.Consumers()
	stats := make([]ConsumerStatus, 0, len(consumers))
	for _, consumer := range consumers {
		stats = append(stats, ConsumerStatus{
			Name:    consumer.Name(),
			Status:  consumer.Status(),
			Enabled: consumer.Enabled(m.Config()),
		})
	}
	return stats
}

// ProducerStatus 生产者统计
type ProducerStatus struct {
	Name string `json:"name"`
}

// ProducerStatus 获取生产者统计列表
func (m *Manager) ProducerStatus() []ProducerStatus {
	producers := m.Producers()
	stats := make([]ProducerStatus, 0, len(producers))
	for _, producer := range producers {
		stats = append(stats, ProducerStatus{Name: producer.Name()})
	}
	return stats
}

// nilProducer 未注册生产者
type nilProducer struct {
	name string
}

// Name 生产者名称
func (n *nilProducer) Name() string { return n.name }

// Description 生产者描述
func (n *nilProducer) Description() string { return "not registered" }

// Connection 生产者连接
func (n *nilProducer) Connection() string { return "unknown" }

// IsDelay 是否延迟
func (n *nilProducer) IsDelay() bool { return false }

// DelayMs 延迟毫秒
func (n *nilProducer) DelayMs() int64 { return 0 }

// Publish 发布消息
func (n *nilProducer) Publish(ctx context.Context, msg any) error {
	return fmt.Errorf("queue producer [%s] not registered", n.name)
}

// Close 关闭生产者
func (n *nilProducer) Close() error { return nil }
