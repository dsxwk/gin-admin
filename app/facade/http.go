package facade

import (
	"context"
	"gin/pkg/http"
)

// Http 泛型函数,返回对应类型的门面实例
// 使用示例:
//
//	user, err := facade.Http[UserResponse]().SendAsJson(ctx, "https://api.example.com/user/1", nil)
//	data, err := facade.Http[any]().Get(ctx, "https://api.example.com/user/1", nil)
func Http[T any]() HttpFacade[T] {
	return HttpFacade[T]{}
}

type HttpFacade[T any] struct{}

// Send 发送HTTP请求
func (h HttpFacade[T]) Send(ctx context.Context, method, uri string, opt *http.Option) ([]byte, error) {
	return http.Http[T]().Send(ctx, method, uri, opt)
}

// Get 发送GET请求
func (h HttpFacade[T]) Get(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.Http[T]().Get(ctx, uri, opt)
}

// Post 发送POST请求
func (h HttpFacade[T]) Post(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.Http[T]().Post(ctx, uri, opt)
}

// Put 发送PUT请求
func (h HttpFacade[T]) Put(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.Http[T]().Put(ctx, uri, opt)
}

// Delete 发送DELETE请求
func (h HttpFacade[T]) Delete(ctx context.Context, uri string, opt *http.Option) ([]byte, error) {
	return http.Http[T]().Delete(ctx, uri, opt)
}

// AsJson 将响应体解析为T类型
func (h HttpFacade[T]) AsJson(data []byte) (*T, error) {
	return http.Http[T]().AsJson(data)
}

// SendAsJson 发送请求并解析为T类型
func (h HttpFacade[T]) SendAsJson(ctx context.Context, method, uri string, opt *http.Option) (*T, error) {
	return http.Http[T]().SendAsJson(ctx, method, uri, opt)
}
