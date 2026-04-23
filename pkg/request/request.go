package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"strconv"
	"strings"
)

type Utils struct{}

// Path 泛型获取请求路径参数
func Path[T any](ctx *gin.Context, key string, defaultValue T) T {
	val := ctx.Param(key)
	if val == "" {
		return defaultValue
	}

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

// Bind 绑定请求参数
func (r Utils) Bind(ctx *gin.Context, v any) error {
	ct := ctx.ContentType()

	switch {
	case strings.HasPrefix(ct, "application/json"):
		return ctx.ShouldBindJSON(v)

	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"),
		strings.HasPrefix(ct, "multipart/form-data"):
		return ctx.ShouldBind(v)

	default:
		return ctx.ShouldBind(v)
	}
}

// BindValidate 绑定参数并验证
func (r Utils) BindValidate(ctx *gin.Context, v any, scene string) error {
	if err := r.Bind(ctx, v); err != nil {
		return err
	}
	return r.Validate(v, scene)
}

// ValidateWithMessages 验证并自定义错误消息
func (r Utils) ValidateWithMessages(data interface{}, scene string, messages map[string]string) error {
	v := validate.Struct(data, scene)
	v.WithMessages(messages)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ValidateWithTranslates 验证并自定义字段翻译
func (r Utils) ValidateWithTranslates(data interface{}, scene string, translates map[string]string) error {
	v := validate.Struct(data, scene)
	v.WithTranslates(translates)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// GetValidator 获取验证器实例
func (r Utils) GetValidator(data interface{}, scene string) *validate.Validation {
	return validate.Struct(data, scene)
}

// Validate 通用验证函数
func (r Utils) Validate(data interface{}, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}
