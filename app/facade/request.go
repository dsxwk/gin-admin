package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider/request"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

// Request 门面函数,返回门面实例
// 使用示例:
//
//	name := facade.Request().Path[string](ctx, "name", "default")
//	age := facade.Request().Path[int](ctx, "age", 18)
//	userID := facade.Request().Path[int64](ctx, "id", 0)
func Request() RequestFacade {
	return RequestFacade{
		client: container.Default().Request(),
	}
}

// RequestFacade 请求门面
type RequestFacade struct {
	client *request.Client
}

// Path 获取请求路径参数
// 使用示例:
//
//	name := facade.Request().Path[string](ctx, "name", "default")
func (r RequestFacade) Path[T any](ctx *gin.Context, key string, defaultValue T) T {
	if r.client == nil {
		return defaultValue
	}
	return r.client.Path[T](ctx, key, defaultValue)
}

// GetHeader 获取请求头
func (r RequestFacade) GetHeader[T any](ctx *gin.Context, key string, defaultValue T) T {
	if r.client == nil {
		return defaultValue
	}
	return r.client.GetHeader[T](ctx, key, defaultValue)
}

// Header 设置请求头
func (r RequestFacade) Header[T any](ctx *gin.Context, key string, value T) {
	if r.client == nil {
		return
	}
	r.client.Header[T](ctx, key, value)
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

// ValidateWithMessages 验证并自定义错误消息
func (r RequestFacade) ValidateWithMessages(data any, scene string, messages map[string]string) error {
	if r.client == nil {
		return nil
	}
	return r.client.ValidateWithMessages(data, scene, messages)
}

// ValidateWithTranslates 验证并自定义字段翻译
func (r RequestFacade) ValidateWithTranslates(data any, scene string, translates map[string]string) error {
	if r.client == nil {
		return nil
	}
	return r.client.ValidateWithTranslates(data, scene, translates)
}

// Validator 获取验证器实例
func (r RequestFacade) Validator(data any, scene string) *validate.Validation {
	if r.client == nil {
		return nil
	}
	return r.client.Validator(data, scene)
}
