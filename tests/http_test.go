package tests

import (
	"context"
	"encoding/json"
	"gin/app/errcode"
	"gin/app/facade"
	"gin/common/ctxkey"
	h "gin/pkg/serviceprovider/http"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// setupTestServer 创建测试服务器
func setupTestServer() *httptest.Server {
	r := gin.Default()

	// GET 测试接口
	r.GET("/ping", func(c *gin.Context) {
		facade.Response().Success(c, errcode.Success().WithMsg("pong"))
	})

	// GET 带参数测试
	r.GET("/echo", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		facade.Response().Success(c, errcode.Success().WithData(gin.H{
			"name": name,
			"age":  age,
		}))
	})

	// POST JSON测试
	r.POST("/echo", func(c *gin.Context) {
		var body map[string]any
		if err := c.BindJSON(&body); err != nil {
			facade.Response().Error(c, errcode.ArgsError().WithMsg(err.Error()))
			return
		}
		facade.Response().Success(c, errcode.Success().WithData(body))
	})

	// POST 表单测试
	r.POST("/form", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		facade.Response().Success(c, errcode.Success().WithData(gin.H{
			"name":  name,
			"email": email,
		}))
	})

	// 单文件上传测试
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			facade.Response().Error(c, errcode.ArgsError().WithMsg(err.Error()))
			return
		}

		description := c.PostForm("description")

		facade.Response().Success(c, errcode.Success().
			WithMsg("upload success").
			WithData(gin.H{
				"filename":    file.Filename,
				"size":        file.Size,
				"description": description,
			}))
	})

	// 多文件上传测试
	r.POST("/multi-upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			facade.Response().Error(c, errcode.ArgsError().WithMsg(err.Error()))
			return
		}

		files := form.File["files"]
		var fileInfos []gin.H
		for _, file := range files {
			fileInfos = append(fileInfos, gin.H{
				"filename": file.Filename,
				"size":     file.Size,
			})
		}

		facade.Response().Success(c, errcode.Success().
			WithMsg("multi upload success").
			WithData(gin.H{
				"files":      fileInfos,
				"file_count": len(fileInfos),
			}))
	})

	// 带自定义字段名的文件上传
	r.POST("/upload-custom", func(c *gin.Context) {
		file, err := c.FormFile("custom_file")
		if err != nil {
			facade.Response().Error(c, errcode.ArgsError().WithMsg(err.Error()))
			return
		}

		facade.Response().Success(c, errcode.Success().
			WithMsg("upload success").
			WithData(gin.H{
				"filename": file.Filename,
				"size":     file.Size,
			}))
	})

	// 错误响应测试
	r.GET("/error", func(c *gin.Context) {
		facade.Response().Error(c, errcode.SystemError().WithMsg("internal server error"))
	})

	// 延迟响应测试
	r.GET("/delay", func(c *gin.Context) {
		delay := c.Query("delay")
		if delay != "" {
			d, _ := time.ParseDuration(delay)
			time.Sleep(d)
		}
		facade.Response().Success(c, errcode.Success().WithMsg("ok"))
	})

	return httptest.NewServer(r)
}

// TestHttpRequest 测试基本HTTP请求
func TestHttpRequest(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIDKey, "test-trace-id")

	// 测试GET请求
	resp := facade.Http().Send(ctx, "GET", ts.URL+"/ping")
	require.NoError(t, resp.ErrMsg)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	err := json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	require.Equal(t, float64(0), result["code"])
	require.Equal(t, "pong", result["msg"])
	require.Empty(t, result["data"])
}

// TestHttpTraceCollected 测试普通HTTP请求记录调试信息
func TestHttpTraceCollected(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	traceID := "test-http-trace-collected"
	store := facade.Debugger().Store()
	store.Delete(traceID)
	t.Cleanup(func() {
		store.Delete(traceID)
	})

	ctx := context.WithValue(t.Context(), ctxkey.TraceIDKey, traceID)
	resp := facade.Http().Send(ctx, http.MethodGet, ts.URL+"/ping")
	require.NoError(t, resp.ErrMsg)

	trace, ok := store.Get(traceID)
	require.True(t, ok)
	require.Len(t, trace.HTTP, 1)
	require.Equal(t, http.MethodGet, trace.HTTP[0].Method)
	require.Equal(t, ts.URL+"/ping", trace.HTTP[0].URL)
	require.Equal(t, http.StatusOK, trace.HTTP[0].Status)
	require.GreaterOrEqual(t, trace.HTTP[0].Ms, float64(0))
}

// TestHttpWithoutTrace 测试普通HTTP请求跳过调试信息
func TestHttpWithoutTrace(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	traceID := "test-http-without-trace"
	store := facade.Debugger().Store()
	store.Delete(traceID)
	t.Cleanup(func() {
		store.Delete(traceID)
	})

	ctx := context.WithValue(t.Context(), ctxkey.TraceIDKey, traceID)
	resp := facade.Http().Send(h.WithoutTrace(ctx), http.MethodGet, ts.URL+"/ping")
	require.NoError(t, resp.ErrMsg)

	trace, ok := store.Get(traceID)
	require.False(t, ok)
	require.Nil(t, trace)
}

// TestHttpClientBuilder 测试HTTP客户端链式配置
func TestHttpClientBuilder(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	type EchoResponse struct {
		Code int            `json:"code"`
		Msg  string         `json:"msg"`
		Data map[string]any `json:"data"`
	}

	queryResp, response := facade.Http().
		WithQuery(map[string]any{"name": "张三", "age": "18"}).
		SendAsJson[EchoResponse](t.Context(), http.MethodGet, ts.URL+"/echo")
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "张三", queryResp.Data["name"])
	require.Equal(t, "18", queryResp.Data["age"])

	bodyResp, response := facade.Http().
		WithHeader("Content-Type", "application/json").
		WithBody(map[string]any{"name": "李四", "email": "lisi@example.com"}).
		SendAsJson[EchoResponse](t.Context(), http.MethodPost, ts.URL+"/echo")
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "李四", bodyResp.Data["name"])
	require.Equal(t, "lisi@example.com", bodyResp.Data["email"])

	type FormResponse struct {
		Code int `json:"code"`
		Data struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
	}

	formResp, response := facade.Http().
		WithForm(map[string]any{"name": "王五", "email": "wangwu@example.com"}).
		SendAsJson[FormResponse](t.Context(), http.MethodPost, ts.URL+"/form")
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "王五", formResp.Data.Name)
	require.Equal(t, "wangwu@example.com", formResp.Data.Email)

	headerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"x-custom-header": r.Header.Get("X-Custom-Header"),
				"authorization":   r.Header.Get("Authorization"),
			},
		})
	}))
	defer headerServer.Close()

	headerResp, response := facade.Http().
		WithHeaders(map[string]string{
			"X-Custom-Header": "custom-value",
			"Authorization":   "Bearer token123",
		}).
		SendAsJson[EchoResponse](t.Context(), http.MethodGet, headerServer.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "custom-value", headerResp.Data["x-custom-header"])
	require.Equal(t, "Bearer token123", headerResp.Data["authorization"])
}

// TestHttpClientBuilderIsolation 测试链式配置不会污染基础客户端
func TestHttpClientBuilderIsolation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"base":   r.Header.Get("X-Base"),
				"first":  r.Header.Get("X-First"),
				"second": r.Header.Get("X-Second"),
			},
		})
	}))
	defer ts.Close()

	type HeaderResponse struct {
		Data map[string]any `json:"data"`
	}

	base := facade.Http().WithHeader("X-Base", "base")

	first, response := base.
		WithHeader("X-First", "first").
		SendAsJson[HeaderResponse](t.Context(), http.MethodGet, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "base", first.Data["base"])
	require.Equal(t, "first", first.Data["first"])
	require.Empty(t, first.Data["second"])

	second, response := base.
		WithHeader("X-Second", "second").
		SendAsJson[HeaderResponse](t.Context(), http.MethodGet, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "base", second.Data["base"])
	require.Empty(t, second.Data["first"])
	require.Equal(t, "second", second.Data["second"])

	plain, response := base.SendAsJson[HeaderResponse](t.Context(), http.MethodGet, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "base", plain.Data["base"])
	require.Empty(t, plain.Data["first"])
	require.Empty(t, plain.Data["second"])
}

// TestHttpBuildURLMergesExistingQuery 测试query参数合并已有参数
func TestHttpBuildURLMergesExistingQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"a":     r.URL.Query().Get("a"),
				"keep":  r.URL.Query().Get("keep"),
				"added": r.URL.Query().Get("added"),
			},
		})
	}))
	defer ts.Close()

	type QueryResponse struct {
		Data map[string]any `json:"data"`
	}

	resp, response := facade.Http().
		WithQuery(map[string]any{"a": "new", "added": "yes"}).
		SendAsJson[QueryResponse](t.Context(), http.MethodGet, ts.URL+"/query?keep=yes&a=old")
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "new", resp.Data["a"])
	require.Equal(t, "yes", resp.Data["keep"])
	require.Equal(t, "yes", resp.Data["added"])
}

// TestHttpSendAsJsonInvalidJson 测试非JSON响应解析失败
func TestHttpSendAsJsonInvalidJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("not-json"))
	}))
	defer ts.Close()

	resp, response := facade.Http().
		SendAsJson[map[string]any](t.Context(), http.MethodGet, ts.URL)
	require.Error(t, response.ErrMsg)
	require.Contains(t, response.ErrMsg.Error(), "json解析失败")
	require.Nil(t, resp)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

// TestHttpSendAsJsonEmptyBody 测试空响应体解析
func TestHttpSendAsJsonEmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	type EmptyResponse struct {
		Value string `json:"value"`
	}

	resp, response := facade.Http().
		SendAsJson[EmptyResponse](t.Context(), http.MethodGet, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, http.StatusNoContent, response.StatusCode)
	require.NotNil(t, resp)
	require.Empty(t, resp.Value)
}

// TestHttpMethods 测试常用请求方法及nil上下文
func TestHttpMethods(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"method": r.Method,
		})
	}))
	defer ts.Close()

	getResp := facade.Http().Get(nil, ts.URL+"/method")
	require.NoError(t, getResp.ErrMsg)
	require.Contains(t, string(getResp.Body), http.MethodGet)

	postResp := facade.Http().Post(t.Context(), ts.URL+"/method")
	require.NoError(t, postResp.ErrMsg)
	require.Contains(t, string(postResp.Body), http.MethodPost)

	putResp := facade.Http().Put(t.Context(), ts.URL+"/method")
	require.NoError(t, putResp.ErrMsg)
	require.Contains(t, string(putResp.Body), http.MethodPut)

	deleteResp := facade.Http().Delete(t.Context(), ts.URL+"/method")
	require.NoError(t, deleteResp.ErrMsg)
	require.Contains(t, string(deleteResp.Body), http.MethodDelete)
}

// TestHttpContentType 测试请求体类型自动识别和自定义覆盖
func TestHttpContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"content_type": r.Header.Get("Content-Type"),
		})
	}))
	defer ts.Close()

	type ContentTypeResponse struct {
		ContentType string `json:"content_type"`
	}

	jsonResp, response := facade.Http().
		WithBody(map[string]any{"name": "test"}).
		SendAsJson[ContentTypeResponse](t.Context(), http.MethodPost, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "application/json", jsonResp.ContentType)

	formResp, response := facade.Http().
		WithForm(map[string]any{"name": "test"}).
		SendAsJson[ContentTypeResponse](t.Context(), http.MethodPost, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "application/x-www-form-urlencoded", formResp.ContentType)

	textResp, response := facade.Http().
		WithBody("plain text").
		SendAsJson[ContentTypeResponse](t.Context(), http.MethodPost, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "text/plain", textResp.ContentType)

	customResp, response := facade.Http().
		WithHeader("Content-Type", "application/custom").
		WithBody(map[string]any{"name": "test"}).
		SendAsJson[ContentTypeResponse](t.Context(), http.MethodPost, ts.URL)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, "application/custom", customResp.ContentType)
}

// TestHttpSendAsJsonJson 测试JSON响应解析
func TestHttpSendAsJsonJson(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIDKey, "test-trace-id")

	// 测试GET请求并解析JSON
	resp, response := facade.Http().SendAsJson[errcode.SuccessResponse](
		ctx,
		"GET",
		ts.URL+"/ping",
	)
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, "pong", resp.Msg)
	require.Empty(t, resp.Data)
}

// TestHttpRequestWithQuery 测试带Query参数的请求
func TestHttpRequestWithQuery(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	type EchoResponse struct {
		Code int            `json:"code"`
		Msg  string         `json:"msg"`
		Data map[string]any `json:"data"`
	}

	resp, response := facade.Http().
		WithQuery(map[string]any{"name": "张三", "age": "18"}).
		SendAsJson[EchoResponse](ctx, "GET", ts.URL+"/echo")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "张三", resp.Data["name"])
	require.Equal(t, "18", resp.Data["age"])
}

// TestHttpForm 测试普通表单提交
func TestHttpForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	type FormResponse struct {
		Code int `json:"code"`
		Data struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithForm(map[string]any{"name": "张三", "email": "zhangsan@example.com"}).
		WithTimeout(30 * time.Second).
		SendAsJson[FormResponse](ctx, "POST", ts.URL+"/form")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "张三", resp.Data.Name)
	require.Equal(t, "zhangsan@example.com", resp.Data.Email)
}

// TestHttpRequestWithBody 测试带Body的POST请求
func TestHttpRequestWithBody(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	bodyData := map[string]any{
		"name":  "李四",
		"email": "lisi@example.com",
	}
	body, _ := json.Marshal(bodyData)

	type EchoResponse struct {
		Code int            `json:"code"`
		Msg  string         `json:"msg"`
		Data map[string]any `json:"data"`
	}

	resp, response := facade.Http().
		WithHeader("Content-Type", "application/json").
		WithBody(body).
		SendAsJson[EchoResponse](ctx, "POST", ts.URL+"/echo")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "李四", resp.Data["name"])
	require.Equal(t, "lisi@example.com", resp.Data["email"])
}

// TestHttpRequestWithForm 测试表单请求
func TestHttpRequestWithForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	type FormResponse struct {
		Code int            `json:"code"`
		Msg  string         `json:"msg"`
		Data map[string]any `json:"data"`
	}

	resp, response := facade.Http().
		WithForm(map[string]any{"name": "王五", "email": "wangwu@example.com"}).
		SendAsJson[FormResponse](ctx, "POST", ts.URL+"/form")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "王五", resp.Data["name"])
	require.Equal(t, "wangwu@example.com", resp.Data["email"])
}

// TestHttpRequestWithHeaders 测试带自定义Headers的请求
func TestHttpRequestWithHeaders(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	// 创建一个可以读取header的测试服务器
	r := gin.Default()
	r.GET("/headers", func(c *gin.Context) {
		customHeader := c.GetHeader("X-Custom-Header")
		authHeader := c.GetHeader("Authorization")
		facade.Response().Success(c, errcode.Success().WithData(gin.H{
			"x-custom-header": customHeader,
			"authorization":   authHeader,
		}))
	})
	ts2 := httptest.NewServer(r)
	defer ts2.Close()

	type HeaderResponse struct {
		Code int            `json:"code"`
		Data map[string]any `json:"data"`
	}

	resp, response := facade.Http().
		WithHeaders(map[string]string{
			"X-Custom-Header": "custom-value",
			"Authorization":   "Bearer token123",
		}).
		SendAsJson[HeaderResponse](ctx, "GET", ts2.URL+"/headers")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, "custom-value", resp.Data["x-custom-header"])
	require.Equal(t, "Bearer token123", resp.Data["authorization"])
}

// TestHttpRequestWithTimeout 测试超时
func TestHttpRequestWithTimeout(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	// 设置100ms超时,但服务器会延迟200ms
	response := facade.Http().
		WithTimeout(100*time.Millisecond).
		Send(ctx, "GET", ts.URL+"/delay?delay=200ms")
	require.Error(t, response.ErrMsg)
	require.Contains(t, response.ErrMsg.Error(), "请求失败")
}

// TestHttpRequestPostForm 测试POST表单请求(使用Request方法)
func TestHttpRequestPostForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	resp := facade.Http().
		WithForm(map[string]any{"name": "赵六", "email": "zhaoliu@example.com"}).
		Send(ctx, "POST", ts.URL+"/form")
	require.NoError(t, resp.ErrMsg)

	var result map[string]any
	err := json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)

	data := result["data"].(map[string]any)
	require.Equal(t, "赵六", data["name"])
	require.Equal(t, "zhaoliu@example.com", data["email"])
}

// TestHttpSendKeepsErrorStatus 测试保留错误状态码和响应体
func TestHttpSendKeepsErrorStatus(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	response := facade.Http().Send(ctx, "GET", ts.URL+"/error")
	require.NoError(t, response.ErrMsg)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Contains(t, string(response.Body), "internal server error")
}

// TestHttpSendToJsonErrorResponse 测试错误状态码JSON响应解析
func TestHttpSendToJsonErrorResponse(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	type ErrorResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	resp, response := facade.Http().SendAsJson[ErrorResponse](
		ctx,
		"GET",
		ts.URL+"/error",
	)
	require.NoError(t, response.ErrMsg)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.NotNil(t, resp)
	require.Equal(t, 500, resp.Code)
	require.Equal(t, "internal server error", resp.Msg)
}

// TestHttpRequestInvalidURL 测试无效URL
func TestHttpRequestInvalidURL(t *testing.T) {
	ctx := t.Context()

	response := facade.Http().Send(ctx, "GET", "http://invalid.url.that.does.not.exist")
	require.Error(t, response.ErrMsg)
}

// 创建测试文件
func createTestFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)

	return tmpFile.Name()
}

// TestHttpUploadFile 测试单文件上传
func TestHttpUploadFile(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	// 创建测试文件
	testContent := "Hello, this is a test file"
	testFilePath := createTestFile(t, testContent)
	defer func() { _ = os.Remove(testFilePath) }()

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFile("file", h.File{
			FilePath:  testFilePath,
			FieldName: "file",
		}).
		WithForm(map[string]any{
			"description": "测试文件上传",
		}).
		WithTimeout(30 * time.Second).
		SendAsJson[UploadResponse](ctx, "POST", ts.URL+"/upload")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "upload success", resp.Msg)
	// 只检查文件名前缀因为临时文件名是动态的
	require.Contains(t, resp.Data.Filename, "test-")
	require.Equal(t, int64(len(testContent)), resp.Data.Size)
	require.Equal(t, "测试文件上传", resp.Data.Description)
}

// TestHttpUploadFileWithData 测试使用字节数据上传文件
func TestHttpUploadFileWithData(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	fileData := []byte("This is file content from byte data")

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFiles(map[string]h.File{
			"file": {
				FileData:  fileData,
				FileName:  "test.txt",
				FieldName: "file",
			},
		}).
		WithForm(map[string]any{
			"description": "使用字节数据上传",
		}).
		WithTimeout(30 * time.Second).
		SendAsJson[UploadResponse](ctx, "POST", ts.URL+"/upload")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "test.txt", resp.Data.Filename)
	require.Equal(t, int64(len(fileData)), resp.Data.Size)
	require.Equal(t, "使用字节数据上传", resp.Data.Description)
}

// TestHttpUploadMultipleFiles 测试多文件上传
func TestHttpUploadMultipleFiles(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	// 创建多个测试文件
	file1Content := "Content of file 1"
	file1Path := createTestFile(t, file1Content)
	defer func() { _ = os.Remove(file1Path) }()

	file2Content := "Content of file 2"
	file2Path := createTestFile(t, file2Content)
	defer func() { _ = os.Remove(file2Path) }()

	file3Data := []byte("Content of file 3")

	type MultiUploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Files     []gin.H `json:"files"`
			FileCount int     `json:"file_count"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFiles(map[string]h.File{
			"file1": {
				FilePath:  file1Path,
				FieldName: "files",
			},
			"file2": {
				FilePath:  file2Path,
				FieldName: "files",
			},
			"file3": {
				FileData:  file3Data,
				FileName:  "file3.txt",
				FieldName: "files",
			},
		}).
		WithTimeout(30 * time.Second).
		SendAsJson[MultiUploadResponse](ctx, "POST", ts.URL+"/multi-upload")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "multi upload success", resp.Msg)
	require.Equal(t, 3, resp.Data.FileCount)
}

// TestHttpUploadFileWithCustomFieldName 测试自定义表单字段名的文件上传
func TestHttpUploadFileWithCustomFieldName(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	testContent := "Test file with custom field name"
	testFilePath := createTestFile(t, testContent)
	defer func() { _ = os.Remove(testFilePath) }()

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename string `json:"filename"`
			Size     int64  `json:"size"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFiles(map[string]h.File{
			"myfile": {
				FilePath:  testFilePath,
				FieldName: "custom_file", // 自定义字段名
			},
		}).
		WithTimeout(30 * time.Second).
		SendAsJson[UploadResponse](ctx, "POST", ts.URL+"/upload-custom")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "upload success", resp.Msg)
}

// TestHttpUploadFileWithRequestMethod 测试使用Request方法上传文件
func TestHttpUploadFileWithRequestMethod(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	testContent := "Test file using Request method"
	testFilePath := createTestFile(t, testContent)
	defer func() { _ = os.Remove(testFilePath) }()

	resp := facade.Http().
		WithFiles(map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		}).
		WithForm(map[string]any{
			"description": "使用Request方法上传",
		}).
		WithTimeout(30*time.Second).
		Send(ctx, "POST", ts.URL+"/upload")
	require.NoError(t, resp.ErrMsg)

	var result map[string]any
	err := json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	require.Equal(t, float64(0), result["code"])
	require.Equal(t, "upload success", result["msg"])
}

// TestHttpUploadFileWithTimeout 测试文件上传超时
func TestHttpUploadFileWithTimeout(t *testing.T) {
	// 创建一个延迟响应的上传服务器
	r := gin.Default()
	r.POST("/slow-upload", func(c *gin.Context) {
		time.Sleep(500 * time.Millisecond)
		file, err := c.FormFile("file")
		if err != nil {
			facade.Response().Error(c, errcode.ArgsError().WithMsg(err.Error()))
			return
		}
		facade.Response().Success(c, errcode.Success().
			WithMsg("upload success").
			WithData(gin.H{"filename": file.Filename}))
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := t.Context()

	testContent := "Test timeout file"
	testFilePath := createTestFile(t, testContent)
	defer func() { _ = os.Remove(testFilePath) }()

	response := facade.Http().
		WithFiles(map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		}).
		WithTimeout(100*time.Millisecond). // 100ms超时,但服务器需要500ms
		Send(ctx, "POST", ts.URL+"/slow-upload")
	require.Error(t, response.ErrMsg)
	require.Contains(t, response.ErrMsg.Error(), "请求失败")
}

// TestHttpUploadFileLargeFile 测试大文件上传(模拟)
func TestHttpUploadFileLargeFile(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	// 创建一个大文件(1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte('A' + i%26)
	}

	testFilePath := createTestFile(t, string(largeContent))
	defer func() { _ = os.Remove(testFilePath) }()

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFiles(map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		}).
		WithForm(map[string]any{
			"description": "大文件上传测试",
		}).
		WithTimeout(60 * time.Second). // 大文件需要更长的超时时间
		SendAsJson[UploadResponse](ctx, "POST", ts.URL+"/upload")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, int64(len(largeContent)), resp.Data.Size)
}

// TestHttpUploadFileWithContext 测试带Context的文件上传
func TestHttpUploadFileWithContext(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIDKey, "upload-trace-id")

	testContent := "Test file with context"
	testFilePath := createTestFile(t, testContent)
	defer func() { _ = os.Remove(testFilePath) }()

	type UploadResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFiles(map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		}).
		WithTimeout(30 * time.Second).
		SendAsJson[UploadResponse](ctx, "POST", ts.URL+"/upload")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
}

// TestHttpUploadFileAndFormData 测试同时上传文件和表单数据
func TestHttpUploadFileAndFormData(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := t.Context()

	testContent := "Test file with form data"
	testFilePath := createTestFile(t, testContent)
	defer func() { _ = os.Remove(testFilePath) }()

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, response := facade.Http().
		WithFiles(map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		}).
		WithForm(map[string]any{
			"description": "文件描述",
			"user_id":     "12345",
			"category":    "test",
		}).
		WithTimeout(30 * time.Second).
		SendAsJson[UploadResponse](ctx, "POST", ts.URL+"/upload")
	require.NoError(t, response.ErrMsg)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "文件描述", resp.Data.Description)
}
