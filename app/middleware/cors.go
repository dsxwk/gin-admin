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
			c.Header("Access-Control-Allow-Origin", conf.Cors.AllowOrigin)
			c.Header("Access-Control-Allow-Headers", conf.Cors.AllowHeaders)
			c.Header("Access-Control-Expose-Headers", conf.Cors.ExposeHeaders)
			c.Header("Access-Control-Allow-Methods", conf.Cors.AllowMethods)
			c.Header("Access-Control-Allow-Credentials", conf.Cors.AllowCredentials)

			// 放行所有OPTIONS方法
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}

		c.Next()
	}
}
