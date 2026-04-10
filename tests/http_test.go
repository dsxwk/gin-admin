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
	"os"
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

	// 单文件上传测试
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
			return
		}

		description := c.PostForm("description")

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "upload success",
			"data": gin.H{
				"filename":    file.Filename,
				"size":        file.Size,
				"description": description,
			},
		})
	})

	// 多文件上传测试
	r.POST("/multi-upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
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

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "multi upload success",
			"data": gin.H{
				"files":      fileInfos,
				"file_count": len(fileInfos),
			},
		})
	})

	// 带自定义字段名的文件上传
	r.POST("/upload-custom", func(c *gin.Context) {
		file, err := c.FormFile("custom_file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "upload success",
			"data": gin.H{
				"filename": file.Filename,
				"size":     file.Size,
			},
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
	resp, err := facade.Http.Send(ctx, "GET", ts.URL+"/ping", nil)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	require.NoError(t, err)
	require.Equal(t, float64(0), result["code"])
	require.Equal(t, "pong", result["msg"])
}

// TestHttpSendToJson 测试JSON响应解析
func TestHttpSendToJson(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-trace-id")

	// 测试GET请求并解析JSON
	resp, err := facade.HttpAs[errcode.SuccessResponse]().Send(
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
		Query: map[string]interface{}{
			"name": "张三",
			"age":  "18",
		},
	}

	type EchoResponse struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}

	resp, err := facade.HttpAs[EchoResponse]().Send(
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

// TestHttpForm 测试普通表单提交
func TestHttpForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	// 使用普通表单接口
	opt := &h.Option{
		Form: map[string]interface{}{
			"name":  "张三",
			"email": "zhangsan@example.com",
		},
		Timeout: 30 * time.Second,
	}

	type FormResponse struct {
		Code int `json:"code"`
		Data struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[FormResponse]().Send(ctx, "POST", ts.URL+"/form", opt)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "张三", resp.Data.Name)
	require.Equal(t, "zhangsan@example.com", resp.Data.Email)
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

	resp, err := facade.HttpAs[EchoResponse]().Send(
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
		Form: map[string]interface{}{
			"name":  "王五",
			"email": "wangwu@example.com",
		},
	}

	type FormResponse struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}

	resp, err := facade.HttpAs[FormResponse]().Send(
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

	resp, err := facade.HttpAs[HeaderResponse]().Send(
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

	_, err := facade.Http.Send(ctx, "GET", ts.URL+"/delay?delay=200ms", opt)
	require.Error(t, err)
	require.Contains(t, err.Error(), "请求失败")
}

// TestHttpRequestPostForm 测试POST表单请求(使用Request方法)
func TestHttpRequestPostForm(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	opt := &h.Option{
		Form: map[string]interface{}{
			"name":  "赵六",
			"email": "zhaoliu@example.com",
		},
	}

	resp, err := facade.Http.Send(ctx, "POST", ts.URL+"/form", opt)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
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

	_, err := facade.Http.Send(ctx, "GET", ts.URL+"/error", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "状态码: 500")
}

// TestHttpSendToJsonErrorResponse 测试JSON解析错误响应
func TestHttpSendToJsonErrorResponse(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	type ErrorResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	resp, err := facade.HttpAs[ErrorResponse]().Send(
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

	_, err := facade.Http.Send(ctx, "GET", "http://invalid.url.that.does.not.exist", nil)
	require.Error(t, err)
}

// 创建测试文件
func createTestFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)

	return tmpFile.Name()
}

// TestHttpUploadFile 测试单文件上传
func TestHttpUploadFile(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	// 创建测试文件
	testContent := "Hello, this is a test file"
	testFilePath := createTestFile(t, testContent)
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		},
		Form: map[string]interface{}{
			"description": "测试文件上传",
		},
		Timeout: 30 * time.Second,
	}

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[UploadResponse]().Send(ctx, "POST", ts.URL+"/upload", opt)
	require.NoError(t, err)
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

	ctx := context.Background()

	fileData := []byte("This is file content from byte data")

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FileData:  fileData,
				FileName:  "test.txt",
				FieldName: "file",
			},
		},
		Form: map[string]interface{}{
			"description": "使用字节数据上传",
		},
		Timeout: 30 * time.Second,
	}

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[UploadResponse]().Send(ctx, "POST", ts.URL+"/upload", opt)
	require.NoError(t, err)
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

	ctx := context.Background()

	// 创建多个测试文件
	file1Content := "Content of file 1"
	file1Path := createTestFile(t, file1Content)
	defer os.Remove(file1Path)

	file2Content := "Content of file 2"
	file2Path := createTestFile(t, file2Content)
	defer os.Remove(file2Path)

	file3Data := []byte("Content of file 3")

	opt := &h.Option{
		Files: map[string]h.File{
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
		},
		Timeout: 30 * time.Second,
	}

	type MultiUploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Files     []gin.H `json:"files"`
			FileCount int     `json:"file_count"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[MultiUploadResponse]().Send(ctx, "POST", ts.URL+"/multi-upload", opt)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "multi upload success", resp.Msg)
	require.Equal(t, 3, resp.Data.FileCount)
}

// TestHttpUploadFileWithCustomFieldName 测试自定义表单字段名的文件上传
func TestHttpUploadFileWithCustomFieldName(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	testContent := "Test file with custom field name"
	testFilePath := createTestFile(t, testContent)
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"myfile": {
				FilePath:  testFilePath,
				FieldName: "custom_file", // 自定义字段名
			},
		},
		Timeout: 30 * time.Second,
	}

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename string `json:"filename"`
			Size     int64  `json:"size"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[UploadResponse]().Send(ctx, "POST", ts.URL+"/upload-custom", opt)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "upload success", resp.Msg)
}

// TestHttpUploadFileWithRequestMethod 测试使用Request方法上传文件
func TestHttpUploadFileWithRequestMethod(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	testContent := "Test file using Request method"
	testFilePath := createTestFile(t, testContent)
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		},
		Form: map[string]interface{}{
			"description": "使用Request方法上传",
		},
		Timeout: 30 * time.Second,
	}

	resp, err := facade.Http.Send(ctx, "POST", ts.URL+"/upload", opt)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
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
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "upload success",
			"data": gin.H{"filename": file.Filename},
		})
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.Background()

	testContent := "Test timeout file"
	testFilePath := createTestFile(t, testContent)
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		},
		Timeout: 100 * time.Millisecond, // 100ms超时,但服务器需要500ms
	}

	_, err := facade.Http.Send(ctx, "POST", ts.URL+"/slow-upload", opt)
	require.Error(t, err)
	require.Contains(t, err.Error(), "请求失败")
}

// TestHttpUploadFileLargeFile 测试大文件上传(模拟)
func TestHttpUploadFileLargeFile(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	// 创建一个大文件(1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte('A' + i%26)
	}

	testFilePath := createTestFile(t, string(largeContent))
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		},
		Form: map[string]interface{}{
			"description": "大文件上传测试",
		},
		Timeout: 60 * time.Second, // 大文件需要更长的超时时间
	}

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[UploadResponse]().Send(ctx, "POST", ts.URL+"/upload", opt)
	require.NoError(t, err)
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
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "upload-trace-id")

	testContent := "Test file with context"
	testFilePath := createTestFile(t, testContent)
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		},
		Timeout: 30 * time.Second,
	}

	type UploadResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	resp, err := facade.HttpAs[UploadResponse]().Send(ctx, "POST", ts.URL+"/upload", opt)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
}

// TestHttpUploadFileAndFormData 测试同时上传文件和表单数据
func TestHttpUploadFileAndFormData(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()

	ctx := context.Background()

	testContent := "Test file with form data"
	testFilePath := createTestFile(t, testContent)
	defer os.Remove(testFilePath)

	opt := &h.Option{
		Files: map[string]h.File{
			"file": {
				FilePath:  testFilePath,
				FieldName: "file",
			},
		},
		Form: map[string]interface{}{
			"description": "文件描述",
			"user_id":     "12345",
			"category":    "test",
		},
		Timeout: 30 * time.Second,
	}

	type UploadResponse struct {
		Code int `json:"code"`
		Data struct {
			Filename    string `json:"filename"`
			Size        int64  `json:"size"`
			Description string `json:"description"`
		} `json:"data"`
		Msg string `json:"msg"`
	}

	resp, err := facade.HttpAs[UploadResponse]().Send(ctx, "POST", ts.URL+"/upload", opt)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 0, resp.Code)
	require.Equal(t, "文件描述", resp.Data.Description)
}
