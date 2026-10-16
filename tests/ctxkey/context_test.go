package ctxkey

import (
	"context"
	"testing"
	"time"

	servicectxkey "gin/pkg/ctxkey"

	"github.com/stretchr/testify/require"
)

// TestContextValue 测试标准Context读写
func TestContextValue(t *testing.T) {
	startTime := time.Now()
	ctx := context.Background()
	ctx = servicectxkey.WithValue(ctx, servicectxkey.UserIdKey, int64(1))
	ctx = servicectxkey.WithValue(ctx, servicectxkey.TraceIDKey, "trace-id")
	ctx = servicectxkey.WithValue(ctx, servicectxkey.LangKey, "zh")
	ctx = servicectxkey.WithValue(ctx, servicectxkey.StartTimeKey, startTime)
	ctx = servicectxkey.WithValue(ctx, servicectxkey.MsKey, 12.5)

	require.Equal(t, int64(1), servicectxkey.UserID(ctx))
	require.Equal(t, "trace-id", servicectxkey.TraceID(ctx))
	require.Equal(t, "zh", servicectxkey.Lang(ctx))
	require.Equal(t, startTime, servicectxkey.StartTime(ctx))
	require.Equal(t, 12.5, servicectxkey.Ms(ctx))
}

// TestContextValueNil 测试nilContext读取
func TestContextValueNil(t *testing.T) {
	require.Nil(t, servicectxkey.Value(nil, servicectxkey.TraceIDKey))
	require.Equal(t, int64(0), servicectxkey.UserID(nil))
	require.Equal(t, "unknown", servicectxkey.TraceID(nil))
	require.Empty(t, servicectxkey.Lang(nil))
	require.True(t, servicectxkey.StartTime(nil).IsZero())
	require.Zero(t, servicectxkey.Ms(nil))
}
