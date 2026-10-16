package base

import (
	"gin/pkg/ctxkey"
	"gin/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Controller 基础控制器
type Controller struct {
	Response errcode.Response
}

// GetUserID 获取当前登录用户id
func (s *Controller) GetUserID(ctx *gin.Context) int64 {
	return ctx.GetInt64(ctxkey.GinUserIDKey)
}
