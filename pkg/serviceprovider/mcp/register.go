package mcp

import (
	"strings"
	"sync"
)

var (
	tools   = make(map[string]Tool)
	toolsMu sync.RWMutex
)

// Register 注册MCP工具
func Register(t Tool) {
	toolsMu.Lock()
	defer toolsMu.Unlock()
	tools[t.Name()] = t
}

// GetAll 获取所有已注册工具
func GetAll() []Tool {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	result := make([]Tool, 0, len(tools))
	for _, t := range tools {
		result = append(result, t)
	}
	return result
}

// Get 获取指定工具
func Get(name string) (Tool, bool) {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	t, ok := tools[name]
	return t, ok
}

// GetAllDefs 获取所有工具定义(用于tools/list)
func GetAllDefs() []ToolDef {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	defs := make([]ToolDef, 0, len(tools))
	for _, t := range tools {
		defs = append(defs, ToolDef{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	return defs
}

// BuildSystemPrompt 根据已注册工具自动生成系统提示词
func BuildSystemPrompt(role string) string {
	toolsMu.RLock()
	defer toolsMu.RUnlock()

	if len(tools) == 0 {
		return "你是一个" + role + ",当前没有可用工具"
	}

	prompt := "你是" + role + ",只能使用以下MCP工具处理请求:\n"
	for _, t := range tools {
		prompt += "- " + t.Name() + ": " + t.Description() + "\n"
	}
	prompt += "\n规则:\n"
	prompt += "- 仅处理与上述工具功能相关的请求\n"
	prompt += "- 超出范围的问题,直接回复\"抱歉,我当前只能处理以下范围:"

	names := make([]string, 0, len(tools))
	for _, t := range tools {
		names = append(names, shortDesc(t.Description()))
	}
	prompt += strings.Join(names, "、") + "\"\n"
	prompt += "- 不要编造不存在的数据,一切以工具返回为准"

	return prompt
}

// shortDesc 提取简短中文描述(取逗号前部分)
func shortDesc(desc string) string {
	if idx := strings.Index(desc, ","); idx > 0 {
		return desc[:idx]
	}
	return desc
}
