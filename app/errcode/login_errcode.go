package errcode

import (
	"gin/pkg/errcode"
	"net/http"
)

const (
	LoginErrCodePrefix = 100 // 错误码前缀
)

type LoginErrCode struct{}

func (s LoginErrCode) AccountError() errcode.ErrorCode {
	return errcode.NewError(101, "login.accountErr").
		WithHttpCode(http.StatusNotFound).
		WithPrefix(LoginErrCodePrefix)
}

func (s LoginErrCode) PwdError() errcode.ErrorCode {
	return errcode.NewError(102, "login.pwdErr").
		WithHttpCode(http.StatusInternalServerError).
		WithPrefix(LoginErrCodePrefix)
}

func (s LoginErrCode) AccountDisabled() errcode.ErrorCode {
	return errcode.NewError(102, "login.accountDisabled").
		WithHttpCode(http.StatusInternalServerError).
		WithPrefix(LoginErrCodePrefix)
}

func (s LoginErrCode) InvalidToken() errcode.ErrorCode {
	return errcode.NewError(103, "login.invalidToken").
		WithHttpCode(http.StatusUnauthorized).
		WithPrefix(LoginErrCodePrefix)
}
