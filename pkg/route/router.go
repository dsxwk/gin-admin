package route

import "github.com/gin-gonic/gin"

// Router 路由接口
type Router interface {
	RegisterRoutes(router *gin.RouterGroup)
	IsAuth() bool
}
