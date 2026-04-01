package facade

import (
	"context"
	"gin/pkg/http"
)

// Http HTTP客户端门面
var Http = &httpFacade{
	Json: &jsonHttpFacadeHelper{},
}

type httpFacade struct {
	Json *jsonHttpFacadeHelper
}

type jsonHttpFacadeHelper struct{}

// RequestJson 泛型函数,通过helper调用
func RequestJson[T any](ctx context.Context, method, uri string, opt *http.Option) (*T, error) {
	return http.RequestJson[T](ctx, method, uri, opt)
}

// Request 发送HTTP请求
func (h *httpFacade) Request(ctx context.Context, method, uri string, opt *http.Option) (string, error) {
	return http.Request(ctx, method, uri, opt)
}
