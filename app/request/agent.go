package request

import (
	"gin/app/errcode"
	"gin/common/base"

	"github.com/gookit/validate"
)

// AgentAsk AI对话请求验证
type AgentAsk struct {
	base.BaseRequest
	Question  string `json:"question" form:"question" validate:"required" label:"问题"`
	Provider  string `json:"provider" form:"provider" validate:"" label:"模型提供商"`
	SessionId int64  `json:"sessionId" form:"sessionId" validate:"int|min:0" label:"会话ID"`
}

// Validate 请求验证
func (s AgentAsk) Validate(data AgentAsk, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证场景
func (s AgentAsk) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"Ask":     []string{"Question"},
		"Stream":  []string{"Question", "Provider", "SessionId"},
		"History": []string{"SessionId"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s AgentAsk) Messages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
		"min":      "字段 {field} 不能小于 0",
	}
}

// Translates 字段翻译
func (s AgentAsk) Translates() map[string]string {
	return validate.MS{
		"Question":  "问题",
		"Provider":  "模型提供商",
		"SessionId": "会话ID",
	}
}

// AgentDelete 删除AI会话请求验证
type AgentDelete struct {
	base.BaseRequest
	ID int64 `json:"id" form:"id" validate:"required|int|gt:0" label:"会话ID"`
}

// Validate 请求验证
func (s AgentDelete) Validate(data AgentDelete, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证场景
func (s AgentDelete) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"Delete": []string{"ID"},
	})
}

// Messages 验证器错误消息
func (s AgentDelete) Messages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
		"gt":       "字段 {field} 必须大于 0",
	}
}

// Translates 字段翻译
func (s AgentDelete) Translates() map[string]string {
	return validate.MS{
		"ID": "会话ID",
	}
}
