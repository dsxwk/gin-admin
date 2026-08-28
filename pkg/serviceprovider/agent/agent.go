package agent

import (
	"context"
	"errors"
	"gin/pkg/serviceprovider/mcp"
	"time"

	"github.com/goccy/go-json"
)

// Agent AI智能体
type Agent struct {
	provider     Provider
	systemPrompt string
	history      []Message
	tools        []mcp.Tool
	toolDefs     []mcp.ToolDef
	maxSteps     int             // 最大工具调用轮次,防止死循环
	userId       int64           // 当前操作用户ID
	ctx          context.Context // 请求上下文
	sessionId    int64           // 会话ID
	recorder     Recorder        // 会话记录器
}

// New 创建Agent
func New(provider Provider) *Agent {
	return &Agent{
		provider: provider,
		ctx:      context.Background(),
		maxSteps: 5,
		tools:    mcp.GetAll(),
		toolDefs: mcp.GetAllDefs(),
	}
}

// WithSystemPrompt 设置系统提示语
func (a *Agent) WithSystemPrompt(prompt string) *Agent {
	a.systemPrompt = prompt
	return a
}

// WithTools 设置可用工具
func (a *Agent) WithTools(tools []mcp.Tool) *Agent {
	a.tools = tools
	defs := make([]mcp.ToolDef, 0, len(tools))
	for _, t := range tools {
		defs = append(defs, mcp.ToolDef{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	a.toolDefs = defs
	return a
}

// WithMaxSteps 设置最大工具调用轮次
func (a *Agent) WithMaxSteps(n int) *Agent {
	a.maxSteps = n
	return a
}

// WithHistory 设置历史消息(从数据库加载)
func (a *Agent) WithHistory(history []Message) *Agent {
	a.history = history
	return a
}

// WithSession 设置会话ID和记录器
func (a *Agent) WithSession(sessionId int64, recorder Recorder) *Agent {
	a.sessionId = sessionId
	a.recorder = recorder
	return a
}

// WithUserId 设置当前操作用户ID
func (a *Agent) WithUserId(userId int64) *Agent {
	a.userId = userId
	return a
}

// Ask 发送问题并获取回答(自动调用工具)
func (a *Agent) Ask(question string) (string, error) {
	// 构建消息列表
	messages := a.buildMessages(question)

	// 记录用户消息
	a.record(MessageRecord{Role: "user", Content: question})

	for step := 0; step < a.maxSteps; step++ {
		// 调用模型
		start := time.Now()
		result, err := a.provider.Chat(messages, a.toolDefs)
		costMs := float64(time.Since(start).Milliseconds())
		if err != nil {
			return "", err
		}

		// 无工具调用,直接返回
		if len(result.ToolCalls) == 0 {
			a.record(MessageRecord{Role: "assistant", Content: result.Content, CostMs: costMs})
			return result.Content, nil
		}

		// 构建assistant消息(含tool_calls)
		assistantMsg := Message{Role: "assistant", Content: result.Content}
		for _, tc := range result.ToolCalls {
			assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, &MessageToolCall{
				Id:   tc.Id,
				Type: "function",
				Function: struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				}{Name: tc.Name, Arguments: tc.Arguments},
			})
		}
		messages = append(messages, assistantMsg)

		// 记录assistant工具调用消息
		for _, tc := range result.ToolCalls {
			a.record(MessageRecord{
				Role:       "assistant",
				Content:    result.Content,
				ToolName:   tc.Name,
				ToolCallId: tc.Id,
				ToolArgs:   a.parseArgs(tc.Arguments),
				CostMs:     costMs,
			})
		}

		// 执行每个工具调用并追加结果
		for _, tc := range result.ToolCalls {
			// 解析参数
			args := a.parseArgs(tc.Arguments)

			toolResult, err := a.callTool(tc.Name, args)
			if err != nil {
				toolResult = map[string]any{"error": err.Error()}
			}
			resultJson, _ := json.Marshal(toolResult)

			messages = append(messages, Message{
				Role:       "tool",
				ToolCallId: tc.Id,
				Content:    string(resultJson),
			})

			// 记录tool结果消息
			a.record(MessageRecord{
				Role:       "tool",
				Content:    string(resultJson),
				ToolName:   tc.Name,
				ToolCallId: tc.Id,
			})
		}
	}

	return "", errors.New("达到最大工具调用轮次")
}

// parseArgs 解析工具参数
func (a *Agent) parseArgs(arguments string) map[string]any {
	var args map[string]any
	if arguments != "" {
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			args = map[string]any{"raw": arguments}
		}
	}
	return args
}

// record 记录消息
func (a *Agent) record(record MessageRecord) {
	if a.recorder == nil || a.sessionId == 0 {
		return
	}
	_ = a.recorder.RecordMessage(a.sessionId, record)
}

// buildMessages 构建消息列表
func (a *Agent) buildMessages(question string) []Message {
	messages := make([]Message, 0)

	if a.systemPrompt != "" {
		messages = append(messages, Message{Role: "system", Content: a.systemPrompt})
	}

	// 追加历史
	messages = append(messages, a.history...)

	// 追加当前问题
	messages = append(messages, Message{Role: "user", Content: question})

	return messages
}

// callTool 调用工具
func (a *Agent) callTool(name string, args map[string]any) (any, error) {
	tool, ok := mcp.Get(name)
	if !ok {
		return nil, errors.New("tool not found: " + name)
	}
	return mcp.CallTool(a.ctx, a.userId, tool, args)
}

// Reset 重置对话历史
func (a *Agent) Reset() {
	a.history = nil
}

// WithContext 设置请求上下文
func (a *Agent) WithContext(ctx context.Context) *Agent {
	if ctx != nil {
		a.ctx = ctx
	}
	return a
}
