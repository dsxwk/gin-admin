package route

import "github.com/gin-gonic/gin"

// Router 路由接口
type Router interface {
	Register(router *gin.RouterGroup)
	IsAuth() bool
}

// Mount 挂载路由模块
func Mount(routerGroup *gin.RouterGroup, authMiddleware []gin.HandlerFunc, routers ...Router) {
	if routerGroup == nil {
		return
	}

	public := routerGroup.Group("")
	auth := routerGroup.Group("", authMiddleware...)

	for _, r := range routers {
		if r == nil {
			continue
		}

		if r.IsAuth() {
			r.Register(auth)
			continue
		}

		r.Register(public)
	}
}
