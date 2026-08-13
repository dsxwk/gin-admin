package router

import (
	v1 "gin/app/controller/v1"
	"gin/app/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

func init() {
	Register(&AgentRouter{})
}

// AgentRouter Agent路由
type AgentRouter struct{}

func (r *AgentRouter) RegisterRoutes(routerGroup *gin.RouterGroup) {
	var agent v1.AgentController
	agentTimeout := middleware.Timeout{}.Handle(3 * time.Minute)

	router := routerGroup.Group("/api/v1/agent")
	{
		// 跳过全局超时,使用独立的3分钟超时
		router.POST("/ask", r.skipGlobalTimeout, agentTimeout, agent.Ask)
		router.GET("/history", agent.History)
		router.GET("/sessions", agent.Sessions)
	}
}

// skipGlobalTimeout 标记跳过全局超时中间件
func (r *AgentRouter) skipGlobalTimeout(c *gin.Context) {
	c.Set(middleware.SkipTimeoutKey, true)
	c.Next()
}

func (r *AgentRouter) IsAuth() bool {
	return true
}
