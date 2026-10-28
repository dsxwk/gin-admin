package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// File 文件请求验证
type File struct {
	base.Request
	ID      int64  `json:"id" form:"id" label:"ID"`
	GroupId *int64 `json:"groupId" form:"groupId" label:"分组ID"`
	PageListValidate
}

// FileUpdateGroup 修改文件分组请求验证
type FileUpdateGroup struct {
	base.Request
	ID      int64  `json:"id" form:"id" label:"ID"`
	GroupId *int64 `json:"groupId" form:"groupId" label:"分组ID"`
}

// ValidateWithScene 请求验证
func (s File) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s File) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s File) Translates() validatex.Translates {
	return validatex.Translates{
		"id":       "ID",
		"groupId":  "分组ID",
		"page":     "页码",
		"pageSize": "每页数量",
	}
}

// ValidateWithScene 请求验证
func (s FileUpdateGroup) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "UpdateGroup":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s FileUpdateGroup) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s FileUpdateGroup) Translates() validatex.Translates {
	return validatex.Translates{
		"id":      "ID",
		"groupId": "分组ID",
	}
}
