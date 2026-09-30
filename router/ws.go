package router

import (
	"strings"

	"gin/app/facade"
	"gin/app/middleware"
	"gin/common/ctxkey"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider/ws"

	"github.com/gin-gonic/gin"
)

// WsRouter WebSocket路由
type WsRouter struct{}

// Register 注册路由
func (r *WsRouter) Register(routerGroup *gin.RouterGroup) {
	cfg := facade.Config()
	if cfg == nil || !cfg.Ws.Enabled {
		return
	}

	manager := facade.WS()
	if manager == nil {
		return
	}

	path := cfg.Ws.Path
	if path == "" {
		path = "/ws"
	}

	routerGroup.GET(path, func(c *gin.Context) {
		userID, ok := websocketUserID(c)
		if !ok {
			return
		}

		if userID > 0 {
			ctx := ws.WithUserID(c.Request.Context(), userID)
			ctx = ctxkey.WithValue(ctx, ctxkey.UserIdKey, userID)
			c.Request = c.Request.WithContext(ctx)
		}

		manager.ServeHTTP(c.Writer, c.Request)
	})
}

// IsAuth 是否需要鉴权
func (r *WsRouter) IsAuth() bool {
	return false
}

// websocketUserID 解析WebSocket登录用户
func websocketUserID(c *gin.Context) (int64, bool) {
	token := websocketToken(c)
	if token == "" {
		return 0, true
	}

	claims, err := (middleware.Jwt{}).Decode(token)
	if err != nil {
		facade.Response().Error(c, errcode.Unauthorized().WithMsg(err.Error()))
		return 0, false
	}

	userID := jwtUserID(claims)
	if userID <= 0 {
		facade.Response().Error(c, errcode.Unauthorized().WithMsg("登录用户无效"))
		return 0, false
	}

	return userID, true
}

// websocketToken 获取WebSocket登录凭证
func websocketToken(c *gin.Context) string {
	token := strings.TrimSpace(c.Query("token"))
	if token != "" {
		return token
	}

	token = strings.TrimSpace(c.GetHeader("token"))
	if token != "" {
		return token
	}

	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if len(authorization) <= 7 || !strings.EqualFold(authorization[:7], "Bearer ") {
		return ""
	}
	return strings.TrimSpace(authorization[7:])
}

// jwtUserID 获取JWT用户ID
func jwtUserID(claims map[string]any) int64 {
	switch value := claims["id"].(type) {
	case float64:
		return int64(value)
	case int64:
		return value
	case int:
		return int64(value)
	default:
		return 0
	}
}
