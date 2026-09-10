package middleware

import (
	"gin/app/facade"
	"gin/app/model"
	"gin/common/base"
	"gin/common/ctxkey"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

type OperatorLog struct {
	base.BaseMiddleware
}

// Handle 操作日志中间件
func (s OperatorLog) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			// 是否记录操作日
			if !facade.Config().OperatorRecord.Enable {
				return
			}
			// 从gin上下文获取userId
			var userId int64
			if v, exists := c.Get(ctxkey.UserIdKey); exists {
				if id, ok := v.(int64); ok {
					userId = id
				}
			}

			// 从request上下文获取traceId、lang、params
			ctx := c.Request.Context()
			traceId, _ := ctx.Value(ctxkey.TraceIdKey).(string)
			lang, _ := ctx.Value(ctxkey.LangKey).(string)
			params := ctx.Value(ctxkey.ParamsKey)
			header, _ := json.Marshal(c.Request.Header)

			// 构建操作日志
			log := model.OperatorLog{
				Ip:         c.ClientIP(),
				Method:     c.Request.Method,
				Header:     string(header),
				Uri:        c.Request.URL.Path,
				Lang:       lang,
				Params:     &model.JsonValue{Data: params},
				UserId:     userId,
				TraceId:    traceId,
				StatusCode: int64(c.Writer.Status()),
				UserAgent:  c.Request.UserAgent(),
				CostMs:     float64(time.Since(start).Milliseconds()),
			}

			// 异步写入避免阻塞请求
			go func() {
				facade.DB(log.Connection()).Model(&model.OperatorLog{}).Create(&log)
			}()
		}()

		c.Next()
	}
}
