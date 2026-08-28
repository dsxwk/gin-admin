package router

import (
	"gin/app/controller/v1"

	"github.com/gin-gonic/gin"
)

// DepartmentRouter Department
type DepartmentRouter struct{}

func init() {
	Register(&DepartmentRouter{})
}

// RegisterRoutes 注册路由
func (r *DepartmentRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		department v1.DepartmentController
	)

	router := routerGroup.Group("/api/v1/department")
	{
		// 列表
		router.GET("", department.List)
		// 创建
		router.POST("", department.Create)
		// 更新
		router.PUT("/:id", department.Update)
		// 删除
		router.DELETE("/:id", department.Delete)
		// 详情
		router.GET("/:id", department.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *DepartmentRouter) IsAuth() bool {
	return true
}
