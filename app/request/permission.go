package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// Permission 请求验证
type Permission struct {
	base.Request
	ID     int64  `json:"id" form:"id" label:"ID"`
	Key    string `json:"key" form:"key" label:"权限标识"`
	Method string `json:"method" form:"method" label:"请求方式"`
	Uri    string `json:"uri" form:"uri" label:"路由地址"`
	PageListValidate
}

// ValidateWithScene 请求验证
func (s Permission) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Permission) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s Permission) Translates() validatex.Translates {
	return validatex.Translates{
		"ID":                        "ID",
		"Key":                       "权限标识",
		"Method":                    "请求方式",
		"Uri":                       "路由地址",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
