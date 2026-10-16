package mcp

import servicemcp "gin/pkg/serviceprovider/mcp"

// All MCP工具列表
func All() []servicemcp.Tool {
	return []servicemcp.Tool{
		NewCliExecTool(),
		NewDepartmentQueryTool(),
		NewDictQueryTool(),
		NewRolePermissionTool(),
		NewOperatorStatsTool(),
		NewSystemConfigQueryTool(),
		NewUserSearchTool(),
	}
}
