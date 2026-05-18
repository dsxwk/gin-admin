package base

import (
	"gin/app/facade"
	"gin/common/ctxkey"
	"gin/common/response"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
	Response response.Response
}

func init() {
	response.SetLogger(facade.Log())
}

// GetUserId 获取当前登录用户id
func (s *BaseController) GetUserId(ctx *gin.Context) int64 {
	id, _ := ctx.Get(ctxkey.UserIdKey)
	uid := id.(float64)
	return int64(uid)
}
