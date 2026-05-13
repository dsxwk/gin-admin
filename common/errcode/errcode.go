package errcode

import (
	"fmt"
	"github.com/samber/lo"
	"net/http"
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

func NewError(code int64, msg string, httpCode ...int) ErrorCode {
	httpCodeValue := lo.FirstOr(httpCode, http.StatusOK)
	if httpCodeValue == 0 {
		httpCodeValue = http.StatusOK
	}

	return ErrorCode{
		Code:     code,
		Msg:      msg,
		HttpCode: httpCodeValue,
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
	return ErrorCode{Code: 0, Msg: "Success", HttpCode: http.StatusOK}
}

func Redirect() ErrorCode {
	return ErrorCode{
		Code:     301,
		Msg:      "Redirect",
		HttpCode: http.StatusMovedPermanently,
	}
}

func ArgsError() ErrorCode {
	return ErrorCode{
		Code:     400,
		Msg:      "Invalid arguments",
		HttpCode: http.StatusBadRequest,
	}
}

func Unauthorized() ErrorCode {
	return ErrorCode{
		Code:     401,
		Msg:      "Unauthorized",
		HttpCode: http.StatusUnauthorized,
	}
}

func NotFound() ErrorCode {
	return ErrorCode{
		Code:     404,
		Msg:      "Resource not found",
		HttpCode: http.StatusNotFound,
	}
}

func RateLimitError() ErrorCode {
	return ErrorCode{
		Code:     429,
		Msg:      "Rate limit exceeded",
		HttpCode: http.StatusTooManyRequests,
	}
}

func SystemError() ErrorCode {
	return ErrorCode{
		Code:     500,
		Msg:      "Internal server error",
		HttpCode: http.StatusInternalServerError,
	}
}

func TimeoutError() ErrorCode {
	return ErrorCode{
		Code:     504,
		Msg:      "Request Timeout",
		HttpCode: http.StatusGatewayTimeout,
	}
}
