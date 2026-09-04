package v1

import (
	"errors"
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/common/errcode"
	"gin/pkg/serviceprovider/mcp"

	"github.com/gin-gonic/gin"
)

type AgentController struct {
	base.BaseController
	service service.AgentService
}

// Ask AI对话
// @Tags AI助手
// @Summary AI对话
// @Description AI Agent对话接口,自动调用MCP工具,数据持久化到数据库
// @Param token header string true "认证Token"
// @Param data body request.AgentAsk true "对话参数"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Router /api/v1/agent/ask [post]
func (s *AgentController) Ask(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.AgentAsk
	)

	if err := facade.Request().BindValidate(c, &req, "Ask"); err != nil {
		s.Response.Error(c, err)
		return
	}

	// 解析提供商信息
	providerName, modelName, ok := facade.AgentProvider(req.Provider)
	if !ok {
		s.Response.Error(c, errors.New("agent未启用或配置错误"))
		return
	}

	// 获取用户ID和追踪ID
	userId := s.GetUserId(c)
	traceId, _ := ctx.Value(ctxkey.TraceIdKey).(string)

	// 创建或加载会话
	sessionId := req.SessionId
	if sessionId == 0 {
		title := req.Question
		if len([]rune(title)) > 200 {
			title = string([]rune(title)[:200])
		}
		var err error
		sessionId, err = s.service.CreateSession(ctx, userId, title, providerName, modelName, traceId)
		if err != nil {
			s.Response.Error(c, err)
			return
		}
	}

	// 加载历史消息
	history, err := s.service.LoadHistory(ctx, sessionId)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	// 构建Agent并注入会话记录器
	a := facade.Agent(providerName)
	if a == nil {
		s.Response.Error(c, errors.New("agent未启用或配置错误"))
		return
	}
	a.WithSystemPrompt(mcp.BuildSystemPrompt("Gin-Admin后台AI助手")).
		WithSession(sessionId, &s.service).
		WithHistory(history).
		WithUserId(userId).
		WithContext(ctx)

	result, err := a.Ask(req.Question)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(map[string]any{
		"sessionId": sessionId,
		"provider":  providerName,
		"model":     modelName,
		"answer":    result,
	}))
}

// Sessions 会话列表
// @Tags AI助手
// @Summary 会话列表
// @Description 获取当前用户的会话列表
// @Param token header string true "认证Token"
// @Success 200 {object} errcode.SuccessResponse{data=[]model.AgentSession} "成功"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/agent/sessions [get]
func (s *AgentController) Sessions(c *gin.Context) {
	ctx := c.Request.Context()

	userId := s.GetUserId(c)
	sessions, err := s.service.ListSessions(ctx, userId)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(sessions))
}

// History 会话历史记录
// @Tags AI助手
// @Summary 会话历史记录
// @Description 获取指定会话的历史消息记录
// @Param token header string true "认证Token"
// @Param sessionId query int true "会话ID"
// @Success 200 {object} errcode.SuccessResponse{data=[]model.AgentMessage} "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 500 {object} errcode.SystemErrorResponse "系统错误"
// @Router /api/v1/agent/history [get]
func (s *AgentController) History(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.AgentAsk
	)

	if err := facade.Request().BindValidate(c, &req, "History"); err != nil {
		s.Response.Error(c, err)
		return
	}

	messages, err := s.service.History(ctx, req.SessionId)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success().WithData(messages))
}
