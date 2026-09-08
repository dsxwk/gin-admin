package errcode

import (
	"gin/pkg/errcode"
	"net/http"
)

func Success() errcode.ErrorCode {
	return errcode.NewError(0, "Success").WithHttpCode(http.StatusOK)
}

func Redirect() errcode.ErrorCode {
	return errcode.NewError(301, "Redirect").WithHttpCode(http.StatusMovedPermanently)
}

func ArgsError() errcode.ErrorCode {
	return errcode.NewError(400, "Bad Request").WithHttpCode(http.StatusBadRequest)
}

func Unauthorized() errcode.ErrorCode {
	return errcode.NewError(401, "Unauthorized").WithHttpCode(http.StatusUnauthorized)
}

func Forbidden() errcode.ErrorCode {
	return errcode.NewError(403, "Forbidden").WithHttpCode(http.StatusForbidden)
}

func NotFound() errcode.ErrorCode {
	return errcode.NewError(404, "Resource Not Found").WithHttpCode(http.StatusNotFound)
}

func RateLimitError() errcode.ErrorCode {
	return errcode.NewError(429, "Too Many Requests").WithHttpCode(http.StatusTooManyRequests)
}

func SystemError() errcode.ErrorCode {
	return errcode.NewError(500, "System Internal Error").WithHttpCode(http.StatusInternalServerError)
}

func TimeoutError() errcode.ErrorCode {
	return errcode.NewError(504, "Gateway Timeout").WithHttpCode(http.StatusGatewayTimeout)
}
