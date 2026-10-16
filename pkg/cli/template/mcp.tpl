package {{.Package}}

import (
	"context"
	"gin/pkg/serviceprovider/mcp"
)

// {{.Name}}Args {{.Description}}参数
type {{.Name}}Args struct {
}

// {{.Name}}Tool {{.Description}}工具
type {{.Name}}Tool struct{}

// New{{.Name}}Tool 创建{{.Description}}
func New{{.Name}}Tool() *mcp.ToolAdapter[{{.Name}}Args] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "{{.ToolName}}",
			Description: "{{.Description}}",
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					// 待添加参数
				},
				Required: []string{},
			},
		},
		(&{{.Name}}Tool{}).Handle,
	)
}

// Handle 处理{{.Description}}
func (t *{{.Name}}Tool) Handle(_ context.Context, _ int64, args *{{.Name}}Args) (any, error) {
	// 待处理业务逻辑
	return map[string]any{}, nil
}
