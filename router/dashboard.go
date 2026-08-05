package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// DashboardRouter 仪表盘
type DashboardRouter struct{}

func init() {
	Register(&DashboardRouter{})
}

// RegisterRoutes 注册路由
func (r *DashboardRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		dashboard v1.DashboardController
	)

	router := routerGroup.Group("/api/v1/dashboard")
	{
		// 卡片统计
		router.GET("/cards", dashboard.Cards)
		// 操作日志统计
		router.GET("/statistics", dashboard.Statistics)
		// 系统资源
		router.GET("/system-resource", dashboard.SystemResource)
	}
}

// IsAuth 是否需要鉴权
func (r *DashboardRouter) IsAuth() bool {
	return true
}
