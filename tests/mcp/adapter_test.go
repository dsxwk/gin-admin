package mcp

import (
	"context"
	"testing"

	servicemcp "gin/pkg/serviceprovider/mcp"

	"github.com/stretchr/testify/require"
)

type adapterTestArgs struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// TestToolAdapter 测试工具适配器
func TestToolAdapter(t *testing.T) {
	var userId int64
	tool := servicemcp.NewToolAdapter(
		servicemcp.ToolConfig{
			Name:        "adapter_test",
			Description: "adapter test",
			Auth:        true,
			InputSchema: servicemcp.InputSchema{
				Type: "object",
				Properties: map[string]servicemcp.Property{
					"name": {Type: "string", Description: "name"},
				},
			},
		},
		func(_ context.Context, uid int64, args *adapterTestArgs) (any, error) {
			userId = uid
			return map[string]any{"name": args.Name, "count": args.Count}, nil
		},
	)

	require.Equal(t, "adapter_test", tool.Name())
	require.True(t, tool.IsAuth())
	require.Equal(t, "object", tool.InputSchema().Type)

	result, err := servicemcp.CallTool(context.Background(), 7, tool, map[string]any{
		"name":  "adapter",
		"count": 2,
	})
	require.NoError(t, err)
	require.Equal(t, int64(7), userId)
	require.Equal(t, map[string]any{"name": "adapter", "count": 2}, result)
}
