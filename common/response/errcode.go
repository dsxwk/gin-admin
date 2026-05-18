package response

import (
	"gin/common/errcode"
	"gin/pkg/provider/logger"
	"github.com/gin-gonic/gin"
)

var (
	log *logger.Logger
)

// Response 通用响应结构体
type Response struct {
	Code int64       `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 返回数据
	log  *logger.Logger
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
func (r Response) Success(c *gin.Context, e errcode.ErrorCode) {
	var (
		resp Response
	)

	if e.Msg == "" {
		e.Msg = errcode.Success().Msg
	}
	if e.Data == nil {
		e.Data = []string{}
	}
	if e.HttpCode == 0 {
		e.HttpCode = 200
	}
	resp = Response{
		Code: e.Code,
		Msg:  e.Msg,
		Data: e.Data,
	}

	resp.json(c, e.HttpCode)
}

// Error 返回失败响应,可传ErrorCode
func (r Response) Error(c *gin.Context, e errcode.ErrorCode) {
	var (
		resp Response
		ctx  = c.Request.Context()
	)

	if e.Msg != "" {
		log.WithDebugger(ctx).Error(e.Msg)
	} else {
		log.WithDebugger(ctx).Error(errcode.SystemError().Msg)
	}

	if e.Msg == "" {
		e.Msg = errcode.SystemError().Msg
	}
	if e.Data == nil {
		e.Data = []string{}
	}
	if e.HttpCode == 0 {
		e.HttpCode = 500
	}
	resp = Response{
		Code: e.Code,
		Msg:  e.Msg,
		Data: e.Data,
	}

	resp.json(c, e.HttpCode)
}
