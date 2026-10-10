package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin/app/errcode"
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/mcp"
	"net/http"
	"time"

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
	traceId, _ := ctx.Value(ctxkey.TraceIDKey).(string)

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
	a.WithSystemPrompt(mcp.BuildSystemPrompt("Gin-Admin后台AI助手", facade.MCP().Tools())).
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

// Stream 流式对话(SSE)
// @Tags AI助手
// @Summary 流式对话(SSE)
// @Description 使用SSE返回模型流式输出,支持客户端断开取消
// @Param token header string true "认证Token"
// @Param question query string true "问题"
// @Param provider query string false "模型提供商"
// @Param sessionId query int false "会话ID"
// @Success 200 {string} string "SSE流"
// @Router /api/v1/agent/stream [get]
func (s *AgentController) Stream(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.AgentAsk
	)

	//流式请求单独限制总时长,避免上游长时间无响应导致连接永不结束
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	if err := facade.Request().BindValidate(c, &req, "Stream"); err != nil {
		s.Response.Error(c, err)
		return
	}

	providerName, modelName, ok := facade.AgentProvider(req.Provider)
	if !ok {
		s.Response.Error(c, errors.New("agent未启用或配置错误"))
		return
	}

	userId := s.GetUserId(c)
	traceId, _ := ctx.Value(ctxkey.TraceIDKey).(string)
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

	history, err := s.service.LoadHistory(ctx, sessionId)
	if err != nil {
		s.Response.Error(c, err)
		return
	}

	a := facade.Agent(providerName)
	if a == nil {
		s.Response.Error(c, errors.New("agent未启用或配置错误"))
		return
	}
	mcpTools := make([]mcp.Tool, 0)
	if mcpHandler := facade.MCP(); mcpHandler != nil {
		mcpTools = mcpHandler.Tools()
	}
	a.WithSystemPrompt(mcp.BuildSystemPrompt("Gin-Admin后台AI助手", mcpTools)).
		WithSession(sessionId, &s.service).
		WithHistory(history).
		WithUserId(userId).
		WithContext(ctx)

	sseConfig := facade.Config().Agent.Sse
	heartbeat := sseConfig.Heartbeat
	if heartbeat <= 0 {
		heartbeat = 15 * time.Second
	}
	retry := sseConfig.Retry
	if retry <= 0 {
		retry = 3000
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	_, _ = fmt.Fprintf(c.Writer, "retry: %d\n\n", retry)
	c.Writer.Flush()

	type streamEvent struct {
		chunk  *agent.StreamChunk
		answer string
		err    error
		done   bool
	}

	events := make(chan streamEvent, 32)
	go func() {
		answer, streamErr := a.StreamAsk(req.Question, func(chunk agent.StreamChunk) error {
			select {
			case events <- streamEvent{chunk: &chunk}:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
		select {
		case events <- streamEvent{answer: answer, err: streamErr, done: true}:
		case <-ctx.Done():
		}
		close(events)
	}()

	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = fmt.Fprint(c.Writer, ": ping\n\n")
			c.Writer.Flush()
		case event, ok := <-events:
			if !ok {
				return
			}
			if event.chunk != nil {
				if event.chunk.Err != nil {
					_ = writeAgentSSE(c, "error", map[string]any{"message": event.chunk.Err.Error()})
					return
				}
				if event.chunk.Content != "" {
					if err = writeAgentSSE(c, "message", map[string]any{"content": event.chunk.Content}); err != nil {
						return
					}
				}
				continue
			}
			if event.done {
				if event.err != nil {
					_ = writeAgentSSE(c, "error", map[string]any{"message": event.err.Error()})
					return
				}
				_ = writeAgentSSE(c, "done", map[string]any{
					"sessionId": sessionId,
					"provider":  providerName,
					"model":     modelName,
					"answer":    event.answer,
				})
				return
			}
		}
	}
}

// writeAgentSSE 写入SSE事件
func writeAgentSSE(c *gin.Context, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, payload); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
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

// Delete 删除会话及全部消息
// @Tags AI助手
// @Summary 删除会话
// @Description 删除当前用户指定会话及其全部消息
// @Param token header string true "认证Token"
// @Param id path int true "会话ID"
// @Success 200 {object} errcode.SuccessResponse "成功"
// @Failure 400 {object} errcode.ArgsErrorResponse "参数错误"
// @Failure 404 {object} errcode.SuccessResponse "会话不存在"
// @Router /api/v1/agent/sessions/{id} [delete]
func (s *AgentController) Delete(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req request.AgentDelete
	)

	req.ID = facade.Request().Path[int64](c, "id", 0)
	if err := facade.Request().BindValidate(c, &req, "Delete"); err != nil {
		s.Response.Error(c, err)
		return
	}

	if err := s.service.DeleteSession(ctx, s.GetUserId(c), req.ID); err != nil {
		s.Response.Error(c, err)
		return
	}

	s.Response.Success(c, errcode.Success())
}
