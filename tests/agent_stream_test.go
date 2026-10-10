package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/config"
	"gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/agent/providers"
	"gin/pkg/serviceprovider/mcp"
	wsprovider "gin/pkg/serviceprovider/ws"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// newAgentStreamServer 创建模拟OpenAI SSE服务
func newAgentStreamServer(t *testing.T, requestBody chan map[string]any) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(request.Body).Decode(&body); err == nil && requestBody != nil {
			requestBody <- body
		}

		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		flusher, ok := writer.(http.Flusher)
		if !ok {
			return
		}

		flusher.Flush()
		_, _ = fmt.Fprint(writer, "data: {\"choices\":[{\"delta\":{\"content\":\"你好\"}}]}\n\n")
		flusher.Flush()
		_, _ = fmt.Fprint(writer, "data: {\"choices\":[{\"delta\":{\"content\":\"，流式对话\"}}]}\n\n")
		flusher.Flush()
		_, _ = fmt.Fprint(writer, "data: [DONE]\n\n")
		flusher.Flush()
	}))

	t.Cleanup(server.Close)
	return server
}

// TestOpenAICompatStreamChat OpenAI兼容SSE流式对话测试
func TestOpenAICompatStreamChat(t *testing.T) {
	requestBody := make(chan map[string]any, 1)
	server := newAgentStreamServer(t, requestBody)

	provider := providers.NewOpenAICompat("stream_test", config.AgentProvider{
		BaseUrl: server.URL,
		Model:   "stream-test",
	}, 1024, 0.7, 30*time.Second)

	stream, err := provider.StreamChat(t.Context(), []agent.Message{
		{Role: "user", Content: "你好"},
	}, nil)
	require.NoError(t, err)

	var answer string
	done := false
	for chunk := range stream {
		require.NoError(t, chunk.Err)
		answer += chunk.Content
		if chunk.Done {
			done = true
		}
	}

	require.True(t, done)
	require.Equal(t, "你好，流式对话", answer)

	select {
	case body := <-requestBody:
		require.Equal(t, true, body["stream"])
		require.Equal(t, "stream-test", body["model"])
	case <-time.After(time.Second):
		t.Fatal("未收到模型请求")
	}
}

// agentStreamTestProvider 测试流式Provider
type agentStreamTestProvider struct{}

// Name 提供商名称
func (p *agentStreamTestProvider) Name() string { return "test" }

// Chat 普通对话
func (p *agentStreamTestProvider) Chat(_ context.Context, _ []agent.Message, _ []mcp.ToolDef) (*agent.ChatResult, error) {
	return &agent.ChatResult{Content: "你好，流式对话"}, nil
}

// StreamChat 流式对话
func (p *agentStreamTestProvider) StreamChat(_ context.Context, _ []agent.Message, _ []mcp.ToolDef) (<-chan agent.StreamChunk, error) {
	stream := make(chan agent.StreamChunk, 3)
	stream <- agent.StreamChunk{Content: "你好"}
	stream <- agent.StreamChunk{Content: "，流式对话"}
	stream <- agent.StreamChunk{Done: true}
	close(stream)
	return stream, nil
}

// agentStreamTestHandler 测试流式处理器
type agentStreamTestHandler struct{}

// Match 匹配全部消息
func (h *agentStreamTestHandler) Match(_ wsprovider.Message) bool { return true }

// Handle 处理流式对话
func (h *agentStreamTestHandler) Handle(ctx context.Context, connection wsprovider.Connection, _ wsprovider.Message) error {
	if err := wsprovider.SendJSON(connection, "agent.start", map[string]any{"sessionId": 1}); err != nil {
		return err
	}

	a := agent.New(&agentStreamTestProvider{}, nil).WithContext(ctx)
	answer, err := a.StreamAsk("你好", func(chunk agent.StreamChunk) error {
		if chunk.Content == "" {
			return nil
		}
		return wsprovider.SendJSON(connection, "agent.chunk", map[string]any{"content": chunk.Content})
	})
	if err != nil {
		return wsprovider.SendError(connection, err)
	}

	return wsprovider.SendJSON(connection, "agent.done", map[string]any{"answer": answer})
}

// TestOpenAICompatChatTimeout 非流式请求超时测试
func TestOpenAICompatChatTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(100 * time.Millisecond)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(writer, `{"choices":[{"message":{"content":"ok"}}]}`)
	}))
	t.Cleanup(server.Close)

	provider := providers.NewOpenAICompat("timeout_test", config.AgentProvider{
		BaseUrl: server.URL,
		Model:   "timeout-test",
	}, 128, 0.7, 20*time.Millisecond)

	_, err := provider.Chat(t.Context(), []agent.Message{{Role: "user", Content: "你好"}}, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

// TestWebSocketAgentStream WebSocket流式对话测试
func TestWebSocketAgentStream(t *testing.T) {
	manager := wsprovider.NewManager(
		wsprovider.Options{Shards: 1, MaxConnections: 2},
		&agentStreamTestHandler{},
	)
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	conn := dialWsTestServer(t, newWsAuthTestServer(t, manager, 1).URL)
	writeWsMessage(t, conn, websocket.TextMessage, []byte(`{"type":"agent.chat","data":{"question":"你好"}}`))

	var answer string
	done := false
	for index := 0; index < 10; index++ {
		_, data := readWsMessage(t, conn)
		var response wsResponse
		require.NoError(t, json.Unmarshal(data, &response))

		switch response.Type {
		case "agent.start":
			require.Contains(t, string(response.Data), "sessionId")
		case "agent.chunk":
			var chunk struct {
				Content string `json:"content"`
			}
			require.NoError(t, json.Unmarshal(response.Data, &chunk))
			answer += chunk.Content
		case "agent.done":
			done = true
			var result struct {
				Answer string `json:"answer"`
			}
			require.NoError(t, json.Unmarshal(response.Data, &result))
			require.Equal(t, "你好，流式对话", result.Answer)
		case "error":
			t.Fatalf("WebSocket流式对话失败: %s", string(response.Data))
		}

		if done {
			break
		}
	}

	require.True(t, done)
	require.Equal(t, "你好，流式对话", answer)
}
