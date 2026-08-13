package mcp

import (
	"context"
	"gin/app/service"
	"gin/pkg/serviceprovider/mcp"
)

func init() { mcp.Register(&RolePermission{}) }

// RolePermission 角色权限设置工具
type RolePermission struct{}

func (t *RolePermission) Name() string { return "set_role_permissions" }
func (t *RolePermission) Description() string {
	return "为指定角色设置权限,permissions为空时授予全部权限"
}

func (t *RolePermission) InputSchema() mcp.InputSchema {
	return mcp.InputSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"role":        {Type: "string", Description: "角色名称,如admin"},
			"permissions": {Type: "array", Description: "权限标识列表,如GET:/api/v1/user,为空授予全部权限"},
		},
		Required: []string{"role"},
	}
}

func (t *RolePermission) Call(ctx context.Context, args map[string]any) (any, error) {
	return t.CallWithUser(ctx, 0, args)
}

// CallWithUser 带用户上下文执行(实现mcp.UserContextCall接口)
func (t *RolePermission) CallWithUser(ctx context.Context, userId int64, args map[string]any) (any, error) {
	role, _ := args["role"].(string)
	if role == "" {
		return map[string]any{"error": "请指定角色名称"}, nil
	}

	// 解析权限列表
	var permissions []string
	if p, ok := args["permissions"].([]any); ok {
		for _, v := range p {
			if s, ok := v.(string); ok {
				permissions = append(permissions, s)
			}
		}
	}

	svc := service.PermissionService{}
	svc.WithContext(ctx)
	count, err := svc.SetRolePermissions(userId, role, permissions)
	if err != nil {
		return map[string]any{"error": err.Error()}, nil
	}

	return map[string]any{
		"role":        role,
		"granted":     count,
		"permissions": permissions,
		"message":     "权限设置完成",
	}, nil
}

// IsAuth 需要鉴权
func (t *RolePermission) IsAuth() bool {
	return true
}
