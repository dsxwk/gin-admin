package router

import (
	"gin/pkg/route"
)

// All 路由模块
func All() []route.Router {
	return []route.Router{
		&LoginRouter{},
		&AgentRouter{},
		&ArticleRouter{},
		&ConfigCategoryRouter{},
		&DashboardRouter{},
		&DepartmentRouter{},
		&DictRouter{},
		&ImportRecordsRouter{},
		&MenuRouter{},
		&OperatorLogRouter{},
		&PermissionRouter{},
		&RoleRouter{},
		&SystemConfigRouter{},
		&UserRouter{},
		&McpRouter{},
	}
}
