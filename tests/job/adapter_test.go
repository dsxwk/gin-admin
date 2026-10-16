package job

import (
	"context"
	"encoding/json"
	"testing"

	servicejob "gin/pkg/serviceprovider/job"

	"github.com/stretchr/testify/require"
)

type adapterTestPayload struct {
	Name string `json:"name"`
}

type adapterContextKey string

// TestJobAdapter 测试任务适配器
func TestJobAdapter(t *testing.T) {
	handled := false
	adapter := servicejob.NewAdapter(
		servicejob.Config{
			Name:        "adapter_test",
			Description: "adapter test",
			Connection:  "sync",
			Retry:       1,
			Delay:       100,
		},
		func(_ context.Context, payload *adapterTestPayload) error {
			handled = payload.Name == "adapter"
			return nil
		},
	)

	require.Equal(t, "adapter_test", adapter.Name())
	require.Equal(t, 1, adapter.Retry())
	require.Equal(t, int64(100), adapter.Delay())
	require.IsType(t, &adapterTestPayload{}, adapter.NewPayload())
	require.NoError(t, adapter.Handle(&adapterTestPayload{Name: "adapter"}))
	require.True(t, handled)
}

// TestJobAdapterContext 测试任务上下文传递
func TestJobAdapterContext(t *testing.T) {
	key := adapterContextKey("trace")
	handled := false
	adapter := servicejob.NewAdapter(
		servicejob.Config{
			Name:       "adapter_context_test",
			Connection: "sync",
		},
		func(ctx context.Context, payload *adapterTestPayload) error {
			handled = ctx.Value(key) == "context" && payload.Name == "adapter"
			return nil
		},
	)

	payload, err := json.Marshal(adapterTestPayload{Name: "adapter"})
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), key, "context")
	require.NoError(t, servicejob.Execute(ctx, adapter, servicejob.NewMessage("adapter_context_test", payload, 0)))
	require.True(t, handled)
}
