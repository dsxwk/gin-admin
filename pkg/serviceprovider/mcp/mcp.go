package mcp

import (
	"errors"
	"gin/common/errcode"
	"gin/common/response"
	"gin/config"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"strings"
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

// ServeHTTP Gin适配器,处理JSONRPC请求
func (h *Handler) ServeHTTP(c *gin.Context) {
	// 认证检查
	if !h.checkAuth(c) {
		return
	}

	// 读取并解析JSONRPC请求
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.writeError(c, nil, -32700, "Parse error")
		return
	}

	var req JsonRpcRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.writeError(c, nil, -32700, "Parse error")
		return
	}

	// 路由JSONRPC方法
	switch req.Method {
	case "initialize":
		h.handleInitialize(c, req)
	case "tools/list":
		h.handleToolsList(c, req)
	case "tools/call":
		h.handleToolsCall(c, req)
	default:
		h.writeError(c, req.Id, -32601, "Method not found: "+req.Method)
	}
}

// checkAuth 认证检查
func (h *Handler) checkAuth(c *gin.Context) bool {
	if h.Config == nil || !h.Config.Mcp.Auth.Enabled || h.Config.Mcp.Auth.Token == "" {
		return true
	}

	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token != h.Config.Mcp.Auth.Token {
		resp := response.Response{}
		resp.Error(c, errcode.Forbidden().WithMsg("Invalid token"))
		return false
	}
	return true
}

// handleInitialize 处理initialize请求
func (h *Handler) handleInitialize(c *gin.Context, req JsonRpcRequest) {
	h.writeResult(c, req.Id, map[string]any{
		"protocolVersion": "2024-11-05",
		"serverInfo":      h.Info,
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
	})
}

// handleToolsList 处理tools/list请求
func (h *Handler) handleToolsList(c *gin.Context, req JsonRpcRequest) {
	h.writeResult(c, req.Id, map[string]any{
		"tools": GetAllDefs(),
	})
}

// handleToolsCall 处理tools/call请求
func (h *Handler) handleToolsCall(c *gin.Context, req JsonRpcRequest) {
	name, _ := req.Params["name"].(string)
	if name == "" {
		h.writeError(c, req.Id, -32602, "Missing tool name")
		return
	}

	tool, ok := Get(name)
	if !ok {
		h.writeError(c, req.Id, -32602, "Tool not found: "+name)
		return
	}

	args, _ := req.Params["arguments"].(map[string]any)
	if args == nil {
		args = make(map[string]any)
	}

	ctx := c.Request.Context()
	userId := int64(0)

	if at, ok := tool.(AuthTool); ok && at.IsAuth() {
		uid, err := h.authUserId(c)
		if err != nil || uid == 0 {
			h.writeError(c, req.Id, -32001, "无权限")
			return
		}
		userId = uid
	}

	var (
		result any
		err    error
	)
	if uc, ok := tool.(UserContextCall); ok {
		result, err = uc.CallWithUser(ctx, userId, args)
	} else {
		result, err = tool.Call(ctx, args)
	}

	if err != nil {
		h.writeResult(c, req.Id, CallResult{
			Content: []ContentItem{{Type: "text", Text: err.Error()}},
			IsError: true,
		})
		return
	}

	resultJson, _ := json.Marshal(result)
	h.writeResult(c, req.Id, CallResult{
		Content: []ContentItem{{Type: "text", Text: string(resultJson)}},
		IsError: false,
	})
}

// writeResult 写入成功响应
func (h *Handler) writeResult(c *gin.Context, id any, result any) {
	c.JSON(200, JsonRpcResponse{
		Jsonrpc: "2.0",
		Id:      id,
		Result:  result,
	})
}

// writeError 写入错误响应
func (h *Handler) writeError(c *gin.Context, id any, code int, message string) {
	c.JSON(200, JsonRpcResponse{
		Jsonrpc: "2.0",
		Id:      id,
		Error: &JsonRpcError{
			Code:    code,
			Message: message,
		},
	})
}

// authUserId 从请求头解析JWT获取用户ID
func (h *Handler) authUserId(c *gin.Context) (int64, error) {
	if h.Config == nil {
		return 0, errors.New("mcp config is nil")
	}

	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		token = c.GetHeader("token")
	}
	if token == "" {
		return 0, errors.New("missing token")
	}

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
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
