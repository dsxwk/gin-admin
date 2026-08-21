package facade

import (
	"context"
	"gin/pkg/serviceprovider/http"
)

// Http 门面函数,返回门面实例
// 使用示例:
//
//	user, err := facade.Http().SendAsJson[UserResponse](ctx, "https://api.example.com/user/1", nil)
//	data, err := facade.Http().Get(ctx, "https://api.example.com/user/1", nil)
func Http() HttpFacade {
	return HttpFacade{}
}

type HttpFacade struct{}

// Send 发送HTTP请求
func (h HttpFacade) Send(ctx context.Context, method, uri string, opt *http.Option) ([]byte, error) {
	return http.NewClient().Send(ctx, method, uri, opt)
}

// Get 发送GET请求
func (h HttpFacade) Get(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.NewClient().Get(ctx, uri, opt)
}

// Post 发送POST请求
func (h HttpFacade) Post(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.NewClient().Post(ctx, uri, opt)
}

// Put 发送PUT请求
func (h HttpFacade) Put(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.NewClient().Put(ctx, uri, opt)
}

// Delete 发送DELETE请求
func (h HttpFacade) Delete(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.NewClient().Delete(ctx, uri, opt)
}

// AsJson 将响应体解析为T类型
func (h HttpFacade) AsJson[T any](data []byte) (*T, error) {
	return http.NewClient().AsJson[T](data)
}

// SendAsJson 发送请求并解析为T类型
func (h HttpFacade) SendAsJson[T any](ctx context.Context, method, uri string, opt *http.Option) (*T, error) {
	return http.NewClient().SendAsJson[T](ctx, method, uri, opt)
}
