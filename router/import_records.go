package router

import (
	"gin/app/controller/v1"

	"github.com/gin-gonic/gin"
)

// ImportRecordsRouter ImportRecords
type ImportRecordsRouter struct{}

// Register 注册路由
func (r *ImportRecordsRouter) Register(routerGroup *gin.RouterGroup) {
	var (
		importRecords v1.ImportRecordsController
	)

	router := routerGroup.Group("/api/v1/import-records")
	{
		// 列表
		router.GET("", importRecords.List)
		// 删除
		router.DELETE("/:id", importRecords.Delete)
	}
}

// IsAuth 是否需要鉴权
func (r *ImportRecordsRouter) IsAuth() bool {
	return true
}
