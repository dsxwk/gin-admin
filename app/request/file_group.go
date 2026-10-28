package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// FileGroup 文件分组请求验证
type FileGroup struct {
	base.Request
	ID   int64  `json:"id" form:"id" label:"ID"`
	Pid  int64  `json:"pid" form:"pid" label:"父级ID"`
	Name string `json:"name" form:"name" label:"分组名称"`
	Sort int64  `json:"sort" form:"sort" label:"排序"`
	PageListValidate
}

// ValidateWithScene 请求验证
func (s FileGroup) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return fields.Str("name", s.Name).Required().Max(20).Errors()
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Str("name", s.Name).Required().Max(20).Errors(),
		)
	case "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s FileGroup) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s FileGroup) Translates() validatex.Translates {
	return validatex.Translates{
		"id":       "ID",
		"pid":      "父级ID",
		"name":     "分组名称",
		"sort":     "排序",
		"page":     "页码",
		"pageSize": "每页数量",
	}
}
