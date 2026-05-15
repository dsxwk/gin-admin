package base

import (
	"gin/common/ctxkey"
	"gin/common/response"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
	Response response.Response
}

// GetUserId 获取当前登录用户id
func (s *BaseController) GetUserId(ctx *gin.Context) int64 {
	id, _ := ctx.Get(ctxkey.UserIdKey)
	uid := id.(float64)
	return int64(uid)
}
