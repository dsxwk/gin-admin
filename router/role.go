package router

import (
	v1 "gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// RoleRouter 角色路由
type RoleRouter struct{}

func init() {
	Register(&RoleRouter{})
}

// RegisterRoutes 注册路由
func (r *RoleRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		role v1.RoleController
	)

	router := routerGroup.Group("api/v1")
	{
		// 列表
		router.GET("role", role.List)
	}
}

// IsAuth 是否需要鉴权
func (r *RoleRouter) IsAuth() bool {
	return true
}
