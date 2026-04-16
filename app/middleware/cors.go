package middleware

import (
	"gin/app/facade"
	"gin/common/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Cors struct {
	base.BaseMiddleware
}

var conf = facade.Config.Get()

// Handle 跨域中间件
func (s Cors) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf.Cors.Enabled {
			origin := c.Request.Header.Get("Origin")

			// 获取当前请求对应的跨域配置
			corsConfig := conf.Cors.GetConfig(origin)
			// 不在白名单中拒绝请求
			if corsConfig == nil {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			c.Header("Access-Control-Allow-Origin", corsConfig.AllowOrigin)
			c.Header("Access-Control-Allow-Headers", corsConfig.AllowHeaders)
			c.Header("Access-Control-Expose-Headers", corsConfig.ExposeHeaders)
			c.Header("Access-Control-Allow-Methods", corsConfig.AllowMethods)
			c.Header("Access-Control-Allow-Credentials", corsConfig.AllowCredentials)

			// 放行所有OPTIONS方法
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}

		c.Next()
	}
}
