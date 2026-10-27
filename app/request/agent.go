package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// AgentAsk AI对话请求验证
type AgentAsk struct {
	base.Request
	Question  string `json:"question" form:"question" label:"问题"`
	Provider  string `json:"provider" form:"provider" label:"模型提供商"`
	SessionId int64  `json:"sessionId" form:"sessionId" label:"会话ID"`
}

// ValidateWithScene 请求验证
func (s AgentAsk) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Ask":
		return fields.Str("Question", s.Question).Required().Errors()
	case "Stream":
		return validatex.Check(
			fields.Str("Question", s.Question).Required().Errors(),
			fields.Int("SessionId", s.SessionId).Gte(0).Errors(),
		)
	case "History":
		return fields.Int("SessionId", s.SessionId).Gte(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s AgentAsk) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s AgentAsk) Translates() validatex.Translates {
	return validatex.Translates{
		"Question":  "问题",
		"Provider":  "模型提供商",
		"SessionId": "会话ID",
	}
}

// AgentDelete 删除AI会话请求验证
type AgentDelete struct {
	base.Request
	ID int64 `json:"id" form:"id" label:"会话ID"`
}

// ValidateWithScene 请求验证
func (s AgentDelete) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "Delete":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s AgentDelete) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s AgentDelete) Translates() validatex.Translates {
	return validatex.Translates{"ID": "会话ID"}
}
