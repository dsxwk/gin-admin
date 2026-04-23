package facade

import (
	"gin/pkg/request"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

// Request 请求验证门面
// 使用示例:
//
//	type UserRequest struct {
//	    request.Utils
//	    Name string `json:"name" validate:"required"`
//	}
//
//	req := UserRequest{Name: "test"}
//	err := facade.Request.Validate(req, "create") # err := facade.Request.BindValidate(ctx, &req, "create")
var Request = &requestFacade{}

type requestFacade struct{}

// RequestPath 获取请求参数（泛型方法，返回指定类型）
// 使用示例:
//
//	name := facade.RequestPath[string](ctx, "name", "default")
//	age := facade.RequestPath[int](ctx, "age", 18)
//	userID := facade.RequestPath[int64](ctx, "userId", 0)
func RequestPath[T any](ctx *gin.Context, key string, defaultValue T) T {
	return request.Path[T](ctx, key, defaultValue)
}

// Bind 绑定请求参数
func Bind(ctx *gin.Context, v any) error {
	return request.Utils{}.Bind(ctx, v)
}

// Validate 验证请求数据
func (r *requestFacade) Validate(data interface{}, scene string) error {
	return request.Utils{}.Validate(data, scene)
}

// BindValidate 绑定参数并验证
func (r *requestFacade) BindValidate(ctx *gin.Context, v any, scene string) error {
	return request.Utils{}.BindValidate(ctx, v, scene)
}

// ValidateWithMessages 验证并自定义错误消息
func (r *requestFacade) ValidateWithMessages(data interface{}, scene string, messages map[string]string) error {
	return request.Utils{}.ValidateWithMessages(data, scene, messages)
}

// ValidateWithTranslates 验证并自定义字段翻译
func (r *requestFacade) ValidateWithTranslates(data interface{}, scene string, translates map[string]string) error {
	return request.Utils{}.ValidateWithTranslates(data, scene, translates)
}

// GetValidator 获取验证器实例
func (r *requestFacade) GetValidator(data interface{}, scene string) *validate.Validation {
	return request.Utils{}.GetValidator(data, scene)
}
