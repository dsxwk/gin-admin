package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// PermissionRouter Permission
type PermissionRouter struct{}

func init() {
	Register(&PermissionRouter{})
}

// RegisterRoutes 注册路由
func (r *PermissionRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		permission v1.PermissionController
	)

	router := routerGroup.Group("/api/v1/permission")
	{
		// 列表
		router.GET("", permission.List)
	}
}

// IsAuth 是否需要鉴权
func (r *PermissionRouter) IsAuth() bool {
	return true
}
