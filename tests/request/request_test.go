package request

import (
	"gin/app/facade"
	"gin/pkg/errcode"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dsxwk/validatex"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// requestTestData 请求测试数据
type requestTestData struct {
	Name string `json:"name" form:"name"`
	Age  int    `json:"age" form:"age"`
}

// newRequestTestContext 创建请求测试上下文
func newRequestTestContext(method, target, body, contentType string) *gin.Context {
	ctx, _ := newRequestTestContextWithRecorder(method, target, body, contentType)
	return ctx
}

// newRequestTestContextWithRecorder 创建带响应记录器的请求测试上下文
func newRequestTestContextWithRecorder(method, target, body, contentType string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, target, strings.NewReader(body))
	if contentType != "" {
		ctx.Request.Header.Set("Content-Type", contentType)
	}
	return ctx, recorder
}

// requireValidationErrors 获取结构化校验错误
func requireValidationErrors(t *testing.T, err error) validatex.Errors {
	t.Helper()
	var errorCode errcode.ErrorCode
	require.ErrorAs(t, err, &errorCode)
	data, ok := errorCode.Data.(validatex.ValidationErrorData)
	require.True(t, ok)
	require.NotEmpty(t, data.Errors)
	require.Equal(t, data.Errors[0].Msg, errorCode.Msg)
	return data.Errors
}

// validationFieldPaths 提取校验字段路径
func validationFieldPaths(errors validatex.Errors) []string {
	paths := make([]string, 0, len(errors))
	for _, item := range errors {
		paths = append(paths, item.Path)
	}
	return paths
}

// validationErrorRules 提取校验错误规则
func validationErrorRules(errors validatex.Errors) []string {
	rules := make([]string, 0, len(errors))
	for _, item := range errors {
		rules = append(rules, item.Rule)
	}
	return rules
}

// TestRequestForm 表单参数获取测试
func TestRequestForm(t *testing.T) {
	ctx := newRequestTestContext(
		"POST",
		"/",
		"name=alice&age=20&enabled=true",
		"application/x-www-form-urlencoded",
	)

	require.Equal(t, "alice", facade.Request().Form[string](ctx, "name", ""))
	require.Equal(t, 20, facade.Request().Form[int](ctx, "age", 0))
	require.True(t, facade.Request().Form[bool](ctx, "enabled", false))
	require.Equal(t, "default", facade.Request().Form[string](ctx, "missing", "default"))
}

// TestRequestHeader 请求头获取和设置测试
func TestRequestHeader(t *testing.T) {
	ctx := newRequestTestContext("GET", "/", "", "")
	ctx.Request.Header.Set("X-Trace-Id", "trace-123")

	require.Equal(t, "trace-123", facade.Request().Header[string](ctx, "X-Trace-Id", ""))

	facade.Request().SetHeader(ctx, "X-Test", "value")
	require.Equal(t, "value", ctx.GetHeader("X-Test"))
}

// TestRequestParams 请求参数提取测试
func TestRequestParams(t *testing.T) {
	ctx := newRequestTestContext("POST", "/?page=1", `{"name":"body"}`, "application/json")

	params := facade.Request().Params(ctx)
	require.Equal(t, map[string]any{"name": "body"}, params)

	var data requestTestData
	require.NoError(t, facade.Request().Bind(ctx, &data))
	require.Equal(t, "body", data.Name)
}

// TestRequestBindJSONWithQuery JSON和Query合并绑定测试
func TestRequestBindJSONWithQuery(t *testing.T) {
	ctx := newRequestTestContext("POST", "/?age=18", `{"name":"body"}`, "application/json")
	var data requestTestData

	require.NoError(t, facade.Request().Bind(ctx, &data))
	require.Equal(t, "body", data.Name)
	require.Equal(t, 18, data.Age)
}

// TestRequestBindChunkedJSON 分块JSON请求绑定测试
func TestRequestBindChunkedJSON(t *testing.T) {
	ctx := newRequestTestContext("POST", "/", `{"name":"chunked"}`, "application/json")
	ctx.Request.ContentLength = -1
	var data requestTestData

	require.NoError(t, facade.Request().Bind(ctx, &data))
	require.Equal(t, "chunked", data.Name)
}

// TestRequestBindEmptyJSONBody 空JSON请求体绑定测试
func TestRequestBindEmptyJSONBody(t *testing.T) {
	ctx := newRequestTestContext("POST", "/?name=query", "", "application/json")
	var data requestTestData

	require.NoError(t, facade.Request().Bind(ctx, &data))
	require.Equal(t, "query", data.Name)
}

// TestRequestBindFormWithQuery 表单和Query合并绑定测试
func TestRequestBindFormWithQuery(t *testing.T) {
	ctx := newRequestTestContext(
		"POST",
		"/?name=query",
		"name=form&age=20",
		"application/x-www-form-urlencoded",
	)
	var data requestTestData

	require.NoError(t, facade.Request().Bind(ctx, &data))
	require.Equal(t, "form", data.Name)
	require.Equal(t, 20, data.Age)
}
