package http

import (
	"bytes"
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/message"
	"github.com/goccy/go-json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultTimeout = 5 * time.Second

var defaultClient *http.Client

// InitClient 初始化全局HTTP客户端
func InitClient() {
	if defaultClient == nil {
		defaultClient = &http.Client{
			Transport: &TracingTransport{
				Transport: &http.Transport{
					MaxConnsPerHost: 100,
				},
			},
		}
	}
}

// GetClient 获取全局HTTP客户端
func GetClient() *http.Client {
	if defaultClient == nil {
		InitClient()
	}
	return defaultClient
}

// Client HTTP客户端
type Client[T any] struct {
	timeout time.Duration
}

// NewClient 创建HTTP客户端
func NewClient[T any]() *Client[T] {
	return &Client[T]{
		timeout: defaultTimeout,
	}
}

// WithTimeout 自定义超时的HTTP客户端
func (c *Client[T]) WithTimeout(timeout time.Duration) *Client[T] {
	return &Client[T]{
		timeout: timeout,
	}
}

// File 文件
type File struct {
	FileName  string // 文件名
	FilePath  string // 文件路径
	FileData  []byte // 文件数据
	FieldName string // 表单字段名
}

type Option struct {
	Headers map[string]string      // 请求头
	Query   map[string]interface{} // query参数
	Form    map[string]interface{} // form参数
	Body    interface{}            // 请求体
	Files   map[string]File        // 文件上传字段
	Timeout time.Duration          // 超时时间
}

// Send 发送HTTP请求
func (c *Client[T]) Send(ctx context.Context, method, uri string, opt *Option) ([]byte, error) {
	return c.doSend(ctx, method, uri, opt)
}

// Get 发送GET请求
func (c *Client[T]) Get(ctx context.Context, uri string, opt *Option) ([]byte, error) {
	return c.Send(ctx, "GET", uri, opt)
}

// Post 发送POST请求
func (c *Client[T]) Post(ctx context.Context, uri string, opt *Option) ([]byte, error) {
	return c.Send(ctx, "POST", uri, opt)
}

// Put 发送PUT请求
func (c *Client[T]) Put(ctx context.Context, uri string, opt *Option) ([]byte, error) {
	return c.Send(ctx, "PUT", uri, opt)
}

// Delete 发送DELETE请求
func (c *Client[T]) Delete(ctx context.Context, uri string, opt *Option) ([]byte, error) {
	return c.Send(ctx, "DELETE", uri, opt)
}

// AsJson 将响应体解析为T类型
func (c *Client[T]) AsJson(data []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("json解析失败: %w\n响应内容:\n%s", err, data)
	}
	return &result, nil
}

// SendAsJson 发送请求并解析为T类型
func (c *Client[T]) SendAsJson(ctx context.Context, method, uri string, opt *Option) (*T, error) {
	data, err := c.doSend(ctx, method, uri, opt)
	if err != nil {
		return nil, err
	}
	return c.AsJson(data)
}

// doSend 发送请求
func (c *Client[T]) doSend(ctx context.Context, method, uri string, opt *Option) ([]byte, error) {
	if opt == nil {
		opt = &Option{}
	}

	requestTimeout := opt.Timeout
	if requestTimeout == 0 {
		requestTimeout = c.timeout
	}

	method = strings.ToUpper(method)
	uri = c.buildUrl(uri, opt.Query)

	// 判断是否有文件上传
	if opt.Files != nil && len(opt.Files) > 0 {
		return c.doFileUpload(ctx, uri, opt, requestTimeout)
	}

	// 构建请求体
	var bodyReader io.Reader
	contentType := ""
	if opt.Form != nil && len(opt.Form) > 0 {
		values := url.Values{}
		for k, v := range opt.Form {
			values.Set(k, fmt.Sprintf("%v", v))
		}
		bodyReader = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else if opt.Body != nil {
		switch v := opt.Body.(type) {
		case []byte:
			bodyReader = bytes.NewReader(v)
			contentType = "application/octet-stream"
		case string:
			bodyReader = strings.NewReader(v)
			contentType = "text/plain"
		case *bytes.Buffer:
			bodyReader = v
			contentType = "application/octet-stream"
		case *strings.Reader:
			bodyReader = v
			contentType = "text/plain"
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("JSON序列化失败: %w", err)
			}
			bodyReader = bytes.NewReader(jsonBytes)
			contentType = "application/json"
		}
	}

	// 通过上下文应用超时
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, uri, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认Content-Type
	if opt.Headers == nil {
		opt.Headers = make(map[string]string)
	}
	if contentType != "" && opt.Headers["Content-Type"] == "" {
		opt.Headers["Content-Type"] = contentType
	}
	if contentType == "" && opt.Headers["Content-Type"] == "" {
		opt.Headers["Content-Type"] = "application/json"
	}

	// 设置自定义headers
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := GetClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}

// doFileUpload 文件上传
func (c *Client[T]) doFileUpload(ctx context.Context, uri string, opt *Option, requestTimeout time.Duration) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, file := range opt.Files {
		var fileData []byte
		var fileName string

		if file.FilePath != "" {
			data, err := os.ReadFile(file.FilePath)
			if err != nil {
				return nil, fmt.Errorf("读取文件 %s 失败: %w", file.FilePath, err)
			}
			fileData = data
			fileName = filepath.Base(file.FilePath)
		} else if file.FileData != nil {
			fileData = file.FileData
			fileName = file.FileName
			if fileName == "" {
				fileName = fieldName
			}
		} else {
			continue
		}

		formFieldName := file.FieldName
		if formFieldName == "" {
			formFieldName = fieldName
		}

		part, err := writer.CreateFormFile(formFieldName, fileName)
		if err != nil {
			return nil, fmt.Errorf("创建表单文件失败: %w", err)
		}

		if _, err = part.Write(fileData); err != nil {
			return nil, fmt.Errorf("写入文件数据失败: %w", err)
		}
	}

	for k, v := range opt.Form {
		if err := writer.WriteField(k, fmt.Sprintf("%v", v)); err != nil {
			return nil, fmt.Errorf("写入表单字段失败: %w", err)
		}
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭writer失败: %w", err)
	}

	// 通过上下文应用超时
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", uri, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := GetClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("上传失败, 状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// buildURL 拼接get请求query参数
func (c *Client[T]) buildUrl(baseURL string, query map[string]interface{}) string {
	if len(query) == 0 {
		return baseURL
	}
	q := url.Values{}
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	if strings.Contains(baseURL, "?") {
		return baseURL + "&" + q.Encode()
	}
	return baseURL + "?" + q.Encode()
}

// TracingTransport 传输中间件追踪
type TracingTransport struct {
	Transport http.RoundTripper
}

// RoundTrip 实现http.RoundTripper
func (t *TracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// 读取并缓存请求体用于tracing
	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
	}

	resp, err := t.Transport.RoundTrip(req)

	// 收集响应信息用于tracing
	var (
		respBodyBytes []byte
		respStatus    int
	)
	if err == nil {
		respStatus = resp.StatusCode
		respBodyBytes, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(respBodyBytes))
	}

	// 发布trace事件
	costMs := float64(time.Since(start).Nanoseconds()) / 1e6

	traceId := "unknown"
	if req.Context() != nil {
		if id := req.Context().Value(ctxkey.TraceIdKey); id != nil {
			if s, ok := id.(string); ok && s != "" {
				traceId = s
			}
		}
	}

	var respJson interface{}
	if len(respBodyBytes) > 0 {
		if err = json.Unmarshal(respBodyBytes, &respJson); err != nil {
			respJson = respBodyBytes
		}
	}

	message.NewEvent().Publish(debugger.TopicHttp, debugger.HttpEvent{
		TraceId:  traceId,
		Url:      req.URL.String(),
		Method:   req.Method,
		Header:   headersToMap(req.Header),
		Body:     string(reqBodyBytes),
		Status:   respStatus,
		Response: respJson,
		Ms:       costMs,
	})

	return resp, err
}

// headersToMap 将http.Header转换为map[string]string
func headersToMap(header http.Header) map[string]string {
	result := make(map[string]string, len(header))
	for k, v := range header {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}
