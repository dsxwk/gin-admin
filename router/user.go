package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// UserRouter 用户路由
type UserRouter struct{}

func init() {
	Register(&UserRouter{})
}

// RegisterRoutes 注册路由
func (r *UserRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		user v1.UserController
	)

	router := routerGroup.Group("api/v1/user")
	{
		// 列表
		router.GET("", user.List)
		// 创建
		router.POST("", user.Create)
		// 更新
		router.PUT("/:id", user.Update)
		// 删除
		router.DELETE("/:id", user.Delete)
		// 批量导入
		router.POST("/import", user.Import)
		// 更新密码
		router.PUT("/:id/password", user.Password)
		// 批量删除
		router.POST("/batch-delete", user.BatchDelete)
		// 详情
		router.GET("/:id", user.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *UserRouter) IsAuth() bool {
	return true
}
