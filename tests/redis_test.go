package tests

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/app/queue/consumer"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/cache"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRedisQueuePublish Redis队列消息发布和消费
func TestRedisQueuePublish(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIDKey, "test-redis-queue")

	time.Sleep(500 * time.Millisecond)

	producer := facade.Queue().Producer("redis_demo")
	if producer == nil {
		t.Skip("Redis生产者未注册")
	}

	testCases := []struct {
		name    string
		payload consumer.RedisDemoPayload
	}{
		{"task_1", consumer.RedisDemoPayload{Name: "redis_test_1"}},
		{"task_2", consumer.RedisDemoPayload{Name: "redis_test_2"}},
		{"task_3", consumer.RedisDemoPayload{Name: "redis_test_3"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := producer.Publish(ctx, tc.payload)
			require.NoError(t, err, "发送消息失败: %s", tc.name)
			t.Logf("发送成功: %+v", tc.payload)
		})
	}

	time.Sleep(2 * time.Second)
	t.Log("Redis队列消息测试完成")
}

// TestRedisQueueDelayPublish Redis队列延迟消息发布和消费
func TestRedisQueueDelayPublish(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIDKey, "test-redis-delay")

	time.Sleep(500 * time.Millisecond)

	producer := facade.Queue().Producer("redis_delay_demo")
	if producer == nil {
		t.Skip("Redis延迟生产者未注册")
	}

	startTime := time.Now()
	t.Logf("开始发送延迟消息: %v", startTime)

	err := producer.Publish(ctx, consumer.RedisDelayDemoPayload{Name: "redis_delay_1"})
	require.NoError(t, err, "发送延迟消息1失败")

	err = producer.Publish(ctx, consumer.RedisDelayDemoPayload{Name: "redis_delay_2"})
	require.NoError(t, err, "发送延迟消息2失败")

	t.Log("延迟消息已发送, 等待6秒后消费...")

	time.Sleep(6 * time.Second)
	t.Logf("测试完成, 耗时: %v", time.Since(startTime))
}

// TestRedisQueueStatus 测试Redis消费者状态查询
func TestRedisQueueStatus(t *testing.T) {
	cfg := facade.Config()

	consumers := facade.Queue().Consumers()
	if len(consumers) == 0 {
		t.Skip("未注册消费者")
	}

	for _, c := range consumers {
		name := c.Name()
		status := c.Status()
		enabled := c.Enabled(cfg)

		t.Run(name, func(t *testing.T) {
			t.Logf("消费者: %s, 状态: %s, 启用: %v", name, status, enabled)
			assert.NotEmpty(t, name, "消费者名称不能为空")
		})
	}
}

// TestRedisCacheSetGet Redis基本读写操作
func TestRedisCacheSetGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)
	key := "test:cache:setget"
	value := "hello redis"

	err := redisCache.Set(key, value, 10*time.Second)
	require.NoError(t, err, "Set失败")
	defer func() { _ = redisCache.Delete(key) }()

	val, ok := redisCache.Get(key)
	require.True(t, ok, "键值不存在")
	assert.Equal(t, value, val, "读取值不匹配")

	err = redisCache.Delete(key)
	require.NoError(t, err, "Delete失败")

	val, ok = redisCache.Get(key)
	assert.False(t, ok, "删除后应不存在")
	assert.Nil(t, val)
}

// TestRedisCacheExpire Redis过期测试
func TestRedisCacheExpire(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)
	key := "test:cache:expire"
	value := "expire soon"

	err := redisCache.Set(key, value, 1*time.Second)
	require.NoError(t, err, "Set失败")

	val, ok := redisCache.Get(key)
	require.True(t, ok)
	assert.Equal(t, value, val)

	time.Sleep(1100 * time.Millisecond)

	val, ok = redisCache.Get(key)
	assert.False(t, ok, "过期后应不存在")
	assert.Nil(t, val)
}

// TestRedisCacheLock Redis分布式锁
func TestRedisCacheLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Redis().WithContext(ctx)
	key := "test:cache:lock"

	lock, err := redisCache.Lock(ctx, key, 2*time.Second)
	require.NoError(t, err, "获取锁失败")
	defer func() { _ = lock.Release() }()

	_, err = redisCache.Lock(ctx, key, 2*time.Second)
	require.ErrorIs(t, err, cache.ErrLockExists)

	require.NoError(t, lock.Release(), "释放锁失败")

	next, err := redisCache.Lock(ctx, key, 2*time.Second)
	require.NoError(t, err, "释放后重新获取锁失败")
	require.NoError(t, next.Release())
}

// TestRedisCacheLockWithWatchdog Redis锁自动续期测试
func TestRedisCacheLockWithWatchdog(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)
	key := "test:cache:lock-watchdog"
	defer func() { _ = redisCache.Delete(key) }()

	lock, err := redisCache.Lock(ctx, key, 100*time.Millisecond)
	require.NoError(t, err, "获取锁失败")

	time.Sleep(250 * time.Millisecond)

	_, err = redisCache.Lock(ctx, key, 100*time.Millisecond)
	require.ErrorIs(t, err, cache.ErrLockExists, "看门狗未保持锁")
	require.NoError(t, lock.Release(), "释放锁失败")

	next, err := redisCache.Lock(ctx, key, 100*time.Millisecond)
	require.NoError(t, err, "释放后重新获取锁失败")
	require.NoError(t, next.Release())
}

// TestRedisCacheLockLost Redis锁丢失测试
func TestRedisCacheLockLost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)
	key := "test:cache:lock-lost"

	lock, err := redisCache.Lock(ctx, key, time.Second)
	require.NoError(t, err, "获取锁失败")
	require.NoError(t, redisCache.Delete(key))

	require.ErrorIs(t, lock.Release(), cache.ErrLockNotOwned)
}

// TestRedisCacheSerializationSymmetric Redis序列化对称性测试
func TestRedisCacheSerializationSymmetric(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)
	testCases := []struct {
		name  string
		key   string
		value any
		check func(t *testing.T, value any)
	}{
		{"numeric-string", "test:serialize:string", "12345", func(t *testing.T, value any) {
			assert.Equal(t, "12345", value)
		}},
		{"integer", "test:serialize:int", 12345, func(t *testing.T, value any) {
			assert.Equal(t, int64(12345), value)
		}},
		{"map", "test:serialize:map", map[string]any{"count": 100}, func(t *testing.T, value any) {
			result, ok := value.(map[string]any)
			require.True(t, ok)
			assert.Equal(t, int64(100), result["count"])
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, redisCache.Set(tc.key, tc.value, 10*time.Second))
			defer func() { _ = redisCache.Delete(tc.key) }()

			value, ok := redisCache.Get(tc.key)
			require.True(t, ok)
			tc.check(t, value)
		})
	}
}

// TestRedisCacheSubscribeConcurrent Redis并发订阅测试
func TestRedisCacheSubscribeConcurrent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisCache := facade.Redis().WithContext(ctx)
	const count = 16

	var waitGroup sync.WaitGroup
	errs := make(chan error, count)
	waitGroup.Add(count)

	for index := range count {
		go func(index int) {
			defer waitGroup.Done()

			channel := fmt.Sprintf("test:cache:subscribe:%d", index)
			received := make(chan struct{})
			if err := redisCache.Subscribe(channel, func(_ string, _ string) {
				close(received)
			}); err != nil {
				errs <- err
				return
			}
			defer func() { _ = redisCache.Unsubscribe(channel) }()

			if err := redisCache.Publish(channel, "payload"); err != nil {
				errs <- err
				return
			}

			select {
			case <-received:
			case <-time.After(2 * time.Second):
				errs <- fmt.Errorf("channel %s receive timeout", channel)
			}
		}(index)
	}

	waitGroup.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
}

// TestRedisCacheSetOps Redis集合操作
func TestRedisCacheSetOps(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Redis().WithContext(ctx)
	key := "test:cache:set"
	defer func() { _ = redisCache.Delete(key) }()

	err := redisCache.SAdd(key, "a", "b", "c")
	require.NoError(t, err, "SAdd失败")

	isMember, err := redisCache.SIsMember(key, "a")
	require.NoError(t, err)
	assert.True(t, isMember, "成员a应存在")

	isMember, err = redisCache.SIsMember(key, "x")
	require.NoError(t, err)
	assert.False(t, isMember, "成员x不应存在")
}

// TestRedisCacheDataTypes Redis多种数据类型读写
func TestRedisCacheDataTypes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)

	testCases := []struct {
		name  string
		key   string
		value any
	}{
		{"string", "test:type:string", "hello"},
		{"int", "test:type:int", 12345},
		{"float", "test:type:float", 3.14159},
		{"bool", "test:type:bool", true},
		{"map", "test:type:map", map[string]any{"name": "test", "count": 100}},
		{"slice", "test:type:slice", []string{"a", "b", "c"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := redisCache.Set(tc.key, tc.value, 10*time.Second)
			require.NoError(t, err, "Set失败")
			defer func() { _ = redisCache.Delete(tc.key) }()

			val, ok := redisCache.Get(tc.key)
			require.True(t, ok, "读取失败")
			assert.NotNil(t, val, "读取值不能为空")
		})
	}
}
