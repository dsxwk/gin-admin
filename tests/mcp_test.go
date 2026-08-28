package tests

import (
	"bytes"
	"encoding/json"
	"gin/config"
	"gin/pkg/serviceprovider/mcp"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestHandler() *mcp.Handler {
	return mcp.NewHandler(mcp.ServerInfo{Name: "test", Version: "1.0.0"}, &config.Config{
		Mcp: config.Mcp{Auth: config.McpAuth{Enabled: true, Token: "secret"}},
		Jwt: config.Jwt{Key: "jwt-key"},
	})
}

func postJson(t *testing.T, h *mcp.Handler, token string, body any) *httptest.ResponseRecorder {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) mcp.JsonRpcResponse {
	var resp mcp.JsonRpcResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp
}

// TestMcpInitialize 测试初始化
func TestMcpInitialize(t *testing.T) {
	h := newTestHandler()
	w := postJson(t, h, "secret", mcp.JsonRpcRequest{Jsonrpc: "2.0", Id: 1, Method: "initialize"})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.Nil(t, resp.Error)
	assert.NotNil(t, resp.Result)
}

// TestMcpToolsList 测试工具列表
func TestMcpToolsList(t *testing.T) {
	h := newTestHandler()
	w := postJson(t, h, "secret", mcp.JsonRpcRequest{Jsonrpc: "2.0", Id: 2, Method: "tools/list"})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.Nil(t, resp.Error)
	assert.NotNil(t, resp.Result)
}

// TestMcpAuthInvalid 测试鉴权失败
func TestMcpAuthInvalid(t *testing.T) {
	h := newTestHandler()
	w := postJson(t, h, "wrong", mcp.JsonRpcRequest{Jsonrpc: "2.0", Id: 3, Method: "initialize"})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, float64(-32009), resp["code"])
}

// TestMcpMethodNotFound 测试方法不存在
func TestMcpMethodNotFound(t *testing.T) {
	h := newTestHandler()
	w := postJson(t, h, "secret", mcp.JsonRpcRequest{Jsonrpc: "2.0", Id: 4, Method: "unknown"})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, -32601, resp.Error.Code)
}

// TestMcpInvalidRequest 测试JSONRPC版本无效
func TestMcpInvalidRequest(t *testing.T) {
	h := newTestHandler()
	w := postJson(t, h, "secret", mcp.JsonRpcRequest{Jsonrpc: "1.0", Id: 5, Method: "initialize"})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, -32600, resp.Error.Code)
}

// TestMcpToolNotFound 测试工具不存在
func TestMcpToolNotFound(t *testing.T) {
	h := newTestHandler()
	w := postJson(t, h, "secret", mcp.JsonRpcRequest{
		Jsonrpc: "2.0",
		Id:      6,
		Method:  "tools/call",
		Params:  map[string]any{"name": "not_exist"},
	})

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, -32100, resp.Error.Code)
}
