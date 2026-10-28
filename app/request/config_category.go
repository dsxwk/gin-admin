package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// ConfigCategory 请求验证
type ConfigCategory struct {
	base.Request
	ID   int64  `json:"id" form:"id" label:"ID"`
	Name string `json:"name" form:"name" label:"分类名称"`
	PageListValidate
}

// ConfigCategoryCreate ConfigCategory创建验证
type ConfigCategoryCreate struct {
	Name string `json:"name" form:"name" label:"分类名称"`
}

// ConfigCategoryUpdate ConfigCategory更新验证
type ConfigCategoryUpdate struct {
	Name string `json:"name" form:"name" label:"分类名称"`
}

// ValidateWithScene 请求验证
func (s ConfigCategory) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return fields.Str("name", s.Name).Required().Errors()
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Str("name", s.Name).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s ConfigCategory) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s ConfigCategory) Translates() validatex.Translates {
	return validatex.Translates{
		"id":       "ID",
		"name":     "分类名称",
		"page":     "页码",
		"pageSize": "每页数量",
	}
}
