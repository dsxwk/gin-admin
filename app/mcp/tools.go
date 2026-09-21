package mcp

import servicemcp "gin/pkg/serviceprovider/mcp"

// Tools 获取MCP工具列表
func Tools() []servicemcp.Tool {
	return []servicemcp.Tool{
		&CliExec{},
		&DepartmentQuery{},
		&DictQuery{},
		&RolePermission{},
		&OperatorStats{},
		&SystemConfigQuery{},
		&UserSearch{},
	}
}
