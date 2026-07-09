package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// SystemConfigRouter SystemConfig
type SystemConfigRouter struct{}

func init() {
	Register(&SystemConfigRouter{})
}

// RegisterRoutes 注册路由
func (r *SystemConfigRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		systemConfig v1.SystemConfigController
	)

	router := routerGroup.Group("/api/v1/system-config")
	{
		// 列表
		router.GET("", systemConfig.List)
		// 保存配置
		router.PUT("", systemConfig.UpdateConfig)
		// 创建
		router.POST("", systemConfig.Create)
		// 更新
		router.PUT("/:id", systemConfig.Update)
		// 删除
		router.DELETE("/:id", systemConfig.Delete)
		// 详情
		router.GET("/:id", systemConfig.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *SystemConfigRouter) IsAuth() bool {
	return true
}
