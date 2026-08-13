package {{.Package}}

import (
    "context"
    "gin/pkg/serviceprovider/mcp"
)

func init() { mcp.Register(&{{.Name}}{}) }

// {{.Name}} {{.Description}}工具
type {{.Name}} struct{}

func (t *{{.Name}}) Name() string { return "{{.ToolName}}" }
func (t *{{.Name}}) Description() string { return "{{.Description}}" }

func (t *{{.Name}}) InputSchema() mcp.InputSchema {
    return mcp.InputSchema{
        Type: "object",
        Properties: map[string]mcp.Property{
            // TODO: 添加参数
        },
        Required: []string{},
    }
}

func (t *{{.Name}}) Call(ctx context.Context, args map[string]any) (any, error) {
    // TODO: 实现工具逻辑
    return map[string]any{}, nil
}