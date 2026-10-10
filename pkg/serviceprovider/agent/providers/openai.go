package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/mcp"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// OpenAICompat OpenAI兼容Provider,支持所有OpenAI兼容API
type OpenAICompat struct {
	name           string
	apiKey         string
	baseUrl        string
	model          string
	maxTokens      int
	temperature    float64
	requestTimeout time.Duration
	client         *http.Client
	streamClient   *http.Client
}

// NewOpenAICompat 创建OpenAI兼容Provider
func NewOpenAICompat(name string, cfg config.AgentProvider, maxTokens int, temperature float64, requestTimeout time.Duration) *OpenAICompat {
	return &OpenAICompat{
		name:           name,
		apiKey:         cfg.ApiKey,
		baseUrl:        cfg.BaseUrl,
		model:          cfg.Model,
		maxTokens:      maxTokens,
		temperature:    temperature,
		requestTimeout: requestTimeout,
		client:         newRequestClient(),
		streamClient:   newStreamClient(),
	}
}

// newStreamClient 创建流式客户端,仅限制连接与响应头时间避免长连接被整体中断
func newStreamClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 30 * time.Second
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	return &http.Client{Transport: transport}
}

// newRequestClient 创建非流式客户端,超时由context和requestTimeout控制
func newRequestClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 0
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	return &http.Client{Transport: transport}
}

// Name 提供商名称
func (p *OpenAICompat) Name() string {
	return p.name
}

// Chat 对话
func (p *OpenAICompat) Chat(ctx context.Context, messages []agent.Message, tools []mcp.ToolDef) (*agent.ChatResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if p.requestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.requestTimeout)
		defer cancel()
	}
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseUrl+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("模型请求超时: %w", err)
		}
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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
func (p *OpenAICompat) StreamChat(ctx context.Context, messages []agent.Message, tools []mcp.ToolDef) (<-chan agent.StreamChunk, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	reqBody := chatRequest{
		Model:       p.model,
		Messages:    p.convertMessages(messages),
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
		Stream:      true,
	}
	if len(tools) > 0 {
		reqBody.Tools = p.convertTools(tools)
		reqBody.ToolChoice = "auto"
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseUrl+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	flag.Infof("[Agent] 流式请求模型:%s tools数:%d", p.model, len(tools))

	resp, err := p.streamClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		defer func() { _ = resp.Body.Close() }()
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}
		flag.Errorf("[Agent] 流式响应异常 %d: %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan agent.StreamChunk, 16)
	go p.readStream(ctx, resp.Body, ch)
	return ch, nil
}

// readStream 读取OpenAI SSE流
func (p *OpenAICompat) readStream(ctx context.Context, body io.ReadCloser, ch chan<- agent.StreamChunk) {
	defer close(ch)
	defer func() { _ = body.Close() }()

	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	send := func(chunk agent.StreamChunk) bool {
		select {
		case ch <- chunk:
			return true
		case <-ctx.Done():
			return false
		}
	}

	//按index累积流式工具调用分片
	type toolCallAcc struct {
		id   string
		name string
		args strings.Builder
	}
	accs := make(map[int]*toolCallAcc)
	var order []int

	collectToolCalls := func() []*agent.ToolCallInfo {
		if len(order) == 0 {
			return nil
		}
		calls := make([]*agent.ToolCallInfo, 0, len(order))
		for _, index := range order {
			acc := accs[index]
			calls = append(calls, &agent.ToolCallInfo{
				Id:        acc.id,
				Name:      acc.name,
				Arguments: acc.args.String(),
			})
		}
		return calls
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}
		if data == "[DONE]" {
			send(agent.StreamChunk{Done: true, ToolCalls: collectToolCalls()})
			return
		}

		var streamResp streamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			flag.Errorf("[Agent] 解析流式响应失败: %v, body: %s", err, data)
			send(agent.StreamChunk{Err: err})
			return
		}
		if len(streamResp.Choices) == 0 {
			continue
		}

		delta := streamResp.Choices[0].Delta
		for _, tc := range delta.ToolCalls {
			acc, ok := accs[tc.Index]
			if !ok {
				acc = &toolCallAcc{}
				accs[tc.Index] = acc
				order = append(order, tc.Index)
			}
			if tc.Id != "" {
				acc.id = tc.Id
			}
			acc.name += tc.Function.Name
			acc.args.WriteString(tc.Function.Arguments)
		}

		if delta.Content != "" && !send(agent.StreamChunk{Content: delta.Content}) {
			return
		}
	}

	if err := scanner.Err(); err != nil {
		flag.Errorf("[Agent] 读取流式响应失败: %v", err)
		send(agent.StreamChunk{Err: err})
		return
	}
	send(agent.StreamChunk{Done: true, ToolCalls: collectToolCalls()})
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
	Stream      bool          `json:"stream,omitempty"`
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

// streamResponse 流式响应
type streamResponse struct {
	Choices []streamChoice `json:"choices"`
}

// streamChoice 流式选项
type streamChoice struct {
	Delta        streamDelta `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}

// streamDelta 流式增量
type streamDelta struct {
	Content   string           `json:"content"`
	ToolCalls []streamToolCall `json:"tool_calls"`
}

// streamToolCall 流式工具调用分片
type streamToolCall struct {
	Index    int                `json:"index"`
	Id       string             `json:"id"`
	Type     string             `json:"type"`
	Function streamToolCallFunc `json:"function"`
}

// streamToolCallFunc 流式工具调用函数
type streamToolCallFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
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
