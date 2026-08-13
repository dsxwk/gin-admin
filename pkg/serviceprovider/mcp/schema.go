package mcp

// InputSchema 工具输入参数Schema
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property 参数属性定义
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ToolDef 工具定义(用于返回给客户端)
type ToolDef struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

// ServerInfo 服务信息
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// CallRequest 工具调用请求
type CallRequest struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

// CallResult 工具调用结果
type CallResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError"`
}

// ContentItem 返回内容项
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// JsonRpcRequest JSONRPC请求
type JsonRpcRequest struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

// JsonRpcResponse JSONRPC响应
type JsonRpcResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      any           `json:"id"`
	Result  any           `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

// JsonRpcError JSONRPC错误
type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
