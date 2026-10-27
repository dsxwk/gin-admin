package facade

import (
	"gin/pkg/serviceprovider/request"

	"github.com/gin-gonic/gin"
)

// Request 门面函数,返回门面实例
// 使用示例:
//
//	name := facade.Request().Path[string](ctx, "name", "default")
//	age := facade.Request().Path[int](ctx, "age", 18)
//	userID := facade.Request().Path[int64](ctx, "id", 0)
func Request() RequestFacade {
	return RequestFacade{
		client: App().Request(),
	}
}

// RequestFacade 请求门面
type RequestFacade struct {
	client *request.Client
}

// Query 获取请求查询参数
func (r RequestFacade) Query[T any](ctx *gin.Context, key string, defaultValue T) T {
	if r.client == nil {
		return defaultValue
	}
	return r.client.Query(ctx, key, defaultValue)
}

// Form 获取请求表单参数
func (r RequestFacade) Form[T any](ctx *gin.Context, key string, defaultValue T) T {
	if r.client == nil {
		return defaultValue
	}
	return r.client.Form(ctx, key, defaultValue)
}

// Path 获取请求路径参数
// 使用示例:
//
//	name := facade.Request().Path[string](ctx, "name", "default")
func (r RequestFacade) Path[T any](ctx *gin.Context, key string, defaultValue T) T {
	if r.client == nil {
		return defaultValue
	}
	return r.client.Path(ctx, key, defaultValue)
}

// Header 获取请求头
func (r RequestFacade) Header[T any](ctx *gin.Context, key string, defaultValue T) T {
	if r.client == nil {
		return defaultValue
	}
	return r.client.Header(ctx, key, defaultValue)
}

// SetHeader 设置请求头
func (r RequestFacade) SetHeader[T any](ctx *gin.Context, key string, value T) {
	if r.client == nil {
		return
	}
	r.client.SetHeader(ctx, key, value)
}

// Params 获取请求参数
func (r RequestFacade) Params(ctx *gin.Context) any {
	if r.client == nil {
		return map[string]any{}
	}
	return r.client.Params(ctx)
}

// Bind 绑定请求参数
func (r RequestFacade) Bind(ctx *gin.Context, v any) error {
	if r.client == nil {
		return nil
	}
	return r.client.Bind(ctx, v)
}

// Validate 验证请求数据
func (r RequestFacade) Validate(ctx *gin.Context, data any, scene string) error {
	if r.client == nil {
		return nil
	}
	return r.client.Validate(ctx.Request.Context(), data, scene)
}

// BindValidate 绑定参数并验证
func (r RequestFacade) BindValidate(ctx *gin.Context, v any, scene string) error {
	if r.client == nil {
		return nil
	}
	return r.client.BindValidate(ctx, v, scene)
}
