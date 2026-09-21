package {{.Package}}

import (
    "gin/app/controller/v1"
    "gin/pkg/route"
	"github.com/gin-gonic/gin"
)

// {{.Name}}Router {{.Description}}
type {{.Name}}Router struct {}

func init() {
	route.Register(&{{.Name}}Router{})
}

// RegisterRoutes 注册路由
func (r *{{.Name}}Router) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
   	    {{.NameLower}} v1.{{.Name}}Controller
    )

    router := routerGroup.Group("/api/v1/{{.NameKebabCase}}")
    {
        // 列表
        router.GET("", {{.NameLower}}.List)
        // 创建
        router.POST("", {{.NameLower}}.Create)
        // 更新
        router.PUT("/:id", {{.NameLower}}.Update)
        // 删除
        router.DELETE("/:id", {{.NameLower}}.Delete)
        // 详情
        router.GET("/:id", {{.NameLower}}.Detail)
    }
}

// IsAuth 是否需要鉴权
func (r *{{.Name}}Router) IsAuth() bool {
	return true
}
