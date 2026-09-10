package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gin/app/model"
	"gin/common/ctxkey"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/message"
	t "gin/pkg/time"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultEsTimeout = 10 * time.Second

// Client Elasticsearch客户端
type Client struct {
	baseURL    string        // 基础地址
	username   string        // 用户名
	password   string        // 密码
	apiKey     string        // APIKey
	timeout    time.Duration // 请求超时
	httpClient *http.Client  // HTTP客户端
}

// NewClient 创建ES客户端
func NewClient(conf config.Es) *Client {
	baseURL := "http://127.0.0.1:9200"
	if len(conf.Addresses) > 0 && conf.Addresses[0] != "" {
		baseURL = conf.Addresses[0]
	}
	if !strings.Contains(baseURL, "://") {
		baseURL = "http://" + baseURL
	}

	timeout := conf.Timeout
	if timeout <= 0 {
		timeout = defaultEsTimeout
	}

	return &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		username: conf.Username,
		password: conf.Password,
		apiKey:   conf.ApiKey,
		timeout:  timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// FieldMapping 获取字段类型映射
func FieldMapping(fieldType string) map[string]any {
	switch strings.ToLower(fieldType) {
	case "keyword":
		return map[string]any{"type": "keyword"}
	case "text":
		return map[string]any{"type": "text"}
	case "long":
		return map[string]any{"type": "long"}
	case "integer":
		return map[string]any{"type": "integer"}
	case "short":
		return map[string]any{"type": "short"}
	case "byte":
		return map[string]any{"type": "byte"}
	case "double":
		return map[string]any{"type": "double"}
	case "float":
		return map[string]any{"type": "float"}
	case "boolean":
		return map[string]any{"type": "boolean"}
	case "date":
		return map[string]any{"type": "date"}
	case "object":
		return map[string]any{"type": "object"}
	case "nested":
		return map[string]any{"type": "nested"}
	default:
		return map[string]any{"type": "keyword"}
	}
}

// Mappings 批量生成字段映射
func Mappings(fields map[string]string) map[string]any {
	result := make(map[string]any, len(fields))
	for field, fieldType := range fields {
		result[field] = FieldMapping(fieldType)
	}
	return result
}

// Ping 检查ES连接
func (c *Client) Ping(ctx context.Context) error {
	_, _, err := c.execute(ctx, "ping", "", http.MethodGet, "/", "application/json", nil)
	return err
}

// IndexExists 判断索引是否存在
func (c *Client) IndexExists(ctx context.Context, index string) (bool, error) {
	status, _, err := c.request(ctx, "exists", index, http.MethodHead, "/"+url.PathEscape(index), "", nil)
	if err != nil {
		return false, err
	}
	return status == http.StatusOK, nil
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, index string, fieldTypes map[string]string) error {
	exists, err := c.IndexExists(ctx, index)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if fieldTypes == nil {
		fieldTypes = map[string]string{}
	}
	body := map[string]any{
		"mappings": map[string]any{
			"properties": Mappings(fieldTypes),
		},
	}
	return c.CreateIndexBody(ctx, index, body)
}

// CreateIndexBody 使用原始请求体创建索引
func (c *Client) CreateIndexBody(ctx context.Context, index string, body any) error {
	_, _, err := c.execute(ctx, "create_index", index, http.MethodPut, "/"+url.PathEscape(index), "application/json", body)
	return err
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(ctx context.Context, index string) error {
	_, _, err := c.execute(ctx, "delete_index", index, http.MethodDelete, "/"+url.PathEscape(index), "application/json", nil)
	return err
}

// Index 新增或覆盖文档
func (c *Client) Index(ctx context.Context, index, id string, doc any) error {
	path := "/" + url.PathEscape(index) + "/_doc/" + url.PathEscape(id) + "?refresh=true"
	_, _, err := c.execute(ctx, "index", index, http.MethodPut, path, "application/json", doc)
	return err
}

// Update 更新文档
func (c *Client) Update(ctx context.Context, index, id string, doc any) error {
	path := "/" + url.PathEscape(index) + "/_update/" + url.PathEscape(id) + "?refresh=true"
	_, _, err := c.execute(ctx, "update", index, http.MethodPost, path, "application/json", map[string]any{
		"doc": doc,
	})
	return err
}

// Delete 删除文档
func (c *Client) Delete(ctx context.Context, index, id string) error {
	path := "/" + url.PathEscape(index) + "/_doc/" + url.PathEscape(id) + "?refresh=true"
	_, _, err := c.execute(ctx, "delete", index, http.MethodDelete, path, "application/json", nil)
	return err
}

// GetDocument 获取文档
func (c *Client) GetDocument[T any](ctx context.Context, index, id string) (*DocumentResponse[T], error) {
	path := "/" + url.PathEscape(index) + "/_doc/" + url.PathEscape(id)
	data, err := c.get(ctx, "detail", index, path)
	if err != nil {
		return nil, err
	}

	var result DocumentResponse[T]
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("解析ES文档响应失败: %w", err)
	}
	return &result, nil
}

// Search 搜索文档
func (c *Client) Search[T any](ctx context.Context, index string, body map[string]any) (*SearchResponse[T], error) {
	if body == nil {
		body = map[string]any{}
	}
	path := "/" + url.PathEscape(index) + "/_search"
	data, err := c.postBody(ctx, "search", index, path, body)
	if err != nil {
		return nil, err
	}

	var result SearchResponse[T]
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("解析ES搜索响应失败: %w", err)
	}
	return &result, nil
}

// DeleteByQuery 按条件删除文档
func (c *Client) DeleteByQuery(ctx context.Context, index string, query map[string]any) error {
	if query == nil {
		query = map[string]any{}
	}
	path := "/" + url.PathEscape(index) + "/_delete_by_query?refresh=true"
	_, _, err := c.execute(ctx, "delete_by_query", index, http.MethodPost, path, "application/json", map[string]any{
		"query": query,
	})
	return err
}

// Refresh 刷新索引
func (c *Client) Refresh(ctx context.Context, index string) error {
	_, _, err := c.execute(ctx, "refresh", index, http.MethodPost, "/"+url.PathEscape(index)+"/_refresh", "application/json", nil)
	return err
}

// Bulk 批量写入文档
func (c *Client) Bulk(ctx context.Context, index, actions string) error {
	path := "/_bulk"
	if index != "" {
		path = "/" + url.PathEscape(index) + "/_bulk"
	}
	_, _, err := c.execute(ctx, "bulk", index, http.MethodPost, path, "application/x-ndjson", actions)
	return err
}

// JsonValue JSON字段内容
func (c *Client) JsonValue(v *model.JsonValue) any {
	if v == nil {
		return nil
	}
	return v.Data
}

// FormatDateTime 格式化时间
func (c *Client) FormatDateTime(date *model.DateTime) string {
	if date == nil {
		return ""
	}
	return t.FromTime(time.Time(*date)).Format("Y-m-dTH:i:sP")
}

// get GET请求并返回响应体
func (c *Client) get(ctx context.Context, action, index, path string) ([]byte, error) {
	_, data, err := c.execute(ctx, action, index, http.MethodGet, path, "application/json", nil)
	return data, err
}

// postBody POST请求并返回响应体
func (c *Client) postBody(ctx context.Context, action, index, path string, body any) ([]byte, error) {
	_, data, err := c.execute(ctx, action, index, http.MethodPost, path, "application/json", body)
	return data, err
}

// execute 执行请求并处理HTTP错误
func (c *Client) execute(ctx context.Context, action, index, method, path, contentType string, body any) (int, []byte, error) {
	status, data, err := c.request(ctx, action, index, method, path, contentType, body)
	if err != nil {
		return status, data, err
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return status, data, fmt.Errorf("es请求失败,HTTP状态码:%d,响应:%s", status, strings.TrimSpace(string(data)))
	}
	return status, data, nil
}

// request 执行HTTP请求并记录调试信息
func (c *Client) request(ctx context.Context, action, index, method, path, contentType string, body any) (int, []byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	fullURL := c.baseURL + "/" + strings.TrimLeft(path, "/")
	var (
		reader  io.Reader
		reqData []byte
		err     error
	)
	if body != nil {
		switch v := body.(type) {
		case []byte:
			reqData = v
		case string:
			reqData = []byte(v)
		default:
			reqData, err = json.Marshal(body)
			if err != nil {
				return 0, nil, fmt.Errorf("序列化ES请求失败: %w", err)
			}
		}
		reader = bytes.NewReader(reqData)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reader)
	if err != nil {
		return 0, nil, fmt.Errorf("创建ES请求失败: %w", err)
	}

	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	} else if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	start := time.Now()
	resp, requestErr := c.httpClient.Do(req)
	var (
		status int
		data   []byte
	)
	if requestErr == nil {
		status = resp.StatusCode
		data, requestErr = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}

	costMs := float64(time.Since(start).Nanoseconds()) / 1e6
	publishTrace(ctx, action, index, reqData, status, data, requestErr, costMs)
	return status, data, requestErr
}

// publishTrace 发布ES调试事件
func publishTrace(ctx context.Context, action, index string, reqData []byte, status int, respData []byte, reqErr error, ms float64) {
	var request, response any
	if len(reqData) > 0 {
		if err := json.Unmarshal(reqData, &request); err != nil {
			request = string(reqData)
		}
	}
	if len(respData) > 0 {
		if err := json.Unmarshal(respData, &response); err != nil {
			response = string(respData)
		}
	}

	code := "error"
	if reqErr == nil {
		code = "OK"
		if status < http.StatusOK || status >= http.StatusMultipleChoices {
			code = strconv.Itoa(status)
		}
	}

	traceId := "unknown"
	if ctx != nil {
		if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
			if s, ok := id.(string); ok && s != "" {
				traceId = s
			}
		}
	}

	message.NewEvent().Publish(debugger.TopicEs, debugger.EsEvent{
		TraceId:  traceId,
		Action:   action,
		Index:    index,
		Request:  request,
		Response: response,
		Code:     code,
		Ms:       ms,
	})
}
