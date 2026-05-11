package http

import (
	"bytes"
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/pkg/provider/debugger"
	"gin/pkg/provider/message"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const timeout = 5 * time.Second

var (
	defaultClient *fasthttp.Client
	once          sync.Once
)

// InitClient 初始化全局HTTP客户端
func InitClient() {
	once.Do(func() {
		defaultClient = &fasthttp.Client{
			MaxConnsPerHost: 100,
			ReadTimeout:     timeout,
			WriteTimeout:    timeout,
		}
	})
}

// GetClient 获取全局HTTP客户端
func GetClient() *fasthttp.Client {
	if defaultClient == nil {
		InitClient()
	}
	return defaultClient
}

// Client HTTP客户端
type Client[T any] struct {
	client  *fasthttp.Client
	timeout time.Duration
}

// NewClient 创建HTTP客户端
func NewClient[T any]() *Client[T] {
	return &Client[T]{
		timeout: timeout,
		client:  GetClient(),
	}
}

// WithTimeout 自定义超时的HTTP客户端
func (c *Client[T]) WithTimeout(timeout time.Duration) *Client[T] {
	return &Client[T]{
		timeout: timeout,
		client:  c.client,
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
	start := time.Now()
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	if opt == nil {
		opt = &Option{}
	}
	requestTimeout := opt.Timeout
	if requestTimeout == 0 {
		requestTimeout = timeout
	}

	method = strings.ToUpper(method)
	uri = c.buildUrl(uri, opt.Query)
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)

	// 判断是否有文件上传
	if opt.Files != nil && len(opt.Files) > 0 {
		return c.doFileUpload(ctx, req, resp, uri, opt, requestTimeout, start)
	}

	// 构建请求体
	contentType := ""
	if opt.Form != nil && len(opt.Form) > 0 {
		values := url.Values{}
		for k, v := range opt.Form {
			values.Set(k, fmt.Sprintf("%v", v))
		}
		req.SetBodyString(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	if opt.Body != nil {
		method = strings.ToUpper(method)
		if method != http.MethodGet && method != http.MethodHead {
			switch v := opt.Body.(type) {
			case []byte:
				req.SetBody(v)
				contentType = "application/octet-stream"
			case string:
				req.SetBodyString(v)
				contentType = "text/plain"
			case *bytes.Buffer:
				req.SetBody(v.Bytes())
				contentType = "application/octet-stream"
			case *strings.Reader:
				data, err := io.ReadAll(v)
				if err != nil {
					return nil, fmt.Errorf("读取Body失败: %w", err)
				}
				req.SetBody(data)
				contentType = "text/plain"
			default:
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("JSON序列化失败: %w", err)
				}
				req.SetBody(jsonBytes)
				contentType = "application/json"
			}
		}
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

	c.client.ReadTimeout = requestTimeout
	c.client.WriteTimeout = requestTimeout

	var (
		respJson interface{}
		status   int
	)
	defer func() {
		cost := time.Since(start)
		costMs := float64(cost.Nanoseconds()) / 1e6

		traceId := "unknown"
		if ctx != nil {
			if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
				if s, ok := id.(string); ok && s != "" {
					traceId = s
				}
			}
		}

		message.NewEvent().Publish(debugger.TopicHttp, debugger.HttpEvent{
			TraceId:  traceId,
			Url:      uri,
			Method:   method,
			Header:   opt.Headers,
			Body:     string(req.Body()),
			Status:   status,
			Response: respJson,
			Ms:       costMs,
		})
	}()

	if err := c.client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	respBody := resp.Body()
	err := json.Unmarshal(respBody, &respJson)
	if err != nil {
		respJson = respBody
	}
	status = resp.StatusCode()
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.Header.StatusCode(), respBody)
	}

	return respBody, nil
}

// doFileUpload 文件上传内部函数
func (c *Client[T]) doFileUpload(ctx context.Context, req *fasthttp.Request, resp *fasthttp.Response, uri string, opt *Option, requestTimeout time.Duration, start time.Time) ([]byte, error) {
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

	req.Header.SetMethod("POST")
	req.SetRequestURI(uri)
	req.SetBody(body.Bytes())
	req.Header.Set("Content-Type", contentType)

	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	c.client.ReadTimeout = requestTimeout
	c.client.WriteTimeout = requestTimeout

	var (
		respJson interface{}
		status   int
	)
	defer func() {
		cost := time.Since(start)
		costMs := float64(cost.Nanoseconds()) / 1e6

		traceId := "unknown"
		if ctx != nil {
			if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
				if s, ok := id.(string); ok && s != "" {
					traceId = s
				}
			}
		}

		message.NewEvent().Publish(debugger.TopicHttp, debugger.HttpEvent{
			TraceId:  traceId,
			Url:      uri,
			Method:   "POST",
			Header:   opt.Headers,
			Body:     string(body.Bytes()),
			Status:   status,
			Response: respJson,
			Ms:       costMs,
		})
	}()

	if err := c.client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	respBody := resp.Body()
	status = resp.StatusCode()

	err := json.Unmarshal(respBody, &respJson)
	if err != nil {
		respJson = respBody
	}

	if status != fasthttp.StatusOK {
		return respBody, fmt.Errorf("上传失败, 状态码: %d, 响应: %s", status, string(respBody))
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
