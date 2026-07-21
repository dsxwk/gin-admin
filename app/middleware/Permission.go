package middleware

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/common/errcode"
	"github.com/gin-gonic/gin"
)

type Permission struct {
	base.BaseMiddleware
}

// Handle 权限中间件
func (s Permission) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Request.Method + ":" + c.FullPath()

		if !s.hasPermission(c, code) {
			s.Response.Error(c, errcode.Forbidden().WithMsg("No permission"))
			return
		}
		c.Next()
	}
}

func (s Permission) hasPermission(c *gin.Context, code string) bool {
	userID := c.GetInt64(ctxkey.UserIdKey)

	isMember, err := facade.Cache("redis").
		Redis().
		WithContext(c.Request.Context()).
		SIsMember(
			fmt.Sprintf(
				"permission:user:%d",
				userID,
			),
			code,
		)
	if err != nil {
		return false
	}
	return isMember
}
