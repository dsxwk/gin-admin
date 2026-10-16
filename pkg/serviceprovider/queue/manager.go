package queue

import (
	"context"
	"errors"
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
func NewManager(config func() *config.Config, consumerFactories []ConsumerFactory, producerFactories []ProducerFactory) (*Manager, error) {
	manager := &Manager{
		config:    config,
		consumers: NewRegistry[Consumer](),
		producers: NewRegistry[Producer](),
	}

	cfg := manager.Config()
	var consumerErrs []error
	for _, factory := range consumerFactories {
		item := factory.Create(cfg)
		if any(item) == nil {
			continue
		}
		if err := manager.consumers.Register(item); err != nil {
			consumerErrs = append(consumerErrs, err)
		}
	}
	if err := errors.Join(consumerErrs...); err != nil {
		return nil, fmt.Errorf("消费者注册失败: %w", err)
	}

	var producerErrs []error
	for _, factory := range producerFactories {
		item := factory.Create(cfg)
		if any(item) == nil {
			continue
		}
		if err := manager.producers.Register(item); err != nil {
			producerErrs = append(producerErrs, err)
		}
	}
	if err := errors.Join(producerErrs...); err != nil {
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

// Start 启动全部消费者
func (m *Manager) Start() error {
	if m == nil {
		return fmt.Errorf("queue manager not initialized")
	}

	cfg := m.Config()
	var errs []error
	for _, consumer := range m.Consumers() {
		if consumer == nil || !connectionEnabled(cfg, consumer.Connection()) {
			continue
		}
		flag.Infof("启动消费者: %s", consumer.Name())
		if err := consumer.Start(); err != nil {
			errs = append(errs, fmt.Errorf("启动消费者 %s 失败: %w", consumer.Name(), err))
		}
	}
	return errors.Join(errs...)
}

// Stop 停止消费者并关闭生产者
func (m *Manager) Stop() error {
	if m == nil {
		return nil
	}

	var errs []error
	for _, consumer := range m.Consumers() {
		if consumer == nil {
			continue
		}
		if err := consumer.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("停止消费者 %s 失败: %w", consumer.Name(), err))
		}
	}
	for _, producer := range m.Producers() {
		if producer == nil {
			continue
		}
		if err := producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭生产者 %s 失败: %w", producer.Name(), err))
		}
	}
	if err := errors.Join(errs...); err != nil {
		return err
	}
	flag.Infof("所有队列服务已关闭")
	return nil
}

// RegisterConsumer 注册消费者
func (m *Manager) RegisterConsumer(consumer Consumer) {
	if m == nil || consumer == nil {
		return
	}
	if err := m.consumers.Register(consumer); err != nil {
		flag.Errorf("%v", err)
	}
}

// RegisterProducer 注册生产者
func (m *Manager) RegisterProducer(producer Producer) {
	if m == nil || producer == nil {
		return
	}
	if err := m.producers.Register(producer); err != nil {
		flag.Errorf("%v", err)
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
	return m.producers.All()
}

// Consumers 获取所有消费者
func (m *Manager) Consumers() []Consumer {
	if m == nil {
		return nil
	}
	return m.consumers.All()
}

// ConsumerNames 获取所有消费者名称
func (m *Manager) ConsumerNames() []string {
	if m == nil {
		return nil
	}
	return m.consumers.Names()
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
			Enabled: connectionEnabled(m.Config(), consumer.Connection()),
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
func (n *nilProducer) Publish(_ context.Context, _ any) error {
	return fmt.Errorf("queue producer [%s] not registered", n.name)
}

// Close 关闭生产者
func (n *nilProducer) Close() error { return nil }
