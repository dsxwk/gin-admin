package router

import (
	"gin/app/controller/v1"

	"github.com/gin-gonic/gin"
)

// FileGroupRouter 文件分组路由
type FileGroupRouter struct{}

// Register 注册路由
func (r *FileGroupRouter) Register(routerGroup *gin.RouterGroup) {
	var fileGroup v1.FileGroupController

	router := routerGroup.Group("/api/v1/file-group")
	{
		router.GET("", fileGroup.List)
		router.POST("", fileGroup.Create)
		router.PUT("/:id", fileGroup.Update)
		router.DELETE("/:id", fileGroup.Delete)
	}
}

// IsAuth 是否需要鉴权
func (r *FileGroupRouter) IsAuth() bool {
	return true
}
