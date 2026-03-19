package tests

import (
	"context"
	"gin/pkg/cache"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name  string
		cache *cache.CacheProxy
		key   string
	}{
		{"redis", cache.NewRedisCache().WithContext(ctx), "redis_test"},
		{"disk", cache.NewDiskCache().WithContext(ctx), "disk_test"},
		{"memory", cache.NewMemoryCache().WithContext(ctx), "memory_test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cache.Set(tt.key, 123, 10*time.Second)
			require.NoError(t, err)

			val, ok := tt.cache.Get(tt.key)
			require.True(t, ok)

			switch v := val.(type) {
			case string:
				valInt, err := strconv.Atoi(v)
				require.NoError(t, err)
				require.Equal(t, 123, valInt)
			case []byte:
				valInt, err := strconv.Atoi(string(v))
				require.NoError(t, err)
				require.Equal(t, 123, valInt)
			case int:
				require.Equal(t, 123, v)
			default:
				t.Fatalf("unexpected type %T", v)
			}
		})
	}
}
