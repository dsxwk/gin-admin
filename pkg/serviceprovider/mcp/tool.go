package mcp

import "context"

// Tool MCP工具接口,所有工具需实现此接口
type Tool interface {
	Name() string                                               // 工具名称
	Description() string                                        // 工具描述
	InputSchema() InputSchema                                   // 输入参数Schema
	Call(ctx context.Context, args map[string]any) (any, error) // 执行调用
}

// AuthTool 需要鉴权的工具接口,实现此接口表示调用前必须校验token
type AuthTool interface {
	IsAuth() bool // 是否需要鉴权
}

// UserContextCall 带用户上下文的工具调用接口,需要当前用户信息的工具实现此接口
type UserContextCall interface {
	CallWithUser(ctx context.Context, userId int64, args map[string]any) (any, error)
}

// ToolCaller 工具调用器,返回handled表示已处理该工具
type ToolCaller func(ctx context.Context, userId int64, tool Tool, args map[string]any) (result any, handled bool, err error)

var toolCallers []ToolCaller

// RegisterToolCaller 注册工具调用器
func RegisterToolCaller(caller ToolCaller) {
	toolCallers = append(toolCallers, caller)
}

// CallTool 调用工具,支持带用户上下文的工具
func CallTool(ctx context.Context, userId int64, tool Tool, args map[string]any) (any, error) {
	for _, caller := range toolCallers {
		result, handled, err := caller(ctx, userId, tool, args)
		if handled {
			return result, err
		}
	}
	return tool.Call(ctx, args)
}

func init() {
	RegisterToolCaller(func(ctx context.Context, userId int64, tool Tool, args map[string]any) (any, bool, error) {
		if uc, ok := tool.(UserContextCall); ok {
			result, err := uc.CallWithUser(ctx, userId, args)
			return result, true, err
		}
		return nil, false, nil
	})
}
