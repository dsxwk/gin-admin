package tests

import (
	"gin/app/facade"
	"gin/pkg/serviceprovider/cache"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestCacheSetGet 测试缓存设置和获取
func TestCacheSetGet(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name      string
		cacheType string
		key       string
	}{
		{"redis", "redis", "redis_test"},
		{"memory", "memory", "memory_test"},
		{"disk", "disk", "disk_test"},
		{"default", "", "default_test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 获取带上下文的缓存实例
			var _cache *cache.CacheProxy
			_cache = facade.Cache(tt.cacheType).WithContext(ctx)

			// 设置缓存
			err := _cache.Set(tt.key, 123, 10*time.Second)
			require.NoError(t, err)

			// 获取缓存
			val, ok := _cache.Get(tt.key)
			require.True(t, ok)

			// 验证值
			switch v := val.(type) {
			case string:
				require.Equal(t, "123", v)
			case []byte:
				require.Equal(t, []byte("123"), v) // 或 string(v) == "123"
			case int:
				require.Equal(t, 123, v)
			case int64:
				require.Equal(t, int64(123), v)
			case float64:
				require.Equal(t, float64(123), v)
			default:
				t.Fatalf("unexpected type %T, value: %v", v, v)
			}
		})
	}
}

// TestCacheExpiration 测试缓存过期
func TestCacheExpiration(t *testing.T) {
	ctx := t.Context()
	_cache := facade.Cache().WithContext(ctx)
	key := "expire_test"

	// 设置1秒过期的缓存
	err := _cache.Set(key, "value", 1*time.Second)
	require.NoError(t, err)

	// 立即获取应该存在
	val, ok := _cache.Get(key)
	require.True(t, ok)
	require.Equal(t, "value", val)

	// 等待1.5秒后过期
	time.Sleep(1500 * time.Millisecond)

	// 再次获取应该不存在
	val, ok = _cache.Get(key)
	require.False(t, ok)
	require.Nil(t, val)
}

// TestCacheDelete 测试缓存删除
func TestCacheDelete(t *testing.T) {
	ctx := t.Context()
	_cache := facade.Cache().WithContext(ctx)
	key := "delete_test"

	// 设置缓存
	err := _cache.Set(key, "value", 10*time.Second)
	require.NoError(t, err)

	// 确认存在
	val, ok := _cache.Get(key)
	require.True(t, ok)
	require.Equal(t, "value", val)

	// 删除缓存
	err = _cache.Delete(key)
	require.NoError(t, err)

	// 确认已删除
	val, ok = _cache.Get(key)
	require.False(t, ok)
	require.Nil(t, val)
}

// TestCacheDifferentTypes 测试不同类型缓存
func TestCacheDifferentTypes(t *testing.T) {
	ctx := t.Context()

	// 测试各种类型
	testCases := []struct {
		name  string
		value interface{}
		check func(t *testing.T, got interface{})
	}{
		{"string", "hello world", func(t *testing.T, got interface{}) {
			require.Equal(t, "hello world", got)
		}},
		{"int", 12345, func(t *testing.T, got interface{}) {
			// JSON序列化后int会变成float64
			switch v := got.(type) {
			case int:
				require.Equal(t, 12345, v)
			case int64:
				require.Equal(t, int64(12345), v)
			case float64:
				require.Equal(t, float64(12345), v)
			default:
				t.Fatalf("unexpected type: %T", v)
			}
		}},
		{"int64", int64(12345), func(t *testing.T, got interface{}) {
			switch v := got.(type) {
			case int64:
				require.Equal(t, int64(12345), v)
			case float64:
				require.Equal(t, float64(12345), v)
			default:
				t.Fatalf("unexpected type: %T", v)
			}
		}},
		{"float64", 3.14159, func(t *testing.T, got interface{}) {
			require.Equal(t, 3.14159, got)
		}},
		{"bool", true, func(t *testing.T, got interface{}) {
			require.Equal(t, true, got)
		}},
		{"map", map[string]interface{}{"name": "test", "value": 100}, func(t *testing.T, got interface{}) {
			m, ok := got.(map[string]interface{})
			require.True(t, ok)
			require.Equal(t, "test", m["name"])
			// map中的数字也会变成float64
			switch v := m["value"].(type) {
			case int:
				require.Equal(t, 100, v)
			case float64:
				require.Equal(t, float64(100), v)
			default:
				t.Fatalf("unexpected type for value: %T", v)
			}
		}},
		{"slice", []int{1, 2, 3, 4, 5}, func(t *testing.T, got interface{}) {
			// JSON序列化后slice会变成[]interface{}
			s, ok := got.([]interface{})
			require.True(t, ok)
			require.Len(t, s, 5)
			for i, v := range s {
				switch v.(type) {
				case int:
					require.Equal(t, i+1, v)
				case float64:
					require.Equal(t, float64(i+1), v)
				default:
					t.Fatalf("unexpected type at index %d: %T", i, v)
				}
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_cache := facade.Cache("memory").WithContext(ctx)
			key := "type_test_" + tc.name

			err := _cache.Set(key, tc.value, 10*time.Second)
			require.NoError(t, err)

			val, ok := _cache.Get(key)
			require.True(t, ok)
			tc.check(t, val)
		})
	}
}

// TestCacheWithContext 测试带上下文的缓存
func TestCacheWithContext(t *testing.T) {
	ctx := t.Context()

	// 创建带上下文的缓存
	_cache := facade.Cache().WithContext(ctx)

	// 测试设置和获取
	key := "context_test"
	err := _cache.Set(key, "context_value", 10*time.Second)
	require.NoError(t, err)

	val, ok := _cache.Get(key)
	require.True(t, ok)
	require.Equal(t, "context_value", val)
}

// TestCacheStoreMethods 测试不同的Store方法
func TestCacheStoreMethods(t *testing.T) {
	ctx := t.Context()
	key := "store_method_test"
	value := "test_value"

	tests := []struct {
		name  string
		cache *cache.CacheProxy
	}{
		{"Store()", facade.Cache().WithContext(ctx)},
		{"Redis()", facade.Cache("redis").WithContext(ctx)},
		{"Memory()", facade.Cache("memory").WithContext(ctx)},
		{"Disk()", facade.Cache("disk").WithContext(ctx)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cache.Set(key, value, 10*time.Second)
			require.NoError(t, err)

			val, ok := tt.cache.Get(key)
			require.True(t, ok)

			// 兼容不同类型
			switch v := val.(type) {
			case string:
				require.Equal(t, value, v)
			case []byte:
				require.Equal(t, value, string(v))
			default:
				require.Equal(t, value, v)
			}

			_ = tt.cache.Delete(key)
		})
	}
}
