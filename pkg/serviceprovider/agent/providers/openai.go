package providers

import (
	"bytes"
	"fmt"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/mcp"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

// OpenAICompat OpenAI兼容Provider,支持所有OpenAI兼容API
type OpenAICompat struct {
	name        string
	apiKey      string
	baseUrl     string
	model       string
	maxTokens   int
	temperature float64
	client      *http.Client
}

// NewOpenAICompat 创建OpenAI兼容Provider
func NewOpenAICompat(name string, cfg config.AgentProvider, maxTokens int, temperature float64) *OpenAICompat {
	return &OpenAICompat{
		name:        name,
		apiKey:      cfg.ApiKey,
		baseUrl:     cfg.BaseUrl,
		model:       cfg.Model,
		maxTokens:   maxTokens,
		temperature: temperature,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Name 提供商名称
func (p *OpenAICompat) Name() string {
	return p.name
}

// Chat 对话
func (p *OpenAICompat) Chat(messages []agent.Message, tools []mcp.ToolDef) (*agent.ChatResult, error) {
	// 构建请求体
	reqBody := chatRequest{
		Model:       p.model,
		Messages:    p.convertMessages(messages),
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
	}

	if len(tools) > 0 {
		reqBody.Tools = p.convertTools(tools)
		reqBody.ToolChoice = "auto"
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 调试打印请求体
	flag.Infof("[Agent] 请求模型:%s tools数:%d tool_choice:%s", p.model, len(tools), reqBody.ToolChoice)

	// 发送请求
	req, err := http.NewRequest("POST", p.baseUrl+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err = json.Unmarshal(respBody, &chatResp); err != nil {
		flag.Errorf("[Agent] 解析响应失败: %v, body: %s", err, string(respBody)[:200])
		return nil, err
	}

	// 调试:打印响应摘要
	if len(chatResp.Choices) > 0 {
		c := chatResp.Choices[0]
		flag.Infof("[Agent] 响应 content长度:%d tool_calls:%d", len(c.Message.Content), len(c.Message.ToolCalls))
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from model")
	}

	choice := chatResp.Choices[0]
	result := &agent.ChatResult{
		Content: choice.Message.Content,
	}

	// 转换工具调用
	for _, tc := range choice.Message.ToolCalls {
		result.ToolCalls = append(result.ToolCalls, &agent.ToolCallInfo{
			Id:        tc.Id,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}

	return result, nil
}

// StreamChat 流式对话
func (p *OpenAICompat) StreamChat(messages []agent.Message, tools []mcp.ToolDef) (<-chan string, error) {
	// TODO: SSE流式处理
	ch := make(chan string)
	close(ch)
	return ch, fmt.Errorf("stream not implemented yet")
}

// convertMessages 转换消息格式
func (p *OpenAICompat) convertMessages(messages []agent.Message) []chatMessage {
	result := make([]chatMessage, len(messages))
	for i, msg := range messages {
		cm := chatMessage{Role: msg.Role, Content: msg.Content, ToolCallId: msg.ToolCallId}
		// 转换tool_calls
		if len(msg.ToolCalls) > 0 {
			cm.ToolCalls = make([]chatMessageToolCall, len(msg.ToolCalls))
			for j, tc := range msg.ToolCalls {
				cm.ToolCalls[j] = chatMessageToolCall{
					Id:   tc.Id,
					Type: tc.Type,
					Function: chatToolCallFunc{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				}
			}
		}
		result[i] = cm
	}
	return result
}

// convertTools 转换工具定义
func (p *OpenAICompat) convertTools(tools []mcp.ToolDef) []chatTool {
	result := make([]chatTool, len(tools))
	for i, t := range tools {
		result[i] = chatTool{
			Type: "function",
			Function: chatToolFunc{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		}
	}
	return result
}

// chatRequest OpenAI聊天请求
type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Tools       []chatTool    `json:"tools,omitempty"`
	ToolChoice  string        `json:"tool_choice,omitempty"`
}

// chatMessage 聊天消息
type chatMessage struct {
	Role       string                `json:"role"`
	Content    string                `json:"content,omitempty"`
	ToolCallId string                `json:"tool_call_id,omitempty"`
	ToolCalls  []chatMessageToolCall `json:"tool_calls,omitempty"`
}

// chatMessageToolCall 消息中的工具调用(assistant消息)
type chatMessageToolCall struct {
	Id       string           `json:"id"`
	Type     string           `json:"type"`
	Function chatToolCallFunc `json:"function"`
}

// chatTool 工具定义(发送给API的tools列表)
type chatTool struct {
	Type     string       `json:"type"`
	Function chatToolFunc `json:"function"`
}

// chatToolFunc 工具函数定义(Name+Description+Parameters)
type chatToolFunc struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  mcp.InputSchema `json:"parameters"`
}

// chatToolCallFunc 工具调用函数(Name+Arguments)
type chatToolCallFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// chatResponse OpenAI聊天响应
type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

// chatChoice 聊天选项
type chatChoice struct {
	Message chatRespMessage `json:"message"`
}

// chatRespMessage 响应消息
type chatRespMessage struct {
	Content   string         `json:"content,omitempty"`
	ToolCalls []chatToolCall `json:"tool_calls,omitempty"`
}

// chatToolCall 工具调用(API响应)
type chatToolCall struct {
	Id       string           `json:"id"`
	Function chatToolCallFunc `json:"function"`
}
