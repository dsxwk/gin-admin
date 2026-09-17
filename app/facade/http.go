package facade

import (
	"context"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/http"
	"time"
)

// Http 门面函数,返回门面实例
// 使用示例:
//
//	user, response := facade.Http().SendAsJson[UserResponse](ctx, "https://api.example.com/user/1")
//	response := facade.Http().Get(ctx, "https://api.example.com/user/1")
//	response := facade.Http().WithHeader("Authorization", "Bearer token").WithBody(data).Post(ctx, "https://api.example.com/user")
func Http() HttpFacade {
	return HttpFacade{
		client: container.Default().Get[*http.Client](serviceprovider.ServiceHTTP),
	}
}

type HttpFacade struct {
	client *http.Client
}

// instance 获取HTTP客户端实例
func (h HttpFacade) instance() *http.Client {
	if h.client == nil {
		return http.NewClient()
	}
	return h.client
}

// WithTimeout 自定义超时时间
func (h HttpFacade) WithTimeout(timeout time.Duration) HttpFacade {
	h.client = h.instance().WithTimeout(timeout)
	return h
}

// WithHeader 设置单个请求头
func (h HttpFacade) WithHeader(name, value string) HttpFacade {
	h.client = h.instance().WithHeader(name, value)
	return h
}

// WithHeaders 批量设置请求头
func (h HttpFacade) WithHeaders(headers map[string]string) HttpFacade {
	h.client = h.instance().WithHeaders(headers)
	return h
}

// WithQuery 设置query参数
func (h HttpFacade) WithQuery(query map[string]any) HttpFacade {
	h.client = h.instance().WithQuery(query)
	return h
}

// WithForm 设置表单参数
func (h HttpFacade) WithForm(form map[string]any) HttpFacade {
	h.client = h.instance().WithForm(form)
	return h
}

// WithBody 设置请求体
func (h HttpFacade) WithBody(body any) HttpFacade {
	h.client = h.instance().WithBody(body)
	return h
}

// WithFiles 批量设置上传文件
func (h HttpFacade) WithFiles(files map[string]http.File) HttpFacade {
	h.client = h.instance().WithFiles(files)
	return h
}

// WithFile 设置单个上传文件
func (h HttpFacade) WithFile(field string, file http.File) HttpFacade {
	h.client = h.instance().WithFile(field, file)
	return h
}

// Send 发送HTTP请求
func (h HttpFacade) Send(ctx context.Context, method, uri string) http.Response {
	return h.instance().Send(ctx, method, uri)
}

// Get 发送GET请求
func (h HttpFacade) Get(ctx context.Context, uri string) http.Response {
	return h.instance().Get(ctx, uri)
}

// Post 发送POST请求
func (h HttpFacade) Post(ctx context.Context, uri string) http.Response {
	return h.instance().Post(ctx, uri)
}

// Put 发送PUT请求
func (h HttpFacade) Put(ctx context.Context, uri string) http.Response {
	return h.instance().Put(ctx, uri)
}

// Delete 发送DELETE请求
func (h HttpFacade) Delete(ctx context.Context, uri string) http.Response {
	return h.instance().Delete(ctx, uri)
}

// SendAsJson 发送请求并解析为T类型
func (h HttpFacade) SendAsJson[T any](ctx context.Context, method, uri string) (*T, http.Response) {
	return h.instance().SendAsJson[T](ctx, method, uri)
}
