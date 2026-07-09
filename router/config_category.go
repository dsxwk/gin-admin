package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// ConfigCategoryRouter ConfigCategory
type ConfigCategoryRouter struct{}

func init() {
	Register(&ConfigCategoryRouter{})
}

// RegisterRoutes 注册路由
func (r *ConfigCategoryRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		configCategory v1.ConfigCategoryController
	)

	router := routerGroup.Group("/api/v1/config-category")
	{
		// 列表
		router.GET("", configCategory.List)
		// 创建
		router.POST("", configCategory.Create)
		// 更新
		router.PUT("/:id", configCategory.Update)
		// 删除
		router.DELETE("/:id", configCategory.Delete)
		// 详情
		router.GET("/:id", configCategory.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *ConfigCategoryRouter) IsAuth() bool {
	return true
}
