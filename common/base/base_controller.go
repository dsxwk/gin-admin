package base

import (
	"gin/common/ctxkey"
	"gin/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	Response errcode.Response
}

// GetUserId 获取当前登录用户id
func (s *BaseController) GetUserId(ctx *gin.Context) int64 {
	return ctx.GetInt64(ctxkey.UserIdKey)
}
