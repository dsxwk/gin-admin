package errcode

import (
	"fmt"
	"strconv"
)

// 业务码
const (
	LoginErrorPrefix = 1000 // 登录错误码前缀
	UserErrorPrefix  = 1001 // 用户错误码前缀
	MenuErrorPrefix  = 1001 // 菜单错误码前缀
)

type ErrorCode struct {
	Code int64       `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回数据
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
