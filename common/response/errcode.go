package response

import (
	"errors"
	"gin/common/errcode"
	"gin/pkg/serviceprovider/lang"
	"gin/pkg/serviceprovider/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	log *logger.Logger
)

// Response 通用响应结构体
type Response struct {
	Code int64  `json:"code"` // 错误码
	Msg  string `json:"msg"`  // 提示信息
	Data any    `json:"data"` // 返回数据
}

func SetLogger(l *logger.Logger) {
	log = l
}

// json 输出Json响应
func (r Response) json(c *gin.Context, httpCode int) {
	c.Header("Content-Type", "application/json")
	c.JSON(httpCode, r)
	c.Abort()
}

// Success 返回成功响应,可传ErrorCode
func (r Response) Success(c *gin.Context, e error) {
	var (
		ec  errcode.ErrorCode
		ctx = c.Request.Context()
	)
	if e != nil && errors.As(e, &ec) {
		r.Code = ec.Code
		if ec.Msg != "" {
			r.Msg = ec.Msg
			r.Msg = lang.Trans(ctx, ec.Msg, nil)
		} else {
			r.Msg = errcode.Success().Msg
			r.Msg = lang.Trans(ctx, r.Msg, nil)
		}
		if ec.Data == nil {
			r.Data = []any{}
		} else {
			r.Data = ec.Data
		}
		if ec.HttpCode == 0 {
			ec.HttpCode = 200
		}

		r.json(c, ec.HttpCode)
		return
	}

	// 普通error
	r.Code = 0
	if e != nil {
		r.Msg = e.Error()
		r.Msg = lang.Trans(ctx, r.Msg, nil)
	} else {
		r.Msg = errcode.Success().Msg
		r.Msg = lang.Trans(ctx, r.Msg, nil)
	}
	r.Data = []any{}
	r.json(c, http.StatusOK)
}

// Error 返回失败响应,可传ErrorCode
func (r Response) Error(c *gin.Context, e error) {
	var (
		ec  errcode.ErrorCode
		ctx = c.Request.Context()
	)
	// ErrorCode类型
	if e != nil && errors.As(e, &ec) {
		r.Code = ec.Code
		if ec.Msg != "" {
			r.Msg = ec.Msg
			r.Msg = lang.Trans(ctx, r.Msg, nil)
			log.WithDebugger(ctx).Error(r.Msg)
		} else {
			r.Msg = errcode.Success().Msg
			r.Msg = lang.Trans(ctx, r.Msg, nil)
			log.WithDebugger(ctx).Error(r.Msg)
		}
		if ec.Data == nil {
			r.Data = []any{}
		} else {
			r.Data = ec.Data
		}
		if ec.HttpCode == 0 {
			ec.HttpCode = 500
		}

		r.json(c, ec.HttpCode)
		return
	}

	// 普通error
	r.Code = errcode.SystemError().Code
	if e != nil {
		r.Msg = e.Error()
		r.Msg = lang.Trans(ctx, r.Msg, nil)
		if log != nil {
			log.WithDebugger(ctx).Error(r.Msg)
		}
	} else {
		r.Msg = errcode.SystemError().Msg
		r.Msg = lang.Trans(ctx, r.Msg, nil)
		if log != nil {
			log.WithDebugger(ctx).Error(r.Msg)
		}
	}
	r.Data = []any{}
	r.json(c, http.StatusInternalServerError)
}
