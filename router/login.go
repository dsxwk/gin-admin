package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// LoginRouter 登录路由
type LoginRouter struct{}

func init() {
	Register(&LoginRouter{})
}

// RegisterRoutes 注册路由
func (r *LoginRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		login v1.LoginController
	)

	router := routerGroup.Group("api/v1")
	{
		// 登录
		router.POST("/login", login.Login)
		// 刷新token
		router.POST("/refresh-token", login.RefreshToken)
		// 测试
		router.POST("/test", login.Test)
	}
}

// IsAuth 是否需要鉴权
func (r *LoginRouter) IsAuth() bool {
	return false
}
