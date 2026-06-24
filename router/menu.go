package router

import (
	v1 "gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// MenuRouter 菜单路由
type MenuRouter struct{}

func init() {
	Register(&MenuRouter{})
}

// RegisterRoutes 注册路由
func (r *MenuRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		menu v1.MenuController
	)

	router := routerGroup.Group("api/v1")
	{
		// 列表
		router.GET("/menu", menu.List)

		// 角色菜单
		router.GET("/role/:id/menu", menu.RoleMenu)

		// 详情
		router.GET("/menu/:id", menu.Detail)

		// 创建
		router.POST("/menu", menu.Create)

		// 更新
		router.PUT("/menu/:id", menu.Update)

		// 删除
		router.DELETE("/menu/:id", menu.Delete)

		// 菜单功能
		router.GET("/menu/:id/action", menu.Action)

		// 菜单功能详情
		router.GET("/menu/:id/action/:actionId", menu.ActionDetail)

		// 新增菜单功能
		router.POST("/menu/:id/action", menu.CreateAction)

		// 更新菜单功能
		router.PUT("/menu/:id/action/:actionId", menu.UpdateAction)

		// 删除菜单功能
		router.DELETE("/menu/:id/action/:actionId", menu.DeleteAction)
	}
}

// IsAuth 是否需要鉴权
func (r *MenuRouter) IsAuth() bool {
	return true
}
