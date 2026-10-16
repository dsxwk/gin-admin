package serviceprovider

import (
	"context"
	"testing"

	serviceprovider "gin/pkg/serviceprovider"

	"github.com/stretchr/testify/require"
)

// TestShutdownRunner 测试通用关闭任务
func TestShutdownRunner(t *testing.T) {
	stopped := false
	runner := serviceprovider.NewShutdownRunner("test_shutdown", func() error {
		stopped = true
		return nil
	})

	require.Equal(t, "test_shutdown", runner.Name())
	require.NoError(t, runner.Stop())
	require.True(t, stopped)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.NoError(t, runner.Run(ctx))
}
