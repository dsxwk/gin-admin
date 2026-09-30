package tests

import (
	"context"
	"gin/pkg/serviceprovider/ws"
	ws1 "gin/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

// wsResponse WebSocket测试响应
type wsResponse struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// wsErrorResponse WebSocket测试错误响应
type wsErrorResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// wsMatchedTestHandler WebSocket匹配处理器测试
type wsMatchedTestHandler struct{}

// Match 匹配测试消息
func (h *wsMatchedTestHandler) Match(message ws.Message) bool {
	return string(message.Data) == `{"type":"notice"}`
}

// Handle 处理匹配消息
func (h *wsMatchedTestHandler) Handle(_ context.Context, connection ws.Connection, _ ws.Message) error {
	return connection.Send(ws.Message{
		Type: ws.MessageText,
		Data: []byte("matched"),
	})
}

// wsFallbackTestHandler WebSocket兜底处理器测试
type wsFallbackTestHandler struct{}

// Handle 处理兜底消息
func (h *wsFallbackTestHandler) Handle(_ context.Context, connection ws.Connection, _ ws.Message) error {
	return connection.Send(ws.Message{
		Type: ws.MessageText,
		Data: []byte("fallback"),
	})
}

// newWsTestManager 创建WebSocket测试管理器
func newWsTestManager(options ws.Options) *ws.Manager {
	return ws.NewManager(options, ws1.All()...)
}

// newWsTestServer 创建WebSocket测试服务
func newWsTestServer(t *testing.T, manager *ws.Manager) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/ws", gin.WrapH(manager))

	server := httptest.NewServer(engine)
	t.Cleanup(server.Close)

	return server
}

// newWsAuthTestServer 创建带登录用户的WebSocket测试服务
func newWsAuthTestServer(t *testing.T, manager *ws.Manager, userID int64) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/ws", func(c *gin.Context) {
		ctx := ws.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)
		manager.ServeHTTP(c.Writer, c.Request)
	})

	server := httptest.NewServer(engine)
	t.Cleanup(server.Close)

	return server
}

// dialWsTestServer 连接WebSocket测试服务
func dialWsTestServer(t *testing.T, serverURL string) *websocket.Conn {
	t.Helper()

	target := "ws" + strings.TrimPrefix(serverURL, "http") + "/ws"
	conn, response, err := websocket.DefaultDialer.Dial(target, nil)
	if response != nil {
		defer func() { _ = response.Body.Close() }()
	}
	require.NoError(t, err)
	require.NotNil(t, conn)

	t.Cleanup(func() { _ = conn.Close() })

	return conn
}

// writeWsMessage 写入WebSocket消息
func writeWsMessage(t *testing.T, conn *websocket.Conn, messageType int, data []byte) {
	t.Helper()
	require.NoError(t, conn.WriteMessage(messageType, data))
}

// readWsMessage 读取WebSocket消息
func readWsMessage(t *testing.T, conn *websocket.Conn) (int, []byte) {
	t.Helper()

	require.NoError(t, conn.SetReadDeadline(time.Now().Add(2*time.Second)))
	messageType, data, err := conn.ReadMessage()
	require.NoError(t, err)

	return messageType, data
}

// TestWebSocketMessage WebSocket测试
func TestWebSocketMessage(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, conn, websocket.TextMessage, []byte(`{"type":"ping"}`))
	messageType, data := readWsMessage(t, conn)
	require.Equal(t, websocket.TextMessage, messageType)

	var response wsResponse
	require.NoError(t, json.Unmarshal(data, &response))
	require.Equal(t, "pong", response.Type)
	require.Equal(t, int64(1), manager.Count())
}

// TestWebSocketSendError WebSocket公共错误消息测试
func TestWebSocketSendError(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, conn, websocket.TextMessage, []byte(`invalid json`))
	messageType, data := readWsMessage(t, conn)
	require.Equal(t, websocket.TextMessage, messageType)

	var response wsResponse
	require.NoError(t, json.Unmarshal(data, &response))
	require.Equal(t, "error", response.Type)

	var errorResponse wsErrorResponse
	require.NoError(t, json.Unmarshal(response.Data, &errorResponse))
	require.Equal(t, int64(400), errorResponse.Code)
	require.Equal(t, "消息格式错误", errorResponse.Message)
}

// TestWebSocketChat WebSocket聊天消息测试
func TestWebSocketChat(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, conn, websocket.TextMessage, []byte(`{"type":"chat","data":"hello"}`))
	messageType, data := readWsMessage(t, conn)
	require.Equal(t, websocket.TextMessage, messageType)

	var response wsResponse
	require.NoError(t, json.Unmarshal(data, &response))
	require.Equal(t, "chat", response.Type)
	require.Equal(t, `"hello"`, string(response.Data))
}

// TestWebSocketBroadcast WebSocket广播测试
func TestWebSocketBroadcast(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	count := manager.Broadcast(ws.Message{
		Type: ws.MessageText,
		Data: []byte("broadcast"),
	})
	require.Equal(t, 1, count)

	messageType, data := readWsMessage(t, conn)
	require.Equal(t, websocket.TextMessage, messageType)
	require.Equal(t, []byte("broadcast"), data)
}

// TestWebSocketBinary WebSocket二进制消息测试
func TestWebSocketBinary(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, conn, websocket.BinaryMessage, []byte(`{"type":"ping"}`))
	messageType, data := readWsMessage(t, conn)
	require.Equal(t, websocket.TextMessage, messageType)

	var response wsResponse
	require.NoError(t, json.Unmarshal(data, &response))
	require.Equal(t, "pong", response.Type)
}

// TestWebSocketSendToUser WebSocket指定用户发送测试
func TestWebSocketSendToUser(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	boundConn := dialWsTestServer(t, server.URL)
	otherConn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, boundConn, websocket.TextMessage, []byte(`{"type":"bind","data":{"userId":1}}`))
	_, data := readWsMessage(t, boundConn)
	var response wsResponse
	require.NoError(t, json.Unmarshal(data, &response))
	require.Equal(t, "bound", response.Type)

	count := manager.SendToUser(1, ws.Message{
		Type: ws.MessageText,
		Data: []byte("user message"),
	})
	require.Equal(t, 1, count)

	_, data = readWsMessage(t, boundConn)
	require.Equal(t, []byte("user message"), data)

	require.NoError(t, otherConn.SetReadDeadline(time.Now().Add(100*time.Millisecond)))
	_, _, err := otherConn.ReadMessage()
	require.Error(t, err)
}

// TestWebSocketAuthenticatedUserBinding WebSocket自动绑定登录用户测试
func TestWebSocketAuthenticatedUserBinding(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 2})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsAuthTestServer(t, manager, 7)
	conn := dialWsTestServer(t, server.URL)
	writeWsMessage(t, conn, websocket.TextMessage, []byte(`{"type":"ping"}`))
	_, _ = readWsMessage(t, conn)

	count := manager.SendToUser(7, ws.Message{
		Type: ws.MessageText,
		Data: []byte("authenticated message"),
	})
	require.Equal(t, 1, count)

	_, data := readWsMessage(t, conn)
	require.Equal(t, []byte("authenticated message"), data)
}

// TestWebSocketSendToUsers WebSocket指定多个用户发送测试
func TestWebSocketSendToUsers(t *testing.T) {
	manager := newWsTestManager(ws.Options{Shards: 1, MaxConnections: 3})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	firstConn := dialWsTestServer(t, server.URL)
	secondConn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, firstConn, websocket.TextMessage, []byte(`{"type":"bind","data":{"userId":1}}`))
	_, _ = readWsMessage(t, firstConn)
	writeWsMessage(t, secondConn, websocket.TextMessage, []byte(`{"type":"bind","data":{"userId":2}}`))
	_, _ = readWsMessage(t, secondConn)

	count := manager.SendToUsers([]int64{1, 2}, ws.Message{
		Type: ws.MessageText,
		Data: []byte("users message"),
	})
	require.Equal(t, 2, count)

	_, data := readWsMessage(t, firstConn)
	require.Equal(t, []byte("users message"), data)
	_, data = readWsMessage(t, secondConn)
	require.Equal(t, []byte("users message"), data)
}

// TestWebSocketMaxConnections WebSocket最大连接数测试
func TestWebSocketMaxConnections(t *testing.T) {
	manager := ws.NewManager(ws.Options{Shards: 1, MaxConnections: 1})
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	target := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	_, response, err := websocket.DefaultDialer.Dial(target, nil)
	if response != nil {
		defer func() { _ = response.Body.Close() }()
	}
	require.Error(t, err)
	require.NotNil(t, response)
	require.Equal(t, http.StatusServiceUnavailable, response.StatusCode)

	require.NoError(t, conn.Close())
	require.Eventually(t, func() bool {
		return manager.Count() == 0
	}, time.Second, 10*time.Millisecond)
}

// TestWebSocketMatchedHandler WebSocket匹配处理器优先级测试
func TestWebSocketMatchedHandler(t *testing.T) {
	manager := ws.NewManager(
		ws.Options{Shards: 1, MaxConnections: 2},
		&wsFallbackTestHandler{},
		&wsMatchedTestHandler{},
	)
	t.Cleanup(func() { _ = manager.Close(context.Background()) })

	server := newWsTestServer(t, manager)
	conn := dialWsTestServer(t, server.URL)

	writeWsMessage(t, conn, websocket.TextMessage, []byte(`{"type":"notice"}`))
	_, data := readWsMessage(t, conn)
	require.Equal(t, []byte("matched"), data)
}
