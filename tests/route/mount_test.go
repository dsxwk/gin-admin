package route

import (
	"testing"

	serviceroute "gin/pkg/route"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mountPublicRouter struct{}

// Register 注册公开路由
func (r *mountPublicRouter) Register(router *gin.RouterGroup) {
	router.GET("/public", func(_ *gin.Context) {})
}

// IsAuth 是否需要鉴权
func (r *mountPublicRouter) IsAuth() bool {
	return false
}

type mountAuthRouter struct{}

// Register 注册鉴权路由
func (r *mountAuthRouter) Register(router *gin.RouterGroup) {
	router.GET("/users", func(_ *gin.Context) {})
	router.POST("/users", func(_ *gin.Context) {})
}

// IsAuth 是否需要鉴权
func (r *mountAuthRouter) IsAuth() bool {
	return true
}

// TestMountPermissionKeys 测试挂载时收集鉴权路由
func TestMountPermissionKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	keys := serviceroute.Mount(
		engine,
		engine.Group(""),
		nil,
		&mountPublicRouter{},
		&mountAuthRouter{},
	)

	require.ElementsMatch(t, []string{
		"GET:/users",
		"POST:/users",
	}, keys)
	require.NotContains(t, keys, "GET:/public")
}
