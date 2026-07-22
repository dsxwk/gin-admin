package router

import (
	"github.com/gin-gonic/gin"
	"io"
)

// Router 路由接口
type Router interface {
	RegisterRoutes(router *gin.RouterGroup)
	IsAuth() bool // 是否需要鉴权
}

var routerRegister []Router

// Register 注册路由模块
func Register(r Router) {
	routerRegister = append(routerRegister, r)
}

// AutoLoads 自动注册
func AutoLoads(public *gin.RouterGroup, auth *gin.RouterGroup) {
	for _, r := range routerRegister {
		if r.IsAuth() {
			r.RegisterRoutes(auth)
		} else {
			r.RegisterRoutes(public)
		}
	}
}

// GenerateAuthPermissionKeys 提取需要鉴权的路由权限Key(格式: METHOD:PATH)
func GenerateAuthPermissionKeys() []string {
	// 临时禁用Gin debug输出,避免二次打印路由
	orig := gin.DefaultWriter
	gin.DefaultWriter = io.Discard
	defer func() {
		gin.DefaultWriter = orig
	}()

	engine := gin.New()
	group := engine.Group("")
	for _, r := range routerRegister {
		if r.IsAuth() {
			r.RegisterRoutes(group)
		}
	}
	return extractPermissionKeys(engine)
}

// extractPermissionKeys 从engine提取去重的METHOD:PATH权限Key
func extractPermissionKeys(engine *gin.Engine) []string {
	routes := engine.Routes()
	seen := make(map[string]bool)
	keys := make([]string, 0, len(routes))
	for _, route := range routes {
		key := route.Method + ":" + route.Path
		if !seen[key] {
			keys = append(keys, key)
			seen[key] = true
		}
	}
	return keys
}
