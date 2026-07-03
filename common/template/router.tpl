package {{.Package}}

import (
    {{- if ne .Package "router" }}
    "gin/router"
    {{- end }}
    "gin/app/controller/v1"
	"github.com/gin-gonic/gin"
)

// {{.Name}}Router {{.Description}}
type {{.Name}}Router struct {}

func init() {
	{{- if eq .Package "router" }}
	Register(&{{.Name}}Router{})
	{{- else }}
	router.Register(&{{.Name}}Router{})
	{{- end }}
}

// RegisterRoutes 注册路由
func (r *{{.Name}}Router) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var (
   	    {{.NameLower}} v1.{{.Name}}Controller
    )

    router := routerGroup.Group("/api/v1/{{.NameLower}}")
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
