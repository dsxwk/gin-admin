package {{.Package}}

import (
	"gin/pkg/errcode"
	"net/http"
)

const (
	{{.PrefixName}} = {{.Prefix}} // 错误码前缀
)

// {{.Name}} 业务错误码
type {{.Name}} struct{}

// ExampleError 示例错误
func (s {{.Name}}) ExampleError() errcode.ErrorCode {
	return errcode.NewError(1, "{{.SnakeName}}.exampleError").
		WithHttpCode(http.StatusBadRequest).
		WithPrefix({{.PrefixName}})
}
