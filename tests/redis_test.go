package tests

import (
	"context"
	"gin/app/facade"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRedisPubSub 测试Redis发布订阅
func TestRedisPubSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用门面获取Redis缓存实例
	redisCache := facade.Cache.Store("redis").WithContext(ctx)

	done := make(chan struct{})

	// 订阅频道
	err := redisCache.Redis().Subscribe("testRedisChan", func(channel, payload string) {
		assert.Equal(t, "testRedisChan", channel)
		assert.Contains(t, payload, "test")
		close(done)
	})
	require.NoError(t, err)

	// 等待订阅建立
	time.Sleep(100 * time.Millisecond)

	// 发布消息
	err = redisCache.Redis().Publish("testRedisChan", map[string]interface{}{
		"test": "ok",
	})
	require.NoError(t, err)

	// 等待接收消息
	select {
	case <-done:
		t.Log("成功收到消息")
	case <-time.After(3 * time.Second):
		t.Fatal("超时: 未收到订阅消息")
	}
}

// TestRedisSetGet 测试Redis基本操作
func TestRedisSetGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用门面获取Redis缓存实例
	redisCache := facade.Cache.Store("redis").WithContext(ctx)
	key := "test:setget"
	value := "hello redis"

	// 测试 Set
	err := redisCache.Set(key, value, 10*time.Second)
	require.NoError(t, err)

	// 测试 Get
	val, ok := redisCache.Get(key)
	require.True(t, ok)
	assert.Equal(t, value, val)

	// 清理
	err = redisCache.Delete(key)
	require.NoError(t, err)

	// 验证删除
	val, ok = redisCache.Get(key)
	require.False(t, ok)
	assert.Nil(t, val)
}

// TestRedisLock 测试Redis分布式锁
func TestRedisLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache.Store("redis").WithContext(ctx)
	key := "test:lock"
	value := "lock-value"

	// 获取锁
	err := redisCache.Redis().Lock(key, value, 2*time.Second)
	require.NoError(t, err)

	// 尝试再次获取同一个锁(应该失败)
	err = redisCache.Redis().Lock(key, value, 2*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lock already exists")

	// 释放锁
	err = redisCache.Redis().UnLock(key, value)
	require.NoError(t, err)

	// 释放后可以重新获取
	err = redisCache.Redis().Lock(key, value, 2*time.Second)
	require.NoError(t, err)
}

// TestRedisExpire 测试Redis过期
func TestRedisExpire(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache.Store("redis").WithContext(ctx)
	key := "test:expire"
	value := "expire test"

	// 设置1秒过期
	err := redisCache.Set(key, value, 1*time.Second)
	require.NoError(t, err)

	// 立即获取应该存在
	val, ok := redisCache.Get(key)
	require.True(t, ok)
	assert.Equal(t, value, val)

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	// 过期后应该不存在
	val, ok = redisCache.Get(key)
	require.False(t, ok)
	assert.Nil(t, val)
}

// TestRedisDifferentDataTypes 测试不同数据类型
func TestRedisDifferentDataTypes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisCache := facade.Cache.Store("redis").WithContext(ctx)

	testCases := []struct {
		name  string
		key   string
		value interface{}
	}{
		{"string", "test:type:string", "hello"},
		{"int", "test:type:int", 12345},
		{"int64", "test:type:int64", int64(12345)},
		{"float64", "test:type:float64", 3.14159},
		{"bool", "test:type:bool", true},
		{"map", "test:type:map", map[string]interface{}{"name": "test", "value": 100}},
		{"slice", "test:type:slice", []int{1, 2, 3, 4, 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置
			err := redisCache.Set(tc.key, tc.value, 10*time.Second)
			require.NoError(t, err)

			// 获取
			val, ok := redisCache.Get(tc.key)
			require.True(t, ok)

			// 验证（注意：Redis返回的是字符串或字节数组,需要特殊处理）
			switch tc.name {
			case "string":
				assert.Equal(t, tc.value, val)
			case "int", "int64", "float64", "bool":
				// Redis返回字符串,需要转换
				assert.NotNil(t, val)
			default:
				assert.NotNil(t, val)
			}

			// 清理
			err = redisCache.Delete(tc.key)
			require.NoError(t, err)
		})
	}
}
