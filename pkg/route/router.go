package route

import "github.com/gin-gonic/gin"

// Router 路由接口
type Router interface {
	Register(router *gin.RouterGroup)
	IsAuth() bool
}

// Mount 挂载路由模块并返回鉴权路由Key
func Mount(engine *gin.Engine, routerGroup *gin.RouterGroup, authMiddleware []gin.HandlerFunc, routers ...Router) []string {
	if engine == nil || routerGroup == nil {
		return nil
	}

	public := routerGroup.Group("")
	auth := routerGroup.Group("", authMiddleware...)
	permissionKeys := make([]string, 0)
	seen := make(map[string]struct{})

	for _, r := range routers {
		if r == nil {
			continue
		}

		if !r.IsAuth() {
			r.Register(public)
			continue
		}

		before := routeKeySet(engine.Routes())
		r.Register(auth)

		for _, route := range engine.Routes() {
			key := routeKey(route)
			if _, exists := before[key]; exists {
				continue
			}
			if _, exists := seen[key]; exists {
				continue
			}

			seen[key] = struct{}{}
			permissionKeys = append(permissionKeys, key)
		}
	}

	return permissionKeys
}

// routeKeySet 构建路由Key集合
func routeKeySet(routes []gin.RouteInfo) map[string]struct{} {
	keys := make(map[string]struct{}, len(routes))
	for _, route := range routes {
		keys[routeKey(route)] = struct{}{}
	}
	return keys
}

// routeKey 构建权限Key
func routeKey(route gin.RouteInfo) string {
	return route.Method + ":" + route.Path
}
