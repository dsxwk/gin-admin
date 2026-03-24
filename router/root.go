package router

import (
	"gin/app/middleware"
	"gin/common/errcode"
	"gin/common/response"
	_ "gin/docs"
	"gin/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

var (
	loggerMiddleware    = middleware.Logger{}.Handle()
	corsMiddleware      = middleware.Cors{}.Handle()
	jwtMiddleware       = middleware.Jwt{}.Handle()
	recoverMiddleware   = middleware.Recover{}.Handle()
	rateLimitMiddleware = middleware.RateLimit{}
)

// LoadRouters 加载路由
func LoadRouters(router *gin.Engine) {
	// 静态文件
	router.StaticFS("/public", http.Dir(pkg.GetRootPath()+"/public"))

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由分组
	public := router.Group("", loggerMiddleware, corsMiddleware, recoverMiddleware) // 无需权限
	auth := public.Group("", jwtMiddleware)                                         // 需要权限

	// 健康检查
	// 全局限流:rateLimitMiddleware.Handle() 用户限流:rateLimitMiddleware.UserRateLimit(1, 1) ip限流:rateLimitMiddleware.IpRateLimit(1, 1)
	public.GET("/ping", rateLimitMiddleware.IpRateLimit(1, 1), func(c *gin.Context) {
		err := errcode.NewError(0, "pong")
		response.Success(c, &err)
	})

	// 自动注册
	AutoLoads(public, auth)
}
