package errcode

import (
	errCode "gin/pkg/errcode"
	"net/http"
)

// McpParseError 解析错误
func McpParseError() errCode.ErrorCode {
	return errCode.NewError(-32700, "Parse error").WithHttpCode(http.StatusBadRequest)
}

// McpInvalidRequest 无效请求
func McpInvalidRequest() errCode.ErrorCode {
	return errCode.NewError(-32600, "Invalid Request").WithHttpCode(http.StatusBadRequest)
}

// McpMethodNotFound 方法不存在
func McpMethodNotFound() errCode.ErrorCode {
	return errCode.NewError(-32601, "Method not found")
}

// McpInvalidParams 参数无效
func McpInvalidParams() errCode.ErrorCode {
	return errCode.NewError(-32602, "Invalid params")
}

// McpMissingToken 缺少Token
func McpMissingToken() errCode.ErrorCode {
	return errCode.NewError(-32008, "Missing token").WithHttpCode(http.StatusUnauthorized)
}

// McpInvalidToken Token无效
func McpInvalidToken() errCode.ErrorCode {
	return errCode.NewError(-32009, "Invalid token").WithHttpCode(http.StatusUnauthorized)
}

// McpToolNotFound 工具不存在
func McpToolNotFound() errCode.ErrorCode {
	return errCode.NewError(-32100, "Tool not found")
}
