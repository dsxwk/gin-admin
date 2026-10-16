package mcp

import (
	"context"
	"encoding/json"
	"fmt"
)

// ToolConfig 工具配置
type ToolConfig struct {
	Name        string
	Description string
	InputSchema InputSchema
	Auth        bool
}

// ToolAdapter 泛型工具适配器
type ToolAdapter[T any] struct {
	config ToolConfig
	call   func(ctx context.Context, userId int64, args *T) (any, error)
}

// NewToolAdapter 创建泛型工具适配器
func NewToolAdapter[T any](config ToolConfig, call func(ctx context.Context, userId int64, args *T) (any, error)) *ToolAdapter[T] {
	if config.Description == "" {
		config.Description = config.Name
	}

	return &ToolAdapter[T]{
		config: config,
		call:   call,
	}
}

// Name 获取工具名称
func (t *ToolAdapter[T]) Name() string {
	return t.config.Name
}

// Description 获取工具描述
func (t *ToolAdapter[T]) Description() string {
	return t.config.Description
}

// InputSchema 获取输入参数Schema
func (t *ToolAdapter[T]) InputSchema() InputSchema {
	return t.config.InputSchema
}

// IsAuth 是否需要鉴权
func (t *ToolAdapter[T]) IsAuth() bool {
	return t.config.Auth
}

// Call 执行工具调用
func (t *ToolAdapter[T]) Call(ctx context.Context, args map[string]any) (any, error) {
	return t.CallWithUser(ctx, 0, args)
}

// CallWithUser 带用户上下文执行工具调用
func (t *ToolAdapter[T]) CallWithUser(
	ctx context.Context,
	userId int64,
	args map[string]any,
) (any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if t.call == nil {
		return nil, fmt.Errorf("mcp tool [%s] handler not initialized", t.config.Name)
	}

	payload := new(T)
	if len(args) > 0 {
		body, err := json.Marshal(args)
		if err != nil {
			return nil, fmt.Errorf("mcp tool [%s] arguments encode failed: %w", t.config.Name, err)
		}
		if err = json.Unmarshal(body, payload); err != nil {
			return nil, fmt.Errorf("mcp tool [%s] arguments decode failed: %w", t.config.Name, err)
		}
	}

	return t.call(ctx, userId, payload)
}
