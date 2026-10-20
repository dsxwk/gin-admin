package router

import (
	"gin/app/controller/v1"

	"github.com/gin-gonic/gin"
)

// FileRouter 文件路由
type FileRouter struct{}

// Register 注册路由
func (r *FileRouter) Register(routerGroup *gin.RouterGroup) {
	var file v1.FileController

	router := routerGroup.Group("/api/v1/file")
	{
		router.GET("", file.List)
		router.POST("/upload", file.Upload)
		router.DELETE("/:id", file.Delete)
		router.PUT("/:id/group", file.UpdateGroup)
	}
}

// IsAuth 是否需要鉴权
func (r *FileRouter) IsAuth() bool {
	return true
}
