package mcp

import (
	"errors"
	"gin/common/errcode"
	"gin/common/response"
	"gin/config"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
)

// Handler MCP HTTP处理器
type Handler struct {
	Info   ServerInfo
	Config *config.Config
}

// NewHandler 创建MCP处理器
func NewHandler(info ServerInfo, cfg *config.Config) *Handler {
	return &Handler{Info: info, Config: cfg}
}

// ServeHTTP 纯HTTP处理器,处理JSONRPC请求
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 认证检查
	if !h.checkAuth(w, r) {
		return
	}

	// 读取并解析JSONRPC请求
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorCode(w, nil, errcode.McpParseError())
		return
	}

	var req JsonRpcRequest
	if err = json.Unmarshal(body, &req); err != nil {
		h.writeErrorCode(w, nil, errcode.McpParseError())
		return
	}
	if req.Jsonrpc != "2.0" {
		h.writeErrorCode(w, req.Id, errcode.McpInvalidRequest())
		return
	}

	// 路由JSONRPC方法
	switch req.Method {
	case "initialize":
		h.handleInitialize(w, req)
	case "tools/list":
		h.handleToolsList(w, req)
	case "tools/call":
		h.handleToolsCall(w, r, req)
	default:
		h.writeErrorCode(w, req.Id, errcode.McpMethodNotFound().WithMsg("Method not found: "+req.Method))
	}
}

// checkAuth 认证检查
func (h *Handler) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if h.Config == nil || !h.Config.Mcp.Auth.Enabled || h.Config.Mcp.Auth.Token == "" {
		return true
	}

	if bearerToken(r) != h.Config.Mcp.Auth.Token {
		writeHttpError(w, errcode.McpInvalidToken())
		return false
	}
	return true
}

// handleInitialize 处理initialize请求
func (h *Handler) handleInitialize(w http.ResponseWriter, req JsonRpcRequest) {
	h.writeResult(w, req.Id, map[string]any{
		"protocolVersion": "2024-11-05",
		"serverInfo":      h.Info,
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
	})
}

// handleToolsList 处理tools/list请求
func (h *Handler) handleToolsList(w http.ResponseWriter, req JsonRpcRequest) {
	h.writeResult(w, req.Id, map[string]any{
		"tools": GetAllDefs(),
	})
}

// handleToolsCall 处理tools/call请求
func (h *Handler) handleToolsCall(w http.ResponseWriter, r *http.Request, req JsonRpcRequest) {
	name, _ := req.Params["name"].(string)
	if name == "" {
		h.writeErrorCode(w, req.Id, errcode.McpInvalidParams().WithMsg("Missing tool name"))
		return
	}

	tool, ok := Get(name)
	if !ok {
		h.writeErrorCode(w, req.Id, errcode.McpToolNotFound().WithMsg("Tool not found: "+name))
		return
	}

	args, _ := req.Params["arguments"].(map[string]any)
	if args == nil {
		args = make(map[string]any)
	}

	ctx := r.Context()
	userId := int64(0)

	if at, ok := tool.(AuthTool); ok && at.IsAuth() {
		uid, err := h.authUserId(r)
		if err != nil {
			if bearerToken(r) == "" {
				h.writeErrorCode(w, req.Id, errcode.McpMissingToken())
			} else {
				h.writeErrorCode(w, req.Id, errcode.McpInvalidToken())
			}
			return
		}
		if uid == 0 {
			h.writeErrorCode(w, req.Id, errcode.McpInvalidToken())
			return
		}
		userId = uid
	}

	result, err := CallTool(ctx, userId, tool, args)

	if err != nil {
		h.writeResult(w, req.Id, CallResult{
			Content: []ContentItem{{Type: "text", Text: err.Error()}},
			IsError: true,
		})
		return
	}

	resultJson, _ := json.Marshal(result)
	h.writeResult(w, req.Id, CallResult{
		Content: []ContentItem{{Type: "text", Text: string(resultJson)}},
		IsError: false,
	})
}

// writeResult 写入成功响应
func (h *Handler) writeResult(w http.ResponseWriter, id any, result any) {
	writeJson(w, JsonRpcResponse{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  result,
	})
}

// writeErrorCode 写入JSONRPC错误响应
func (h *Handler) writeErrorCode(w http.ResponseWriter, id any, e errcode.ErrorCode) {
	writeJson(w, JsonRpcResponse{
		Jsonrpc: "2.0",
		Id:      id,
		Error: &JsonRpcError{
			Code:    int(e.Code),
			Message: e.Msg,
		},
	})
}

// writeJson 写入JSON响应
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

// writeHttpError 写入HTTP级错误,格式与controller响应一致
func writeHttpError(w http.ResponseWriter, e errcode.ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	if e.HttpCode == 0 {
		e.HttpCode = http.StatusInternalServerError
	}
	w.WriteHeader(e.HttpCode)
	data := any([]string{})
	if e.Data != nil {
		data = e.Data
	}
	_ = json.NewEncoder(w).Encode(response.Response{Code: e.Code, Msg: e.Msg, Data: data})
}

// bearerToken 从请求头解析Token
func bearerToken(r *http.Request) string {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		token = r.Header.Get("token")
	}
	return token
}

// authUserId 从请求头解析JWT获取用户ID
func (h *Handler) authUserId(r *http.Request) (int64, error) {
	if h.Config == nil {
		return 0, errors.New("mcp config is nil")
	}

	token := bearerToken(r)
	if token == "" {
		return 0, errors.New("missing token")
	}

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unsupported signing method")
		}
		return []byte(h.Config.Jwt.Key), nil
	})
	if err != nil || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid user id")
	}
	return int64(id), nil
}
