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

// McpInternalError 内部错误
func McpInternalError() errCode.ErrorCode {
	return errCode.NewError(-32603, "Internal error")
}

// McpResourceAlreadyExists 资源已存在
func McpResourceAlreadyExists() errCode.ErrorCode {
	return errCode.NewError(-32002, "Resource already exists")
}

// McpResourceUnavailable 资源暂时不可用
func McpResourceUnavailable() errCode.ErrorCode {
	return errCode.NewError(-32003, "Resource temporarily unavailable")
}

// McpResourceNotEmpty 资源不为空
func McpResourceNotEmpty() errCode.ErrorCode {
	return errCode.NewError(-32004, "Resource not empty")
}

// McpResourceInvalid 资源无效
func McpResourceInvalid() errCode.ErrorCode {
	return errCode.NewError(-32005, "Resource not valid")
}

// McpSessionNotFound 会话不存在
func McpSessionNotFound() errCode.ErrorCode {
	return errCode.NewError(-32006, "Session not found")
}

// McpRequestTimeout 请求超时
func McpRequestTimeout() errCode.ErrorCode {
	return errCode.NewError(-32007, "Request timeout").WithHttpCode(http.StatusGatewayTimeout)
}

// McpMissingToken 缺少Token
func McpMissingToken() errCode.ErrorCode {
	return errCode.NewError(-32008, "Missing token").WithHttpCode(http.StatusUnauthorized)
}

// McpInvalidToken Token无效
func McpInvalidToken() errCode.ErrorCode {
	return errCode.NewError(-32009, "Invalid token").WithHttpCode(http.StatusUnauthorized)
}

// McpInvalidSession 会话无效
func McpInvalidSession() errCode.ErrorCode {
	return errCode.NewError(-32010, "Invalid session").WithHttpCode(http.StatusUnauthorized)
}

// McpSessionExists 会话已存在
func McpSessionExists() errCode.ErrorCode {
	return errCode.NewError(-32011, "Session already exists")
}

// McpSessionTimeout 会话超时
func McpSessionTimeout() errCode.ErrorCode {
	return errCode.NewError(-32012, "Session timeout").WithHttpCode(http.StatusGatewayTimeout)
}

// McpContentTooLarge 内容过大
func McpContentTooLarge() errCode.ErrorCode {
	return errCode.NewError(-32050, "Content too large").WithHttpCode(http.StatusRequestEntityTooLarge)
}

// McpToolNotFound 工具不存在
func McpToolNotFound() errCode.ErrorCode {
	return errCode.NewError(-32100, "Tool not found")
}

// McpToolExists 工具已存在
func McpToolExists() errCode.ErrorCode {
	return errCode.NewError(-32101, "Tool already exists")
}

// McpToolNotLoaded 工具未加载
func McpToolNotLoaded() errCode.ErrorCode {
	return errCode.NewError(-32102, "Tool not loaded")
}

// McpToolExecutionFailed 工具执行失败
func McpToolExecutionFailed() errCode.ErrorCode {
	return errCode.NewError(-32103, "Tool execution failed")
}

// McpPromptNotFound 提示词不存在
func McpPromptNotFound() errCode.ErrorCode {
	return errCode.NewError(-32104, "Prompt not found")
}

// McpPromptExists 提示词已存在
func McpPromptExists() errCode.ErrorCode {
	return errCode.NewError(-32105, "Prompt already exists")
}

// McpPromptNotLoaded 提示词未加载
func McpPromptNotLoaded() errCode.ErrorCode {
	return errCode.NewError(-32106, "Prompt not loaded")
}

// McpPromptExecutionFailed 提示词执行失败
func McpPromptExecutionFailed() errCode.ErrorCode {
	return errCode.NewError(-32107, "Prompt execution failed")
}

// McpCompletionNotFound 补全不存在
func McpCompletionNotFound() errCode.ErrorCode {
	return errCode.NewError(-32108, "Completion not found")
}

// McpCompletionExists 补全已存在
func McpCompletionExists() errCode.ErrorCode {
	return errCode.NewError(-32109, "Completion already exists")
}

// McpCompletionNotLoaded 补全未加载
func McpCompletionNotLoaded() errCode.ErrorCode {
	return errCode.NewError(-32110, "Completion not loaded")
}

// McpCompletionExecutionFailed 补全执行失败
func McpCompletionExecutionFailed() errCode.ErrorCode {
	return errCode.NewError(-32111, "Completion execution failed")
}

// McpSamplerNotFound 采样器不存在
func McpSamplerNotFound() errCode.ErrorCode {
	return errCode.NewError(-32112, "Sampler not found")
}

// McpSamplerNotLoaded 采样器未加载
func McpSamplerNotLoaded() errCode.ErrorCode {
	return errCode.NewError(-32113, "Sampler not loaded")
}

// McpSamplerExecutionFailed 采样器执行失败
func McpSamplerExecutionFailed() errCode.ErrorCode {
	return errCode.NewError(-32114, "Sampler execution failed")
}

// McpResourceNotFound 资源不存在
func McpResourceNotFound() errCode.ErrorCode {
	return errCode.NewError(-32115, "Resource not found")
}

// McpResourceExists 资源已存在
func McpResourceExists() errCode.ErrorCode {
	return errCode.NewError(-32116, "Resource already exists")
}

// McpResourceNotLoaded 资源未加载
func McpResourceNotLoaded() errCode.ErrorCode {
	return errCode.NewError(-32117, "Resource not loaded")
}

// McpResourceExecutionFailed 资源执行失败
func McpResourceExecutionFailed() errCode.ErrorCode {
	return errCode.NewError(-32118, "Resource execution failed")
}

// McpResourceSubscriptionNotFound 资源订阅不存在
func McpResourceSubscriptionNotFound() errCode.ErrorCode {
	return errCode.NewError(-32119, "Resource subscription not found")
}

// McpResourceSubscriptionExists 资源订阅已存在
func McpResourceSubscriptionExists() errCode.ErrorCode {
	return errCode.NewError(-32120, "Resource subscription already exists")
}

// McpResourceSubscriptionInvalid 资源订阅无效
func McpResourceSubscriptionInvalid() errCode.ErrorCode {
	return errCode.NewError(-32121, "Resource subscription invalid")
}

// McpResourceSubscriptionExecutionFailed 资源订阅执行失败
func McpResourceSubscriptionExecutionFailed() errCode.ErrorCode {
	return errCode.NewError(-32122, "Resource subscription execution failed")
}
