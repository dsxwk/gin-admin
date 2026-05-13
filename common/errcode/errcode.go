package errcode

import (
	"fmt"
	"strconv"
)

type ErrorCode struct {
	Code     int64       `json:"code"`     // 错误码
	Msg      string      `json:"msg"`      // 错误描述
	Data     interface{} `json:"data"`     // 返回数据
	HttpCode int         `json:"HttpCode"` // http状态码
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
		Code: code,
		Msg:  msg,
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

func (e ErrorCode) WithData(data interface{}) ErrorCode {
	e.Data = data
	return e
}

func (e ErrorCode) WithHttpCode(httpCode int) ErrorCode {
	e.HttpCode = httpCode
	return e
}

func Success() ErrorCode {
	return ErrorCode{Code: 0, Msg: "Success"}
}

func Redirect() ErrorCode {
	return ErrorCode{
		Code: 301,
		Msg:  "Redirect",
	}
}

func ArgsError() ErrorCode {
	return ErrorCode{
		Code: 400,
		Msg:  "Invalid arguments",
	}
}

func Unauthorized() ErrorCode {
	return ErrorCode{
		Code: 401,
		Msg:  "Unauthorized",
	}
}

func NotFound() ErrorCode {
	return ErrorCode{
		Code: 404,
		Msg:  "Resource not found",
	}
}

func RateLimitError() ErrorCode {
	return ErrorCode{
		Code: 429,
		Msg:  "Rate limit exceeded",
	}
}

func SystemError() ErrorCode {
	return ErrorCode{
		Code: 500,
		Msg:  "Internal server error",
	}
}

func TimeoutError() ErrorCode {
	return ErrorCode{
		Code: 504,
		Msg:  "Request Timeout",
	}
}
