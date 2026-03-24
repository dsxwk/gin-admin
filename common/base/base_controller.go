package base

import (
	"gin/common/ctxkey"
	"gin/common/errcode"
	"gin/common/response"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// GetUserId 获取当前登录用户id
func (s *BaseController) GetUserId(ctx *gin.Context) int64 {
	id, _ := ctx.Get(ctxkey.UserIdKey)
	uid := id.(float64)
	return int64(uid)
}

// Success 成功返回
func (s *BaseController) Success(c *gin.Context, e errcode.ErrorCode) {
	response.Success(c, &e)
}

// Error 失败返回
func (s *BaseController) Error(c *gin.Context, e errcode.ErrorCode) {
	response.Error(c, &e)
}
