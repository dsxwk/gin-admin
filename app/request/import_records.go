package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// ImportRecords 请求验证
type ImportRecords struct {
	base.Request
	ID   int64  `json:"id" form:"id" label:"ID"`
	Type int64  `json:"type" form:"type" label:"导入类型"`
	Name string `json:"name" form:"name" label:"类型名称"`
	Data any    `json:"data" form:"data" label:"导入数据"`
	PageListValidate
}

// ValidateWithScene 请求验证
func (s ImportRecords) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
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
func (s ImportRecords) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s ImportRecords) Translates() validatex.Translates {
	return validatex.Translates{
		"id":       "ID",
		"type":     "导入类型",
		"name":     "类型名称",
		"data":     "导入数据",
		"page":     "页码",
		"pageSize": "每页数量",
	}
}
