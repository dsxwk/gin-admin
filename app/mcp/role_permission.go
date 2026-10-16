package mcp

import (
	"context"
	"gin/app/service"
	"gin/pkg/serviceprovider/mcp"
)

// RolePermissionArgs 角色权限设置参数
type RolePermissionArgs struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// RolePermissionTool 角色权限设置工具
type RolePermissionTool struct{}

// NewRolePermissionTool 创建角色权限设置工具
func NewRolePermissionTool() *mcp.ToolAdapter[RolePermissionArgs] {
	return mcp.NewToolAdapter(
		mcp.ToolConfig{
			Name:        "set_role_permissions",
			Description: "为指定角色设置权限,permissions为空时授予全部权限",
			Auth:        true,
			InputSchema: mcp.InputSchema{
				Type: "object",
				Properties: map[string]mcp.Property{
					"role":        {Type: "string", Description: "角色名称,如admin"},
					"permissions": {Type: "array", Description: "权限标识列表,如GET:/api/v1/user,为空授予全部权限"},
				},
				Required: []string{"role"},
			},
		},
		(&RolePermissionTool{}).Handle,
	)
}

// Handle 处理角色权限设置工具
func (t *RolePermissionTool) Handle(ctx context.Context, userId int64, args *RolePermissionArgs) (any, error) {
	if args.Role == "" {
		return map[string]any{"error": "请指定角色名称"}, nil
	}

	svc := service.PermissionService{}
	count, err := svc.SetRolePermissions(ctx, userId, args.Role, args.Permissions)
	if err != nil {
		return map[string]any{"error": err.Error()}, nil
	}

	return map[string]any{
		"role":        args.Role,
		"granted":     count,
		"permissions": args.Permissions,
		"message":     "权限设置完成",
	}, nil
}
