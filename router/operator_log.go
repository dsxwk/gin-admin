package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// OperatorLogRouter 操作日志
type OperatorLogRouter struct{}

func init() {
	Register(&OperatorLogRouter{})
}

// RegisterRoutes 注册路由
func (r *OperatorLogRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		operatorLog v1.OperatorLogController
	)

	router := routerGroup.Group("/api/v1/operator-log")
	{
		// 列表
		router.GET("", operatorLog.List)
		// 详情
		router.GET("/:id", operatorLog.Detail)
		// 删除
		router.DELETE("/:id", operatorLog.Delete)
		// 批量删除
		router.POST("/batch-delete", operatorLog.BatchDelete)
	}
}

// IsAuth 是否需要鉴权
func (r *OperatorLogRouter) IsAuth() bool {
	return true
}
