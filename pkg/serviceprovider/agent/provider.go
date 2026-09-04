package agent

import (
	"context"
	"gin/pkg/serviceprovider/mcp"
)

// ChatResult 对话结果
type ChatResult struct {
	Content   string          // 文本内容,无工具调用时返回
	ToolCalls []*ToolCallInfo // 工具调用信息
}

// ToolCallInfo 工具调用信息
type ToolCallInfo struct {
	Id        string // 工具调用ID
	Name      string // 工具名称
	Arguments string // 参数JSON
}

// MessageToolCall 消息中的工具调用
type MessageToolCall struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// Message 对话消息
type Message struct {
	Role       string             `json:"role"`
	Content    string             `json:"content,omitempty"`
	ToolCallId string             `json:"tool_call_id,omitempty"`
	ToolCalls  []*MessageToolCall `json:"tool_calls,omitempty"`
}

// MessageRecord 消息记录结构
type MessageRecord struct {
	Role       string  // 角色:user/assistant/tool/system
	Content    string  // 消息内容
	ToolName   string  // 工具名称
	ToolCallId string  // 工具调用ID
	ToolArgs   any     // 工具参数
	Tokens     int64   // 消耗Token数
	CostMs     float64 // 耗时(毫秒)
}

// Recorder 会话记录器接口
type Recorder interface {
	RecordMessage(ctx context.Context, sessionId int64, record MessageRecord) error
}

// Provider 模型提供商接口
type Provider interface {
	Name() string                                                              // 提供商名称
	Chat(messages []Message, tools []mcp.ToolDef) (*ChatResult, error)         // 对话
	StreamChat(messages []Message, tools []mcp.ToolDef) (<-chan string, error) // 流式对话
}
