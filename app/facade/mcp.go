package facade

import (
	"gin/pkg/serviceprovider/mcp"
)

// MCP 获取MCP处理器
func MCP() *mcp.Handler {
	return App().MCP()
}
