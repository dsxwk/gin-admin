package facade

import (
	"context"
	"gin/pkg/http"
)

// Http HTTP客户端门面
var Http = &httpFacade{}

type httpFacade struct{}

//// SendToJson 发送HTTP请求并解析JSON响应
//func SendToJson[T any](ctx context.Context, method, uri string, opt *http.Option) (*T, error) {
//	return http.SendToJson[T](ctx, method, uri, opt)
//}

// Send 发送HTTP请求
func (h *httpFacade) Send(ctx context.Context, method, uri string, opt *http.Option) ([]byte, error) {
	return http.Send(ctx, method, uri, opt)
}

// HttpAs 泛型HTTP客户端门面
func HttpAs[T any]() *HttpFacade[T] {
	return &HttpFacade[T]{}
}

type HttpFacade[T any] struct{}

// Send 发送HTTP请求并解析为T类型
func (h *HttpFacade[T]) Send(ctx context.Context, method, uri string, opt *http.Option) (*T, error) {
	return http.SendToJson[T](ctx, method, uri, opt)
}
