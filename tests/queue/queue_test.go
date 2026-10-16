package queue

import (
	"context"
	"gin/app/facade"
	"gin/config"
	"gin/pkg/serviceprovider/queue"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

type queueRegistryTestItem struct {
	name string
}

// Name 注册项名称
func (i queueRegistryTestItem) Name() string {
	return i.name
}

type queueFacadeTestProducer struct {
	name string
}

// Name 生产者名称
func (p *queueFacadeTestProducer) Name() string {
	return p.name
}

// Description 生产者描述
func (p *queueFacadeTestProducer) Description() string {
	return "queue facade test producer"
}

// Connection 生产者连接类型
func (p *queueFacadeTestProducer) Connection() string {
	return "test"
}

// IsDelay 是否为延迟生产者
func (p *queueFacadeTestProducer) IsDelay() bool {
	return false
}

// DelayMs 延迟毫秒数
func (p *queueFacadeTestProducer) DelayMs() int64 {
	return 0
}

// Publish 发送测试消息
func (p *queueFacadeTestProducer) Publish(_ context.Context, _ any) error {
	return nil
}

// Close 关闭生产者
func (p *queueFacadeTestProducer) Close() error {
	return nil
}

type queueConsumerTestDriver struct {
	started bool
	stopped bool
}

// Start 启动测试消费者驱动
func (d *queueConsumerTestDriver) Start(_ queue.ConsumerHandler) {
	d.started = true
}

// Stop 停止测试消费者驱动
func (d *queueConsumerTestDriver) Stop() error {
	d.stopped = true
	return nil
}

// Status 获取测试消费者驱动状态
func (d *queueConsumerTestDriver) Status() string {
	return queue.ConsumerStatusStopped
}

type queueAdapterTestPayload struct {
	Name string `json:"name"`
}

type queueProducerTestDriver struct {
	owner   queue.Producer
	closed  bool
	message any
}

// SetOwner 设置测试生产者归属
func (d *queueProducerTestDriver) SetOwner(owner queue.Producer) {
	d.owner = owner
}

// Publish 发布测试消息
func (d *queueProducerTestDriver) Publish(_ context.Context, message any) error {
	d.message = message
	return nil
}

// Close 关闭测试生产者
func (d *queueProducerTestDriver) Close() error {
	d.closed = true
	return nil
}

// TestQueueRegistryRegisterAndGet 测试注册表注册和查询
func TestQueueRegistryRegisterAndGet(t *testing.T) {
	registry := queue.NewRegistry[queueRegistryTestItem]()
	item := queueRegistryTestItem{name: "queue_registry_test"}

	require.NoError(t, registry.Register(item))

	current, ok := registry.Get(item.Name())
	require.True(t, ok)
	require.Equal(t, item, current)

	require.Error(t, registry.Register(item))

	missing, ok := registry.Get("queue_registry_missing")
	require.False(t, ok)
	require.Equal(t, queueRegistryTestItem{}, missing)
}

// TestQueueConsumerAdapter 测试泛型消费者适配器
func TestQueueConsumerAdapter(t *testing.T) {
	driver := &queueConsumerTestDriver{}
	handled := false
	adapter := queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "queue_adapter_test",
			Description: "queue adapter test",
			Connection:  "test",
			Retry:       3,
		},
		driver,
		func(_ context.Context, payload *queueAdapterTestPayload) error {
			handled = payload.Name == "adapter"
			return nil
		},
	)

	require.Equal(t, "queue_adapter_test", adapter.Name())
	require.Equal(t, 3, adapter.Retry())
	require.False(t, adapter.IsDelay())
	require.IsType(t, &queueAdapterTestPayload{}, adapter.NewPayload())
	require.NoError(t, adapter.Handle(&queueAdapterTestPayload{Name: "adapter"}))
	require.True(t, handled)
	require.NoError(t, adapter.Start())
	require.True(t, driver.started)
	require.NoError(t, adapter.Stop())
	require.True(t, driver.stopped)
}

// TestQueueProducerAdapter 测试生产者适配器
func TestQueueProducerAdapter(t *testing.T) {
	driver := &queueProducerTestDriver{}
	adapter := queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "queue_producer_adapter_test",
			Description: "queue producer adapter test",
			Connection:  "test",
			IsDelay:     true,
			DelayMs:     1000,
		},
		driver,
	)

	require.Same(t, adapter, driver.owner)
	require.Equal(t, int64(1000), adapter.DelayMs())
	require.NoError(t, adapter.Publish(context.Background(), "message"))
	require.Equal(t, "message", driver.message)
	require.NoError(t, adapter.Close())
	require.True(t, driver.closed)
}

// TestQueueManagerStartStop 测试队列管理器启停
func TestQueueManagerStartStop(t *testing.T) {
	consumerDriver := &queueConsumerTestDriver{}
	producerDriver := &queueProducerTestDriver{}
	consumer := queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:       "queue_manager_consumer_test",
			Connection: "test",
		},
		consumerDriver,
		func(_ context.Context, _ *queueAdapterTestPayload) error {
			return nil
		},
	)
	producer := queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:       "queue_manager_producer_test",
			Connection: "test",
		},
		producerDriver,
	)

	manager, err := queue.NewManager(
		func() *config.Config {
			return &config.Config{}
		},
		[]queue.ConsumerFactory{
			queue.NewConsumerFactory("test", func() queue.Consumer {
				return consumer
			}),
		},
		[]queue.ProducerFactory{
			queue.NewProducerFactory("test", func() queue.Producer {
				return producer
			}),
		},
	)
	require.NoError(t, err)
	require.NoError(t, manager.Start())
	require.True(t, consumerDriver.started)
	require.NoError(t, manager.Stop())
	require.True(t, consumerDriver.stopped)
	require.True(t, producerDriver.closed)
}

// TestRedisDelayedMessageEnvelope 测试延迟消息唯一标识和兼容解析
func TestRedisDelayedMessageEnvelope(t *testing.T) {
	body := []byte(`{"name":"same"}`)

	first, err := queue.NewRedisDelayedMessage(body)
	require.NoError(t, err)
	second, err := queue.NewRedisDelayedMessage(body)
	require.NoError(t, err)

	require.NotEqual(t, string(first), string(second))
	require.Equal(t, body, queue.DecodeRedisDelayedMessage(string(first)))
	require.Equal(t, body, queue.DecodeRedisDelayedMessage(string(body)))
}

// TestRedisDelayedMessageAtomicClaim 测试延迟消息只能原子领取一次
func TestRedisDelayedMessageAtomicClaim(t *testing.T) {
	client := facade.Redis().Client()
	key := "test:queue:delayed:atomic"
	require.NoError(t, client.Del(t.Context(), key).Err())
	t.Cleanup(func() { _ = client.Del(context.Background(), key).Err() })

	for range 2 {
		body, err := queue.NewRedisDelayedMessage([]byte(`{"name":"same"}`))
		require.NoError(t, err)
		require.NoError(t, client.ZAdd(t.Context(), key, &redis.Z{
			Score:  float64(time.Now().Add(-time.Second).UnixMilli()),
			Member: string(body),
		}).Err())
	}

	first, err := queue.ClaimRedisDelayedMessages(t.Context(), client, key, float64(time.Now().UnixMilli()), 10)
	require.NoError(t, err)
	require.Len(t, first, 2)

	second, err := queue.ClaimRedisDelayedMessages(t.Context(), client, key, float64(time.Now().UnixMilli()), 10)
	require.NoError(t, err)
	require.Empty(t, second)
}

// TestQueueFacadeProducers 测试门面获取全部生产者
func TestQueueFacadeProducers(t *testing.T) {
	manager, err := queue.NewManager(func() *config.Config {
		return nil
	}, nil, nil)
	require.NoError(t, err)
	facade.App().SetQueue(manager)

	producer := &queueFacadeTestProducer{name: "queue_facade_test"}
	facade.Queue().RegisterProducer(producer)

	current := facade.Queue().Producer(producer.Name())
	require.Same(t, producer, current)

	producers := facade.Queue().Producers()
	found := false
	for _, item := range producers {
		if item.Name() == producer.Name() {
			found = true
			break
		}
	}
	require.True(t, found)
}

// TestQueueConsumerStatusZeroValue 测试消费者状态零值
func TestQueueConsumerStatusZeroValue(t *testing.T) {
	consumer := &queue.RedisConsumer{}
	require.Equal(t, queue.ConsumerStatusStopped, consumer.Status())
}
