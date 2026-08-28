package errcode

import "net/http"

// McpParseError JSONRPC解析错误
func McpParseError() ErrorCode {
	return NewError(-32700, "Parse error").WithHttpCode(http.StatusBadRequest)
}

// McpInvalidRequest JSONRPC无效请求
func McpInvalidRequest() ErrorCode {
	return NewError(-32600, "Invalid Request").WithHttpCode(http.StatusBadRequest)
}

// McpMethodNotFound JSONRPC方法不存在
func McpMethodNotFound() ErrorCode {
	return NewError(-32601, "Method not found")
}

// McpInvalidParams JSONRPC参数无效
func McpInvalidParams() ErrorCode {
	return NewError(-32602, "Invalid params")
}

// McpInternalError JSONRPC内部错误
func McpInternalError() ErrorCode {
	return NewError(-32603, "Internal error")
}

// McpResourceAlreadyExists 资源已存在
func McpResourceAlreadyExists() ErrorCode {
	return NewError(-32002, "Resource already exists")
}

// McpResourceUnavailable 资源暂时不可用
func McpResourceUnavailable() ErrorCode {
	return NewError(-32003, "Resource temporarily unavailable")
}

// McpResourceNotEmpty 资源不为空
func McpResourceNotEmpty() ErrorCode {
	return NewError(-32004, "Resource not empty")
}

// McpResourceInvalid 资源无效
func McpResourceInvalid() ErrorCode {
	return NewError(-32005, "Resource not valid")
}

// McpSessionNotFound 会话不存在
func McpSessionNotFound() ErrorCode {
	return NewError(-32006, "Session not found")
}

// McpRequestTimeout 请求超时
func McpRequestTimeout() ErrorCode {
	return NewError(-32007, "Request timeout").WithHttpCode(http.StatusGatewayTimeout)
}

// McpMissingToken 缺少Token
func McpMissingToken() ErrorCode {
	return NewError(-32008, "Missing token").WithHttpCode(http.StatusUnauthorized)
}

// McpInvalidToken Token无效
func McpInvalidToken() ErrorCode {
	return NewError(-32009, "Invalid token").WithHttpCode(http.StatusUnauthorized)
}

// McpInvalidSession 会话无效
func McpInvalidSession() ErrorCode {
	return NewError(-32010, "Invalid session").WithHttpCode(http.StatusUnauthorized)
}

// McpSessionExists 会话已存在
func McpSessionExists() ErrorCode {
	return NewError(-32011, "Session already exists")
}

// McpSessionTimeout 会话超时
func McpSessionTimeout() ErrorCode {
	return NewError(-32012, "Session timeout").WithHttpCode(http.StatusGatewayTimeout)
}

// McpContentTooLarge 内容过大
func McpContentTooLarge() ErrorCode {
	return NewError(-32050, "Content too large").WithHttpCode(http.StatusRequestEntityTooLarge)
}

// McpToolNotFound 工具不存在
func McpToolNotFound() ErrorCode {
	return NewError(-32100, "Tool not found")
}

// McpToolExists 工具已存在
func McpToolExists() ErrorCode {
	return NewError(-32101, "Tool already exists")
}

// McpToolNotLoaded 工具未加载
func McpToolNotLoaded() ErrorCode {
	return NewError(-32102, "Tool not loaded")
}

// McpToolExecutionFailed 工具执行失败
func McpToolExecutionFailed() ErrorCode {
	return NewError(-32103, "Tool execution failed")
}

// McpPromptNotFound 提示词不存在
func McpPromptNotFound() ErrorCode {
	return NewError(-32104, "Prompt not found")
}

// McpPromptExists 提示词已存在
func McpPromptExists() ErrorCode {
	return NewError(-32105, "Prompt already exists")
}

// McpPromptNotLoaded 提示词未加载
func McpPromptNotLoaded() ErrorCode {
	return NewError(-32106, "Prompt not loaded")
}

// McpPromptExecutionFailed 提示词执行失败
func McpPromptExecutionFailed() ErrorCode {
	return NewError(-32107, "Prompt execution failed")
}

// McpCompletionNotFound 补全不存在
func McpCompletionNotFound() ErrorCode {
	return NewError(-32108, "Completion not found")
}

// McpCompletionExists 补全已存在
func McpCompletionExists() ErrorCode {
	return NewError(-32109, "Completion already exists")
}

// McpCompletionNotLoaded 补全未加载
func McpCompletionNotLoaded() ErrorCode {
	return NewError(-32110, "Completion not loaded")
}

// McpCompletionExecutionFailed 补全执行失败
func McpCompletionExecutionFailed() ErrorCode {
	return NewError(-32111, "Completion execution failed")
}

// McpSamplerNotFound 采样器不存在
func McpSamplerNotFound() ErrorCode {
	return NewError(-32112, "Sampler not found")
}

// McpSamplerNotLoaded 采样器未加载
func McpSamplerNotLoaded() ErrorCode {
	return NewError(-32113, "Sampler not loaded")
}

// McpSamplerExecutionFailed 采样器执行失败
func McpSamplerExecutionFailed() ErrorCode {
	return NewError(-32114, "Sampler execution failed")
}

// McpResourceNotFound 资源不存在
func McpResourceNotFound() ErrorCode {
	return NewError(-32115, "Resource not found")
}

// McpResourceExists 资源已存在
func McpResourceExists() ErrorCode {
	return NewError(-32116, "Resource already exists")
}

// McpResourceNotLoaded 资源未加载
func McpResourceNotLoaded() ErrorCode {
	return NewError(-32117, "Resource not loaded")
}

// McpResourceExecutionFailed 资源执行失败
func McpResourceExecutionFailed() ErrorCode {
	return NewError(-32118, "Resource execution failed")
}

// McpResourceSubscriptionNotFound 资源订阅不存在
func McpResourceSubscriptionNotFound() ErrorCode {
	return NewError(-32119, "Resource subscription not found")
}

// McpResourceSubscriptionExists 资源订阅已存在
func McpResourceSubscriptionExists() ErrorCode {
	return NewError(-32120, "Resource subscription already exists")
}

// McpResourceSubscriptionInvalid 资源订阅无效
func McpResourceSubscriptionInvalid() ErrorCode {
	return NewError(-32121, "Resource subscription invalid")
}

// McpResourceSubscriptionExecutionFailed 资源订阅执行失败
func McpResourceSubscriptionExecutionFailed() ErrorCode {
	return NewError(-32122, "Resource subscription execution failed")
}
