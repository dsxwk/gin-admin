package errcode

import (
	"fmt"
	"net/http"
	"strconv"
)

type ErrorCode struct {
	Code     int64  `json:"code"`     // 错误码
	Msg      string `json:"msg"`      // 错误描述
	Data     any    `json:"data"`     // 返回数据
	HttpCode int    `json:"HttpCode"` // http状态码
}

// Error 实现error接口
func (e ErrorCode) Error() string {
	return fmt.Sprintf("%d=>%s", e.Code, e.Msg)
}

// WithPrefix 设置错误码前缀
func (e ErrorCode) WithPrefix(prefix int64) ErrorCode {
	code := fmt.Sprintf("%d%d", prefix, e.Code)
	e.Code, _ = strconv.ParseInt(code, 10, 64)
	return e
}

func NewError(code int64, msg string) ErrorCode {
	return ErrorCode{
		Code:     code,
		Msg:      msg,
		HttpCode: http.StatusOK,
	}
}

func (e ErrorCode) WithCode(code int64) ErrorCode {
	e.Code = code
	return e
}

func (e ErrorCode) WithMsg(msg string) ErrorCode {
	e.Msg = msg
	return e
}

func (e ErrorCode) WithData(data any) ErrorCode {
	e.Data = data
	return e
}

func (e ErrorCode) WithHttpCode(httpCode int) ErrorCode {
	e.HttpCode = httpCode
	return e
}

func Success() ErrorCode {
	return NewError(0, "Success").WithHttpCode(http.StatusOK)
}

func Redirect() ErrorCode {
	return NewError(301, "Redirect").WithHttpCode(http.StatusMovedPermanently)
}

func ArgsError() ErrorCode {
	return NewError(400, "Bad Request").WithHttpCode(http.StatusBadRequest)
}

func Unauthorized() ErrorCode {
	return NewError(401, "Unauthorized").WithHttpCode(http.StatusUnauthorized)
}

func Forbidden() ErrorCode {
	return NewError(403, "Forbidden").WithHttpCode(http.StatusForbidden)
}

func NotFound() ErrorCode {
	return NewError(404, "Resource Not Found").WithHttpCode(http.StatusNotFound)
}

func RateLimitError() ErrorCode {
	return NewError(429, "Too Many Requests").WithHttpCode(http.StatusTooManyRequests)
}

func SystemError() ErrorCode {
	return NewError(500, "System Internal Error").WithHttpCode(http.StatusInternalServerError)
}

func TimeoutError() ErrorCode {
	return NewError(504, "Gateway Timeout").WithHttpCode(http.StatusGatewayTimeout)
}
