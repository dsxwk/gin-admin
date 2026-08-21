package facade

import (
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
	return RequestFacade{}
}

type RequestFacade struct{}

// Path 获取请求路径参数
// 使用示例:
//
//	name := facade.Request().Path[string](ctx, "name", "default")
func (r RequestFacade) Path[T any](ctx *gin.Context, key string, defaultValue T) T {
	return request.NewClient().Path[T](ctx, key, defaultValue)
}

// GetHeader 获取请求头
func (r RequestFacade) GetHeader[T any](ctx *gin.Context, key string, defaultValue T) T {
	return request.NewClient().GetHeader[T](ctx, key, defaultValue)
}

// Header 设置请求头
func (r RequestFacade) Header[T any](ctx *gin.Context, key string, value T) {
	request.NewClient().Header[T](ctx, key, value)
}

// Bind 绑定请求参数
func (r RequestFacade) Bind(ctx *gin.Context, v any) error {
	return request.NewClient().Bind(ctx, v)
}

// Validate 验证请求数据
func (r RequestFacade) Validate(data interface{}, scene string) error {
	return request.NewClient().Validate(data, scene)
}

// BindValidate 绑定参数并验证
func (r RequestFacade) BindValidate(ctx *gin.Context, v any, scene string) error {
	return request.NewClient().BindValidate(ctx, v, scene)
}

// ValidateWithMessages 验证并自定义错误消息
func (r RequestFacade) ValidateWithMessages(data interface{}, scene string, messages map[string]string) error {
	return request.NewClient().ValidateWithMessages(data, scene, messages)
}

// ValidateWithTranslates 验证并自定义字段翻译
func (r RequestFacade) ValidateWithTranslates(data interface{}, scene string, translates map[string]string) error {
	return request.NewClient().ValidateWithTranslates(data, scene, translates)
}

// GetValidator 获取验证器实例
func (r RequestFacade) GetValidator(data interface{}, scene string) *validate.Validation {
	return request.NewClient().GetValidator(data, scene)
}
