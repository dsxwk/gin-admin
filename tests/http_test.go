package tests

import (
	"context"
	"encoding/json"
	"gin/app/facade"
	"gin/common/ctxkey"
	"gin/common/errcode"
	h "gin/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// setupTestServer 创建测试服务器
func setupTestServer() *httptest.Server {
	r := gin.Default()

	// GET 测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "pong", "data": nil})
	})

	// GET 带参数测试
	r.GET("/echo", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{"name": name, "age": age},
		})
	})

	// POST JSON测试
	r.POST("/echo", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": body,
		})
	})

	// POST 表单测试
	r.POST("/form", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": gin.H{"name": name, "email": email},
		})
	})

	// 错误响应测试
	r.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "internal server error",
		})
	})

	// 延迟响应测试
	r.GET("/delay", func(c *gin.Context) {
		delay := c.Query("delay")
		if delay != "" {
			d, _ := time.ParseDuration(delay)
			time.Sleep(d)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
	})

	return httptest.NewServer(r)
}

// TestHttpRequest 测试基本HTTP请求
func TestHttpRequest(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-trace-id")

	// 测试GET请求
	resp, err := facade.Http.Request(ctx, "GET", ts.URL+"/ping", nil)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp), &result)
	require.NoError(t, err)
	require.Equal(t, float64(0), result["code"])
	require.Equal(t, "pong", result["msg"])
}

// TestHttpRequestJson 测试JSON响应解析
func TestHttpRequestJson(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-trace-id")

	// 测试GET请求并解析JSON
	resp, err := facade.RequestJson[errcode.SuccessResponse](
		ctx,
		"GET",
		ts.URL+"/ping",
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "pong", resp.Msg)
}

// TestHttpRequestWithQuery 测试带Query参数的请求
func TestHttpRequestWithQuery(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	opt := &h.Option{
		Query: map[string]string{
			"name": "张三",
			"age":  "18",
		},
	}

	type EchoResponse struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}

	resp, err := facade.RequestJson[EchoResponse](
		ctx,
		"GET",
		ts.URL+"/echo",
		opt,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "张三", resp.Data["name"])
	require.Equal(t, "18", resp.Data["age"])
}

// TestHttpRequestWithBody 测试带Body的POST请求
func TestHttpRequestWithBody(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	bodyData := map[string]interface{}{
		"name":  "李四",
		"email": "lisi@example.com",
	}
	body, _ := json.Marshal(bodyData)

	opt := &h.Option{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	}

	type EchoResponse struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}

	resp, err := facade.RequestJson[EchoResponse](
		ctx,
		"POST",
		ts.URL+"/echo",
		opt,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "李四", resp.Data["name"])
	require.Equal(t, "lisi@example.com", resp.Data["email"])
}

// TestHttpRequestWithForm 测试表单请求
func TestHttpRequestWithForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	opt := &h.Option{
		Form: map[string]string{
			"name":  "王五",
			"email": "wangwu@example.com",
		},
	}

	type FormResponse struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}

	resp, err := facade.RequestJson[FormResponse](
		ctx,
		"POST",
		ts.URL+"/form",
		opt,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "王五", resp.Data["name"])
	require.Equal(t, "wangwu@example.com", resp.Data["email"])
}

// TestHttpRequestWithHeaders 测试带自定义Headers的请求
func TestHttpRequestWithHeaders(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	opt := &h.Option{
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"Authorization":   "Bearer token123",
		},
	}

	// 创建一个可以读取 header 的测试服务器
	r := gin.Default()
	r.GET("/headers", func(c *gin.Context) {
		customHeader := c.GetHeader("X-Custom-Header")
		authHeader := c.GetHeader("Authorization")
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"x-custom-header": customHeader,
				"authorization":   authHeader,
			},
		})
	})
	ts2 := httptest.NewServer(r)
	defer ts2.Close()

	type HeaderResponse struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}

	resp, err := facade.RequestJson[HeaderResponse](
		ctx,
		"GET",
		ts2.URL+"/headers",
		opt,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "custom-value", resp.Data["x-custom-header"])
	require.Equal(t, "Bearer token123", resp.Data["authorization"])
}

// TestHttpRequestWithTimeout 测试超时
func TestHttpRequestWithTimeout(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	// 设置100ms超时,但服务器会延迟200ms
	opt := &h.Option{
		Timeout: 100 * time.Millisecond,
	}

	_, err := facade.Http.Request(ctx, "GET", ts.URL+"/delay?delay=200ms", opt)
	require.Error(t, err)
	require.Contains(t, err.Error(), "请求失败")
}

// TestHttpRequestPostForm 测试POST表单请求(使用Request方法)
func TestHttpRequestPostForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	opt := &h.Option{
		Form: map[string]string{
			"name":  "赵六",
			"email": "zhaoliu@example.com",
		},
	}

	resp, err := facade.Http.Request(ctx, "POST", ts.URL+"/form", opt)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp), &result)
	require.NoError(t, err)

	data := result["data"].(map[string]interface{})
	require.Equal(t, "赵六", data["name"])
	require.Equal(t, "zhaoliu@example.com", data["email"])
}

// TestHttpRequestErrorResponse 测试错误响应
func TestHttpRequestErrorResponse(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	_, err := facade.Http.Request(ctx, "GET", ts.URL+"/error", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "状态码: 500")
}

// TestHttpRequestJsonErrorResponse 测试JSON解析错误响应
func TestHttpRequestJsonErrorResponse(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	type ErrorResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	resp, err := facade.RequestJson[ErrorResponse](
		ctx,
		"GET",
		ts.URL+"/error",
		nil,
	)
	require.Error(t, err)
	require.Nil(t, resp)
}

// TestHttpRequestInvalidURL 测试无效URL
func TestHttpRequestInvalidURL(t *testing.T) {
	ctx := context.Background()

	_, err := facade.Http.Request(ctx, "GET", "http://invalid.url.that.does.not.exist", nil)
	require.Error(t, err)
}
