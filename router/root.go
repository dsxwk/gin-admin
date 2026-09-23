package router

import (
	"gin/app/facade"
	"gin/app/middleware"
	_ "gin/docs"
	"gin/pkg"
	"gin/pkg/errcode"
	"gin/pkg/route"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouters 加载路由
func NewRouters(router *gin.Engine) {
	cfg := facade.Config()
	if cfg == nil {
		return
	}

	timeoutMiddleware := middleware.Timeout{}.Handle(cfg.App.Timeout)
	loggerMiddleware := middleware.Logger{}.Handle()
	corsMiddleware := middleware.Cors{}.Handle()
	jwtMiddleware := middleware.Jwt{}.Handle()
	recoverMiddleware := middleware.Recover{}.Handle()
	rateLimitMiddleware := middleware.RateLimit{}
	permissionMiddleware := middleware.Permission{}.Handle()
	operatorLogMiddleware := middleware.OperatorLog{}.Handle()

	// 全局中间件
	router.Use(corsMiddleware, timeoutMiddleware, loggerMiddleware, recoverMiddleware, operatorLogMiddleware)

	// 静态文件
	router.StaticFS("/public", http.Dir(pkg.RootPath()+"/public"))

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由分组
	routeGroup := router.Group("")
	authMiddleware := []gin.HandlerFunc{jwtMiddleware, permissionMiddleware}

	// 全局限流:rateLimitMiddleware.Handle() 用户限流:rateLimitMiddleware.UserRateLimit(1, 1) ip限流:rateLimitMiddleware.IpRateLimit(1, 1)
	// 健康检查
	routeGroup.GET("/ping", rateLimitMiddleware.IpRateLimit(1, 1), func(c *gin.Context) {
		facade.Response().Success(c, errcode.NewError(0, "pong"))
	})

	// 业务路由
	route.Mount(routeGroup, authMiddleware, All()...)
}
