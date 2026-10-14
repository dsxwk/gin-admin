package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin/pkg/errcode"
	"io"
	"net/http"
	"net/url"
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

// Form 泛型获取表单参数
func (c Client) Form[T any](ctx *gin.Context, key string, defaultValue T) T {
	val := ctx.PostForm(key)
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

// Header 泛型获取请求头参数
func (c Client) Header[T any](ctx *gin.Context, key string, defaultValue T) T {
	val := ctx.GetHeader(key)
	if val == "" {
		return defaultValue
	}

	return getValue[T](val, defaultValue)
}

// SetHeader 泛型设置请求头参数
func (c Client) SetHeader[T any](ctx *gin.Context, key string, value T) {
	if ctx == nil || ctx.Request == nil {
		return
	}
	ctx.Request.Header.Set(key, fmt.Sprintf("%v", value))
}

// Params 获取请求参数用于日志记录
func (c Client) Params(ctx *gin.Context) any {
	if ctx == nil || ctx.Request == nil {
		return map[string]any{}
	}

	request := ctx.Request
	if isQueryMethod(request.Method) {
		return queryToMap(request.URL.Query())
	}
	if !hasRequestBody(request) {
		return map[string]any{}
	}

	body, err := io.ReadAll(request.Body)
	if err != nil || len(body) == 0 {
		return map[string]any{}
	}
	request.Body = io.NopCloser(bytes.NewReader(body))

	contentType := strings.ToLower(ctx.ContentType())
	switch {
	case isJSONContentType(contentType):
		var result map[string]any
		if json.Unmarshal(body, &result) == nil {
			return result
		}
	case isFormContentType(contentType):
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if request.ParseMultipartForm(32<<20) == nil {
				return queryToMap(request.PostForm)
			}
		} else if values, parseErr := url.ParseQuery(string(body)); parseErr == nil {
			return queryToMap(values)
		}
	}

	return string(body)
}

// queryToMap URL参数转map
func queryToMap(values url.Values) map[string]string {
	result := make(map[string]string, len(values))
	for key, value := range values {
		if len(value) > 0 {
			result[key] = value[0]
		}
	}
	return result
}

// Bind 绑定请求参数
// 绑定策略: JSON/XML先绑定Query再绑定Body,Body同名字段优先,表单交给Gin合并,空Body允许继续验证
func (c Client) Bind(ctx *gin.Context, v any) error {
	if ctx == nil || ctx.Request == nil || v == nil {
		return nil
	}

	if !hasRequestBody(ctx.Request) || isQueryMethod(ctx.Request.Method) {
		return bindQuery(ctx, v)
	}

	contentType := strings.ToLower(ctx.ContentType())
	switch {
	case isJSONContentType(contentType):
		if err := bindQuery(ctx, v); err != nil {
			return err
		}
		return bindBody(func() error { return ctx.ShouldBindJSON(v) }, "bind json error")
	case isXMLContentType(contentType):
		if err := bindQuery(ctx, v); err != nil {
			return err
		}
		return bindBody(func() error { return ctx.ShouldBindXML(v) }, "bind xml error")
	case isFormContentType(contentType):
		return bindBody(func() error { return ctx.ShouldBind(v) }, "bind form error")
	default:
		return bindQuery(ctx, v)
	}
}

// bindQuery 绑定Query参数
func bindQuery(ctx *gin.Context, v any) error {
	if err := ctx.ShouldBindQuery(v); err != nil {
		return fmt.Errorf("bind query error: %w", err)
	}
	return nil
}

// bindBody 绑定请求体,空Body不视为错误
func bindBody(bind func() error, message string) error {
	if err := bind(); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("%s: %w", message, err)
	}
	return nil
}

// hasRequestBody 判断是否存在请求体
func hasRequestBody(request *http.Request) bool {
	if request == nil || request.Body == nil || request.Body == http.NoBody {
		return false
	}
	return request.ContentLength != 0
}

// isQueryMethod 判断是否只允许Query绑定
func isQueryMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

// isJSONContentType 判断JSON请求类型
func isJSONContentType(contentType string) bool {
	mediaType := strings.TrimSpace(strings.Split(contentType, ";")[0])
	return mediaType == "application/json" || strings.HasSuffix(mediaType, "+json")
}

// isXMLContentType 判断XML请求类型
func isXMLContentType(contentType string) bool {
	mediaType := strings.TrimSpace(strings.Split(contentType, ";")[0])
	return mediaType == "application/xml" || mediaType == "text/xml" || strings.HasSuffix(mediaType, "+xml")
}

// isFormContentType 判断表单请求类型
func isFormContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "application/x-www-form-urlencoded") ||
		strings.HasPrefix(contentType, "multipart/form-data")
}

// BindValidate 绑定参数并验证
func (c Client) BindValidate(ctx *gin.Context, v any, scene string) error {
	if err := c.Bind(ctx, v); err != nil {
		return errcode.ArgsError().WithMsg(err.Error())
	}
	return c.withContext(ctx.Request.Context(), v, scene)
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

// Validator 获取验证器实例
func (c Client) Validator(data any, scene string) *validate.Validation {
	return validate.Struct(data, scene)
}

// Validate 通用验证函数,自动注入请求上下文
func (c Client) Validate(ctx context.Context, data any, scene string) error {
	return c.withContext(ctx, data, scene)
}

// WithContext 请求上下文设置接口
type contextSetter interface {
	SetContext(context.Context)
}

// withContext 设置上下文并执行验证
func (c Client) withContext(ctx context.Context, data any, scene string) error {
	if setter, ok := data.(contextSetter); ok {
		setter.SetContext(ctx)
	}

	// 场景为空不验证
	if scene == "" {
		return nil
	}

	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}
