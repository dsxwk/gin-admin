package request

import (
	"context"
	"fmt"
	"gin/pkg/errcode"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct{}

// Query 泛型获取请求查询参数
func (c Client) Query[T any](ctx *gin.Context, key string, defaultValue T) T {
	val := ctx.Query(key)
	if val == "" {
		return defaultValue
	}

	return getValue[T](val, defaultValue)
}

// getValue 获取参数值
func getValue[T any](val string, defaultValue T) T {
	switch any(defaultValue).(type) {
	case string:
		return any(val).(T)

	case int:
		if v, err := strconv.Atoi(val); err == nil {
			return any(v).(T)
		}
	case int64:
		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			return any(v).(T)
		}
	case int32:
		if v, err := strconv.ParseInt(val, 10, 32); err == nil {
			return any(int32(v)).(T)
		}
	case int16:
		if v, err := strconv.ParseInt(val, 10, 16); err == nil {
			return any(int16(v)).(T)
		}
	case int8:
		if v, err := strconv.ParseInt(val, 10, 8); err == nil {
			return any(int8(v)).(T)
		}
	case uint:
		if v, err := strconv.ParseUint(val, 10, 64); err == nil {
			return any(uint(v)).(T)
		}
	case uint64:
		if v, err := strconv.ParseUint(val, 10, 64); err == nil {
			return any(v).(T)
		}
	case uint32:
		if v, err := strconv.ParseUint(val, 10, 32); err == nil {
			return any(uint32(v)).(T)
		}
	case uint16:
		if v, err := strconv.ParseUint(val, 10, 16); err == nil {
			return any(uint16(v)).(T)
		}
	case uint8:
		if v, err := strconv.ParseUint(val, 10, 8); err == nil {
			return any(uint8(v)).(T)
		}
	case bool:
		if v, err := strconv.ParseBool(val); err == nil {
			return any(v).(T)
		}
	case float32:
		if v, err := strconv.ParseFloat(val, 32); err == nil {
			return any(float32(v)).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return any(v).(T)
		}
	}

	return defaultValue
}

// Path 泛型获取请求路径参数
func (c Client) Path[T any](ctx *gin.Context, key string, defaultValue T) T {
	val := ctx.Param(key)
	if val == "" {
		return defaultValue
	}

	return getValue[T](val, defaultValue)
}

// GetHeader 泛型获取请求头参数
func (c Client) GetHeader[T any](ctx *gin.Context, key string, defaultValue T) T {
	val := ctx.GetHeader(key)
	if val == "" {
		return defaultValue
	}

	return getValue[T](val, defaultValue)
}

// Header 泛型设置请求头参数
func (c Client) Header[T any](ctx *gin.Context, key string, value T) {
	ctx.Header(key, fmt.Sprintf("%v", value))
}

// Bind 绑定请求参数
func (c Client) Bind(ctx *gin.Context, v any) error {
	// 先绑定Query参数
	if err := ctx.ShouldBindQuery(v); err != nil {
		return fmt.Errorf("bind query error: %w", err)
	}

	// 如果有请求体再绑定Body
	if ctx.Request.ContentLength > 0 {
		switch {
		case strings.HasPrefix(ctx.ContentType(), "application/json"):
			if err := ctx.ShouldBindJSON(v); err != nil {
				return fmt.Errorf("bind json error: %w", err)
			}
		default:
			if err := ctx.ShouldBind(v); err != nil {
				return fmt.Errorf("bind form error: %w", err)
			}
		}
	}

	return nil
}

type requestContextSetter interface {
	SetContext(context.Context)
}

// setRequestContext 设置请求上下文
func setRequestContext(ctx context.Context, data any) {
	if setter, ok := data.(requestContextSetter); ok {
		setter.SetContext(ctx)
	}
}

// BindValidate 绑定参数并验证
func (c Client) BindValidate(ctx *gin.Context, v any, scene string) error {
	if err := c.Bind(ctx, v); err != nil {
		return errcode.ArgsError().WithMsg(err.Error())
	}
	setRequestContext(ctx.Request.Context(), v)

	// 场景为空不验证
	if scene == "" {
		return nil
	}
	return c.validate(v, scene)
}

// ValidateWithMessages 验证并自定义错误消息
func (c Client) ValidateWithMessages(data any, scene string, messages map[string]string) error {
	v := validate.Struct(data, scene)
	v.WithMessages(messages)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ValidateWithTranslates 验证并自定义字段翻译
func (c Client) ValidateWithTranslates(data any, scene string, translates map[string]string) error {
	v := validate.Struct(data, scene)
	v.WithTranslates(translates)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// GetValidator 获取验证器实例
func (c Client) GetValidator(data any, scene string) *validate.Validation {
	return validate.Struct(data, scene)
}

// Validate 通用验证函数,自动注入请求上下文
func (c Client) Validate(ctx context.Context, data any, scene string) error {
	setRequestContext(ctx, data)
	return c.validate(data, scene)
}

// validate 执行验证
func (c Client) validate(data any, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}
