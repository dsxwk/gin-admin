package facade

import (
	"gin/pkg/request"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

// Request 泛型函数,返回对应类型的门面实例
// 使用示例:
//
//	name := facade.Request[string]().Path(ctx, "name", "default")
//	age := facade.Request[int]().Path(ctx, "age", 18)
//	userID := facade.Request[int64]().Path(ctx, "id", 0)
func Request[T any]() RequestFacade[T] {
	return RequestFacade[T]{}
}

type RequestFacade[T any] struct{}

// Path 获取请求路径参数
// 使用示例:
//
//	name := facade.Request[string]().Path(ctx, "name", "default")
func (r RequestFacade[T]) Path(ctx *gin.Context, key string, defaultValue T) T {
	return request.Utils[T]{}.Path(ctx, key, defaultValue)
}

// Bind 绑定请求参数
func (r RequestFacade[T]) Bind(ctx *gin.Context, v any) error {
	return request.Utils[T]{}.Bind(ctx, v)
}

// Validate 验证请求数据
func (r RequestFacade[T]) Validate(data interface{}, scene string) error {
	return request.Utils[T]{}.Validate(data, scene)
}

// BindValidate 绑定参数并验证
func (r RequestFacade[T]) BindValidate(ctx *gin.Context, v any, scene string) error {
	return request.Utils[T]{}.BindValidate(ctx, v, scene)
}

// ValidateWithMessages 验证并自定义错误消息
func (r RequestFacade[T]) ValidateWithMessages(data interface{}, scene string, messages map[string]string) error {
	return request.Utils[T]{}.ValidateWithMessages(data, scene, messages)
}

// ValidateWithTranslates 验证并自定义字段翻译
func (r RequestFacade[T]) ValidateWithTranslates(data interface{}, scene string, translates map[string]string) error {
	return request.Utils[T]{}.ValidateWithTranslates(data, scene, translates)
}

// GetValidator 获取验证器实例
func (r RequestFacade[T]) GetValidator(data interface{}, scene string) *validate.Validation {
	return request.Utils[T]{}.GetValidator(data, scene)
}
