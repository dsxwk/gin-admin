package router

import (
	"gin/app/facade"
	"gin/app/middleware"
	_ "gin/docs"
	"gin/pkg/errcode"
	"gin/pkg/pathx"
	"gin/pkg/route"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouters 加载路由并返回鉴权路由Key
func NewRouters(router *gin.Engine) []string {
	cfg := facade.Config()
	if cfg == nil {
		return nil
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
	router.StaticFS("/public", http.Dir(pathx.RootPath()+"/public"))
	registerLocalStorageStatic(router, cfg.Filesystem.Local.Root)

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
	return route.Mount(router, routeGroup, authMiddleware, All()...)
}

// registerLocalStorageStatic 注册本地文件静态访问
func registerLocalStorageStatic(router *gin.Engine, localRoot string) {
	publicPath, staticDir := localStorageStatic(localRoot)
	if router == nil || publicPath == "/" || staticDir == "" {
		return
	}
	router.StaticFS(publicPath, gin.Dir(staticDir, false))
}

// localStorageStatic 获取本地文件访问前缀与目录
func localStorageStatic(localRoot string) (string, string) {
	localRoot = strings.TrimSpace(localRoot)
	if localRoot == "" {
		localRoot = filepath.Join("storage", "uploads")
	}
	localRoot = filepath.Clean(localRoot)

	publicPath := "/" + strings.Trim(path.Clean(strings.ReplaceAll(localRoot, "\\", "/")), "/")
	staticDir := localRoot
	if !filepath.IsAbs(staticDir) {
		staticDir = filepath.Join(pathx.RootPath(), staticDir)
	}
	return publicPath, filepath.Clean(staticDir)
}
