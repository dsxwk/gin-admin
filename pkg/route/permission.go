package route

import (
	"io"

	"github.com/gin-gonic/gin"
)

// PermissionKeys 提取需要鉴权的路由权限Key
func PermissionKeys() []string {
	// 临时禁用Gin调试输出避免重复打印路由
	orig := gin.DefaultWriter
	gin.DefaultWriter = io.Discard
	defer func() {
		gin.DefaultWriter = orig
	}()

	engine := gin.New()
	group := engine.Group("")
	for _, r := range snapshot() {
		if r.IsAuth() {
			r.RegisterRoutes(group)
		}
	}

	return extractPermissionKeys(engine)
}

// extractPermissionKeys 提取去重的权限Key
func extractPermissionKeys(engine *gin.Engine) []string {
	routes := engine.Routes()
	seen := make(map[string]struct{}, len(routes))
	keys := make([]string, 0, len(routes))
	for _, route := range routes {
		key := route.Method + ":" + route.Path
		if _, ok := seen[key]; ok {
			continue
		}

		keys = append(keys, key)
		seen[key] = struct{}{}
	}

	return keys
}
