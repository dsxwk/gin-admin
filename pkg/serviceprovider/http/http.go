package http

import (
	"bytes"
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

const defaultTimeout = 5 * time.Second

var (
	defaultClient     *http.Client
	defaultClientOnce sync.Once
)

type traceSkipKey struct{}

// WithoutTrace 跳过当前请求的HTTP调试记录
func WithoutTrace(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceSkipKey{}, true)
}

// skipTrace 判断当前请求是否跳过HTTP调试记录
func skipTrace(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	skip, _ := ctx.Value(traceSkipKey{}).(bool)
	return skip
}

// GetClient 获取全局HTTP客户端
func GetClient() *http.Client {
	defaultClientOnce.Do(func() {
		transport, ok := http.DefaultTransport.(*http.Transport)
		if ok {
			transport = transport.Clone()
		} else {
			transport = &http.Transport{}
		}
		transport.MaxConnsPerHost = 100
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 100
		transport.IdleConnTimeout = 90 * time.Second

		defaultClient = &http.Client{
			Transport: &TracingTransport{
				Transport: transport,
			},
		}
	})
	return defaultClient
}

// NewClient 创建HTTP客户端
func NewClient() *Client {
	return &Client{
		option: Option{
			Timeout: defaultTimeout,
		},
	}
}

// WithTimeout 自定义超时的HTTP客户端
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	client := c.clone()
	client.option.Timeout = timeout
	return client
}

// WithHeader 设置单个请求头
func (c *Client) WithHeader(name, value string) *Client {
	client := c.clone()
	if client.option.Headers == nil {
		client.option.Headers = make(map[string]string)
	}
	client.option.Headers[name] = value
	return client
}

// WithHeaders 批量设置请求头
func (c *Client) WithHeaders(headers map[string]string) *Client {
	client := c.clone()
	if client.option.Headers == nil {
		client.option.Headers = make(map[string]string, len(headers))
	}
	maps.Copy(client.option.Headers, headers)
	return client
}

// WithQuery 设置query参数
func (c *Client) WithQuery(query map[string]any) *Client {
	client := c.clone()
	if client.option.Query == nil {
		client.option.Query = make(map[string]any, len(query))
	}
	maps.Copy(client.option.Query, query)
	return client
}

// WithForm 设置表单参数
func (c *Client) WithForm(form map[string]any) *Client {
	client := c.clone()
	if client.option.Form == nil {
		client.option.Form = make(map[string]any, len(form))
	}
	maps.Copy(client.option.Form, form)
	client.option.Body = nil
	return client
}

// WithBody 设置请求体
func (c *Client) WithBody(body any) *Client {
	client := c.clone()
	client.option.Body = body
	client.option.Form = nil
	return client
}

// WithFiles 批量设置上传文件
func (c *Client) WithFiles(files map[string]File) *Client {
	client := c.clone()
	if client.option.Files == nil {
		client.option.Files = make(map[string]File, len(files))
	}
	maps.Copy(client.option.Files, files)
	return client
}

// WithFile 设置单个上传文件
func (c *Client) WithFile(field string, file File) *Client {
	client := c.clone()
	if client.option.Files == nil {
		client.option.Files = make(map[string]File)
	}
	client.option.Files[field] = file
	return client
}

// File 文件
type File struct {
	FileName  string // 文件名
	FilePath  string // 文件路径
	FileData  []byte // 文件数据
	FieldName string // 表单字段名
}

type Option struct {
	Headers map[string]string // 请求头
	Query   map[string]any    // query参数
	Form    map[string]any    // form参数
	Body    any               // 请求体
	Files   map[string]File   // 文件上传字段
	Timeout time.Duration     // 超时时间
}

type Client struct {
	option Option // 请求配置
}

// Response HTTP响应
type Response struct {
	StatusCode int    // 状态码
	Body       []byte // 响应体
	ErrMsg     error  // 错误信息
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, uri string) Response {
	return c.Send(ctx, "GET", uri)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, uri string) Response {
	return c.Send(ctx, "POST", uri)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, uri string) Response {
	return c.Send(ctx, "PUT", uri)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, uri string) Response {
	return c.Send(ctx, "DELETE", uri)
}

// SendAsJson 发送请求并解析为T类型
func (c *Client) SendAsJson[T any](ctx context.Context, method, uri string) (*T, Response) {
	response := c.Send(ctx, method, uri)
	if response.ErrMsg != nil {
		return nil, response
	}

	var result T
	if len(response.Body) == 0 {
		return &result, response
	}
	if err := json.Unmarshal(response.Body, &result); err != nil {
		response.ErrMsg = fmt.Errorf("json解析失败: %w\n响应内容:\n%s", err, response.Body)
		return nil, response
	}
	return &result, response
}

// Send 发送HTTP请求
func (c *Client) Send(ctx context.Context, method, uri string) Response {
	requestTimeout := c.option.Timeout
	if requestTimeout <= 0 {
		requestTimeout = defaultTimeout
	}

	method = strings.ToUpper(method)
	uri = c.buildURL(uri)

	// 判断是否有文件上传
	if len(c.option.Files) > 0 {
		return c.doFileUpload(ctx, uri, requestTimeout)
	}

	// 构建请求体
	var bodyReader io.Reader
	contentType := ""
	if len(c.option.Form) > 0 {
		values := url.Values{}
		for k, v := range c.option.Form {
			values.Set(k, fmt.Sprintf("%v", v))
		}
		bodyReader = strings.NewReader(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	} else if c.option.Body != nil {
		switch v := c.option.Body.(type) {
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
		case io.Reader:
			bodyReader = v
			contentType = "application/octet-stream"
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return Response{ErrMsg: fmt.Errorf("JSON序列化失败: %w", err)}
			}
			bodyReader = bytes.NewReader(jsonBytes)
			contentType = "application/json"
		}
	}

	// 通过上下文应用超时
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, uri, bodyReader)
	if err != nil {
		return Response{ErrMsg: fmt.Errorf("创建请求失败: %w", err)}
	}

	// 设置自定义headers
	for k, v := range c.option.Headers {
		req.Header.Set(k, v)
	}
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 发送请求
	resp, err := GetClient().Do(req)
	if err != nil {
		return Response{ErrMsg: fmt.Errorf("请求失败: %w", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			StatusCode: resp.StatusCode,
			ErrMsg:     fmt.Errorf("读取响应失败: %w", err),
		}
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}
}

// clone 复制HTTP客户端
func (c *Client) clone() *Client {
	if c == nil {
		return NewClient()
	}

	return &Client{
		option: Option{
			Headers: maps.Clone(c.option.Headers),
			Query:   maps.Clone(c.option.Query),
			Form:    maps.Clone(c.option.Form),
			Body:    c.option.Body,
			Files:   maps.Clone(c.option.Files),
			Timeout: c.option.Timeout,
		},
	}
}

// doFileUpload 文件上传
func (c *Client) doFileUpload(ctx context.Context, uri string, requestTimeout time.Duration) Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, file := range c.option.Files {
		var fileData []byte
		var fileName string

		if file.FilePath != "" {
			data, err := os.ReadFile(file.FilePath)
			if err != nil {
				return Response{ErrMsg: fmt.Errorf("读取文件 %s 失败: %w", file.FilePath, err)}
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
			return Response{ErrMsg: fmt.Errorf("创建表单文件失败: %w", err)}
		}

		if _, err = part.Write(fileData); err != nil {
			return Response{ErrMsg: fmt.Errorf("写入文件数据失败: %w", err)}
		}
	}

	for k, v := range c.option.Form {
		if err := writer.WriteField(k, fmt.Sprintf("%v", v)); err != nil {
			return Response{ErrMsg: fmt.Errorf("写入表单字段失败: %w", err)}
		}
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return Response{ErrMsg: fmt.Errorf("关闭writer失败: %w", err)}
	}

	// 通过上下文应用超时
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", uri, body)
	if err != nil {
		return Response{ErrMsg: fmt.Errorf("创建请求失败: %w", err)}
	}

	for k, v := range c.option.Headers {
		if strings.EqualFold(k, "Content-Type") {
			continue
		}
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)

	// 发送请求
	resp, err := GetClient().Do(req)
	if err != nil {
		return Response{ErrMsg: fmt.Errorf("请求失败: %w", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			StatusCode: resp.StatusCode,
			ErrMsg:     fmt.Errorf("读取响应失败: %w", err),
		}
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}
}

// buildURL 拼接get请求query参数
func (c *Client) buildURL(baseURL string) string {
	if len(c.option.Query) == 0 {
		return baseURL
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := parsedURL.Query()
	for k, v := range c.option.Query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	parsedURL.RawQuery = q.Encode()
	return parsedURL.String()
}

// TracingTransport 传输中间件追踪
type TracingTransport struct {
	Transport http.RoundTripper
}

// RoundTrip 实现http.RoundTripper
func (t *TracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	if skipTrace(req.Context()) {
		return transport.RoundTrip(req)
	}

	start := time.Now()

	// 读取并缓存请求体用于tracing
	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(reqBodyBytes))
	}

	resp, err := transport.RoundTrip(req)

	// 收集响应信息用于tracing
	var (
		respBodyBytes []byte
		respStatus    int
	)
	if err == nil {
		respStatus = resp.StatusCode
		respBodyBytes, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
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

	var respJson any
	if len(respBodyBytes) > 0 {
		if unmarshalErr := json.Unmarshal(respBodyBytes, &respJson); unmarshalErr != nil {
			respJson = respBodyBytes
		}
	}

	eventbus.NewBus().Publish(debugger.TopicHTTP, debugger.HTTPEvent{
		TraceID:  traceId,
		URL:      req.URL.String(),
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
