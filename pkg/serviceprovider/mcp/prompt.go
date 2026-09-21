package mcp

import "strings"

// BuildSystemPrompt 根据工具列表生成系统提示词
func BuildSystemPrompt(role string, tools []Tool) string {
	if len(tools) == 0 {
		return "你是一个" + role + ",当前没有可用工具"
	}

	var prompt strings.Builder
	prompt.WriteString("你是" + role + ",只能使用以下MCP工具处理请求:\n")
	for _, tool := range tools {
		prompt.WriteString("- " + tool.Name() + ": " + tool.Description() + "\n")
	}
	prompt.WriteString("\n规则:\n")
	prompt.WriteString("- 仅处理与上述工具功能相关的请求\n")
	prompt.WriteString("- 超出范围的问题,直接回复\"抱歉,我当前只能处理以下范围:")

	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, shortDesc(tool.Description()))
	}
	prompt.WriteString(strings.Join(names, "、") + "\"\n")
	prompt.WriteString("- 不要编造不存在的数据,一切以工具返回为准")

	return prompt.String()
}

// shortDesc 提取简短中文描述
func shortDesc(desc string) string {
	if idx := strings.Index(desc, ","); idx > 0 {
		return desc[:idx]
	}
	return desc
}
