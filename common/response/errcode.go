package response

import (
	"gin/common/errcode"
	"gin/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 通用响应结构体
type Response struct {
	Code int64       `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 返回数据
}

// Json 输出Json响应
func (r Response) Json(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, r)
	c.Abort()
}

// Success 返回成功响应,可传ErrorCode
func Success(c *gin.Context, e *errcode.ErrorCode) {
	var (
		resp Response
	)

	switch e {

	case nil:
		resp = Response{
			Code: errcode.Success().Code,
			Msg:  errcode.Success().Msg,
			Data: []string{},
		}

	default:
		if e.Data == nil {
			e.Data = []string{}
		}
		resp = Response{
			Code: e.Code,
			Msg:  e.Msg,
			Data: e.Data,
		}
	}

	resp.Json(c)
}

// Error 返回失败响应,可传ErrorCode
func Error(c *gin.Context, e *errcode.ErrorCode) {
	var (
		resp Response
		ctx  = c.Request.Context()
	)

	if e != nil {
		logger.NewLogger().WithDebugger(ctx).Error(e.Msg)
	} else {
		logger.NewLogger().WithDebugger(ctx).Error(errcode.SystemError().Msg)
	}

	switch e {

	case nil:
		resp = Response{
			Code: errcode.SystemError().Code,
			Msg:  errcode.SystemError().Msg,
			Data: []string{},
		}

	default:
		if e.Data == nil {
			e.Data = []string{}
		}
		resp = Response{
			Code: e.Code,
			Msg:  e.Msg,
			Data: e.Data,
		}
	}

	resp.Json(c)
}
