package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// DictRouter Dict
type DictRouter struct{}

func init() {
	Register(&DictRouter{})
}

// RegisterRoutes 注册路由
func (r *DictRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		dict v1.DictController
	)

	router := routerGroup.Group("/api/v1/dict")
	{
		// 列表
		router.GET("", dict.List)
		// 创建
		router.POST("", dict.Create)
		// 更新
		router.PUT("/:id", dict.Update)
		// 删除
		router.DELETE("/:id", dict.Delete)
		// 详情
		router.GET("/:id", dict.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *DictRouter) IsAuth() bool {
	return true
}
