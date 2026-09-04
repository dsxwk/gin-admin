package router

import (
	"context"
	"gin/app/facade"
	"gin/app/middleware"
	"gin/app/service"
	"gin/common/errcode"
	"gin/common/response"
	_ "gin/docs"
	"gin/pkg"
	"gin/pkg/serviceprovider/mcp"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	timeoutMiddleware     = middleware.Timeout{}.Handle(facade.Config().App.Timeout)
	loggerMiddleware      = middleware.Logger{}.Handle()
	corsMiddleware        = middleware.Cors{}.Handle()
	jwtMiddleware         = middleware.Jwt{}.Handle()
	recoverMiddleware     = middleware.Recover{}.Handle()
	rateLimitMiddleware   = middleware.RateLimit{}
	permissionMiddleware  = middleware.Permission{}.Handle()
	operatorLogMiddleware = middleware.OperatorLog{}.Handle()
)

// LoadRouters 加载路由
func LoadRouters(router *gin.Engine) {
	// 全局中间件
	router.Use(corsMiddleware, timeoutMiddleware, loggerMiddleware, recoverMiddleware, operatorLogMiddleware)

	// 静态文件
	router.StaticFS("/public", http.Dir(pkg.GetRootPath()+"/public"))

	// Swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由分组
	public := router.Group("")                                    // 无需权限
	auth := router.Group("", jwtMiddleware, permissionMiddleware) // 需要权限

	// 全局限流:rateLimitMiddleware.Handle() 用户限流:rateLimitMiddleware.UserRateLimit(1, 1) ip限流:rateLimitMiddleware.IpRateLimit(1, 1)
	// 健康检查
	public.GET("/ping", rateLimitMiddleware.IpRateLimit(1, 1), func(c *gin.Context) {
		response.Response{}.Success(c, errcode.NewError(0, "pong"))
	})

	// 自动注册
	AutoLoads(public, auth)

	// MCP服务路由(MCP自带认证,使用public分组)
	if cfg := facade.Config(); cfg != nil && cfg.Mcp.Enabled {
		mcpHandler := facade.Get[*mcp.Handler]("mcp")
		if mcpHandler != nil {
			mcpPath := cfg.Mcp.Path
			if mcpPath == "" {
				mcpPath = "/mcp"
			}
			router.POST(mcpPath, gin.WrapH(mcpHandler))
		}
	}
}

// SyncPermissionRoutes 同步路由权限到数据库(仅服务器启动时调用)
func SyncPermissionRoutes() {
	permissionKeys := GenerateAuthPermissionKeys()
	if len(permissionKeys) > 0 {
		svc := service.PermissionService{}
		_ = svc.SyncRoutePermissions(context.Background(), permissionKeys)
	}
}
