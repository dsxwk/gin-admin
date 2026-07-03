package router

import (
	"gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// ArticleRouter Article
type ArticleRouter struct{}

func init() {
	Register(&ArticleRouter{})
}

// RegisterRoutes 注册路由
func (r *ArticleRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
		article v1.ArticleController
	)

	router := routerGroup.Group("/api/v1/article")
	{
		// 列表
		router.GET("", article.List)
		// 创建
		router.POST("", article.Create)
		// 更新
		router.PUT("/:id", article.Update)
		// 删除
		router.DELETE("/:id", article.Delete)
		// 详情
		router.GET("/:id", article.Detail)
	}
}

// IsAuth 是否需要鉴权
func (r *ArticleRouter) IsAuth() bool {
	return true
}
