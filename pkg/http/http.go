package http

import (
	"bytes"
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/pkg/debugger"
	"gin/pkg/message"
	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const timeout = 5 * time.Second

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

// buildURL 拼接get请求query参数
func buildUrl(baseURL string, query map[string]interface{}) string {
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

// Send 发送http请求
func Send(ctx context.Context, method, uri string, opt *Option) ([]byte, error) {
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
	uri = buildUrl(uri, opt.Query)
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)

	// 判断是否有文件上传
	if opt.Files != nil && len(opt.Files) > 0 {
		return fileUpload(ctx, req, resp, uri, opt, requestTimeout, start)
	}

	// 构建请求体
	contentType := ""
	if opt.Form != nil && len(opt.Form) > 0 {
		// x-www-form-urlencoded
		values := url.Values{}
		for k, v := range opt.Form {
			values.Set(k, fmt.Sprintf("%v", v))
		}

		req.SetBodyString(values.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	if opt.Body != nil {
		// 处理Body参数
		method = strings.ToUpper(method)
		// GET和HEAD请求不应该有body
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
				// 读取strings.Reader的内容
				data, err := io.ReadAll(v)
				if err != nil {
					return nil, fmt.Errorf("读取Body失败: %w", err)
				}
				req.SetBody(data)
				contentType = "text/plain"

			default:
				// 尝试JSON序列化
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
	if contentType != "" && opt.Headers["Content-Type"] == "" {
		opt.Headers["Content-Type"] = contentType
	}
	if contentType == "" && opt.Headers["Content-Type"] == "" {
		opt.Headers["Content-Type"] = "application/json"
	}

	// 设置自定义headers(会覆盖默认设置)
	if opt.Headers != nil {
		for k, v := range opt.Headers {
			req.Header.Set(k, v)
		}
	}

	// 创建客户端并设置超时
	client := &fasthttp.Client{
		MaxConnsPerHost: 100,            // 最大连接数
		ReadTimeout:     requestTimeout, // 读取超时
		WriteTimeout:    requestTimeout, // 写入超时
	}

	var (
		respJson interface{}
		status   int
	)
	defer func() {
		// 耗时
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

		message.GetEventBus().Publish(debugger.TopicHttp, debugger.HttpEvent{
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

	// 发送请求
	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	// 获取响应内容
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

// fileUpload 处理multipart/form-data上传
func fileUpload(ctx context.Context, req *fasthttp.Request, resp *fasthttp.Response, uri string, opt *Option, requestTimeout time.Duration, start time.Time) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
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

		// 使用自定义字段名或默认字段名
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

	// 添加普通表单字段
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

	client := &fasthttp.Client{
		MaxConnsPerHost: 100,
		ReadTimeout:     requestTimeout,
		WriteTimeout:    requestTimeout,
	}

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

		message.GetEventBus().Publish(debugger.TopicHttp, debugger.HttpEvent{
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

	if err := client.Do(req, resp); err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	respBody := resp.Body()
	status = resp.StatusCode()

	// 获取响应内容
	err := json.Unmarshal(respBody, &respJson)
	if err != nil {
		respJson = respBody
	}

	if status != fasthttp.StatusOK {
		return respBody, fmt.Errorf("上传失败, 状态码: %d, 响应: %s", status, string(respBody))
	}

	return respBody, nil
}

// SendToJson 发送请求并解析json响应
func SendToJson[T any](ctx context.Context, method, url string, opt *Option) (*T, error) {
	resp, err := Send(ctx, method, url, opt)
	if err != nil {
		return nil, err
	}

	var result T
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("json解析失败: %w\n响应内容:\n%s", err, resp)
	}

	return &result, nil
}
