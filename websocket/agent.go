package websocket

import (
	"context"
	"errors"
	"gin/app/facade"
	"gin/app/request"
	"gin/app/service"
	"gin/common/ctxkey"
	"gin/pkg/errcode"
	"gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/mcp"
	"gin/pkg/serviceprovider/ws"
	"strings"

	"github.com/goccy/go-json"
)

const (
	agentChatMessage  = "agent.chat"
	agentStartMessage = "agent.start"
	agentChunkMessage = "agent.chunk"
	agentDoneMessage  = "agent.done"
)

// AgentHandler AI流式对话处理器
type AgentHandler struct{}

// Match 匹配AI对话消息
func (h *AgentHandler) Match(message ws.Message) bool {
	var requestData MessageRequest
	if err := json.Unmarshal(message.Data, &requestData); err != nil {
		return false
	}
	return requestData.Type == agentChatMessage
}

// Handle 处理AI流式对话
func (h *AgentHandler) Handle(ctx context.Context, client ws.Connection, message ws.Message) error {
	userID := ws.UserID(ctx)
	if userID <= 0 {
		return ws.SendError(client, errcode.Unauthorized().WithMsg("请先登录"))
	}

	var requestData MessageRequest
	if err := json.Unmarshal(message.Data, &requestData); err != nil {
		return ws.SendError(client, errcode.ArgsError().WithMsg("消息格式错误"))
	}

	var payload request.AgentAsk
	if err := json.Unmarshal(requestData.Data, &payload); err != nil {
		return ws.SendError(client, errcode.ArgsError().WithMsg("对话参数错误"))
	}
	payload.Question = strings.TrimSpace(payload.Question)
	if payload.Question == "" {
		return ws.SendError(client, errcode.ArgsError().WithMsg("问题不能为空"))
	}

	providerName, modelName, ok := facade.AgentProvider(payload.Provider)
	if !ok {
		return ws.SendError(client, errors.New("agent未启用或配置错误"))
	}

	traceID, _ := ctx.Value(ctxkey.TraceIDKey).(string)
	var agentService service.AgentService
	sessionID := payload.SessionId
	if sessionID == 0 {
		title := payload.Question
		if len([]rune(title)) > 200 {
			title = string([]rune(title)[:200])
		}
		var err error
		sessionID, err = agentService.CreateSession(ctx, userID, title, providerName, modelName, traceID)
		if err != nil {
			return ws.SendError(client, err)
		}
	}

	history, err := agentService.LoadHistory(ctx, sessionID)
	if err != nil {
		return ws.SendError(client, err)
	}

	a := facade.Agent(providerName)
	if a == nil {
		return ws.SendError(client, errors.New("agent未启用或配置错误"))
	}
	mcpTools := make([]mcp.Tool, 0)
	if mcpHandler := facade.MCP(); mcpHandler != nil {
		mcpTools = mcpHandler.Tools()
	}
	a.WithSystemPrompt(mcp.BuildSystemPrompt("Gin-Admin后台AI助手", mcpTools)).
		WithSession(sessionID, &agentService).
		WithHistory(history).
		WithUserId(userID).
		WithContext(ctx)

	if err = ws.SendJSON(client, agentStartMessage, map[string]any{
		"sessionId": sessionID,
		"provider":  providerName,
		"model":     modelName,
	}); err != nil {
		return err
	}

	answer, err := a.StreamAsk(payload.Question, func(chunk agent.StreamChunk) error {
		if chunk.Content == "" {
			return nil
		}
		return ws.SendJSON(client, agentChunkMessage, map[string]any{
			"content": chunk.Content,
		})
	})
	if err != nil {
		return ws.SendError(client, err)
	}

	return ws.SendJSON(client, agentDoneMessage, map[string]any{
		"sessionId": sessionID,
		"provider":  providerName,
		"model":     modelName,
		"answer":    answer,
	})
}
