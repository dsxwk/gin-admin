package router

import (
	"gin/app/facade"

	"github.com/gin-gonic/gin"
)

// McpRouter MCP服务路由
type McpRouter struct{}

// Register 注册路由
func (r *McpRouter) Register(routerGroup *gin.RouterGroup) {
	cfg := facade.Config()
	if cfg == nil || !cfg.Mcp.Enabled {
		return
	}

	handler := facade.MCP()
	if handler == nil {
		return
	}

	path := cfg.Mcp.Path
	if path == "" {
		path = "/mcp"
	}

	routerGroup.POST(path, gin.WrapH(handler))
}

// IsAuth 是否需要鉴权
func (r *McpRouter) IsAuth() bool {
	return false
}
