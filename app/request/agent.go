package request

import (
	"gin/app/errcode"
	"gin/common/base"

	"github.com/gookit/validate"
)

// AgentAsk AI对话请求验证
type AgentAsk struct {
	base.BaseRequest
	Question  string `json:"question" validate:"required" label:"问题"`
	Provider  string `json:"provider" validate:"" label:"模型提供商"`
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
