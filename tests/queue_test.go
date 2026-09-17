package tests

import (
	"context"
	"gin/app/facade"
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
func (p *queueFacadeTestProducer) Publish(ctx context.Context, msg any) error {
	return nil
}

// Close 关闭生产者
func (p *queueFacadeTestProducer) Close() error {
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

// TestQueueRegistryRegisterFactory 测试注册表延迟创建工厂
func TestQueueRegistryRegisterFactory(t *testing.T) {
	registry := queue.NewRegistry[queueRegistryTestItem]()
	item := queueRegistryTestItem{name: "queue_registry_factory_test"}
	created := false

	registry.RegisterFactory(func() queueRegistryTestItem {
		created = true
		return item
	})

	require.False(t, created)
	require.NoError(t, registry.RegisterFactories())
	require.True(t, created)

	current, ok := registry.Get(item.Name())
	require.True(t, ok)
	require.Equal(t, item, current)
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

// TestQueueFacadeGetAllProducers 测试门面获取全部生产者
func TestQueueFacadeGetAllProducers(t *testing.T) {
	producer := &queueFacadeTestProducer{name: "queue_facade_test"}
	facade.Queue().Register(producer)

	current := facade.Queue().Producer(producer.Name())
	require.Same(t, producer, current)

	producers := facade.Queue().GetAllProducers()
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
