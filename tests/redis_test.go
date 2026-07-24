package tests

import (
	"context"
	"gin/app/facade"
	"gin/app/queue/consumer"
	_ "gin/app/queue/producer"
	"gin/common/ctxkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestRedisQueuePublish Redis队列消息发布和消费
func TestRedisQueuePublish(t *testing.T) {
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-redis-queue")

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
	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-redis-delay")

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

	consumers := facade.Queue().GetAllConsumers()
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

	redisCache := facade.Cache("redis").WithContext(ctx)
	key := "test:cache:lock"
	value := "lock-owner"
	defer func() { _ = redisCache.Redis().UnLock(key, value) }()

	err := redisCache.Redis().Lock(key, value, 2*time.Second)
	require.NoError(t, err, "获取锁失败")

	err = redisCache.Redis().Lock(key, value, 2*time.Second)
	require.Error(t, err, "重复获取同一锁应该失败")
	assert.Contains(t, err.Error(), "lock already exists")

	err = redisCache.Redis().UnLock(key, value)
	require.NoError(t, err, "释放锁失败")

	err = redisCache.Redis().Lock(key, value, 2*time.Second)
	require.NoError(t, err, "释放后重新获取锁失败")
}

// TestRedisCacheSetOps Redis集合操作
func TestRedisCacheSetOps(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache("redis").WithContext(ctx)
	key := "test:cache:set"
	defer func() { _ = redisCache.Redis().Delete(key) }()

	err := redisCache.Redis().SAdd(key, "a", "b", "c")
	require.NoError(t, err, "SAdd失败")

	isMember, err := redisCache.Redis().SIsMember(key, "a")
	require.NoError(t, err)
	assert.True(t, isMember, "成员a应存在")

	isMember, err = redisCache.Redis().SIsMember(key, "x")
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
		value interface{}
	}{
		{"string", "test:type:string", "hello"},
		{"int", "test:type:int", 12345},
		{"float", "test:type:float", 3.14159},
		{"bool", "test:type:bool", true},
		{"map", "test:type:map", map[string]interface{}{"name": "test", "count": 100}},
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
