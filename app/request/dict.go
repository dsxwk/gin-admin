package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// Dict 请求验证
type Dict struct {
	base.Request
	ID     int64  `json:"id" form:"id" label:"ID"`
	Pid    int64  `json:"pid" form:"pid" label:"父级id"`
	Name   string `json:"name" form:"name" label:"标识"`
	Title  string `json:"title" form:"title" label:"名称"`
	Value  string `json:"value" form:"value" label:"映射值"`
	Status int64  `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" label:"排序"`
	Extend any    `json:"extend" form:"extend" label:"扩展字段"`
	Desc   string `json:"desc" form:"desc" label:"字段描述"`
	PageListValidate
}

// DictCreate 字典创建验证
type DictCreate struct {
	Pid    int64  `json:"pid" form:"pid" label:"父级id"`
	Name   string `json:"name" form:"name" label:"标识"`
	Title  string `json:"title" form:"title" label:"名称"`
	Value  string `json:"value" form:"value" label:"映射值"`
	Status int64  `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" label:"排序"`
	Extend any    `json:"extend" form:"extend" label:"扩展字段"`
	Desc   string `json:"desc" form:"desc" label:"字段描述"`
}

// DictUpdate 字典更新验证
type DictUpdate struct {
	Pid    int64  `json:"pid" form:"pid" label:"父级id"`
	Name   string `json:"name" form:"name" label:"标识"`
	Title  string `json:"title" form:"title" label:"名称"`
	Value  string `json:"value" form:"value" label:"映射值"`
	Status int64  `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort   int64  `json:"sort" form:"sort" label:"排序"`
	Extend any    `json:"extend" form:"extend" label:"扩展字段"`
	Desc   string `json:"desc" form:"desc" label:"字段描述"`
}

// ValidateWithScene 请求验证
func (s Dict) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return validatex.Check(
			fields.Str("name", s.Name).Required().Errors(),
			fields.Str("title", s.Title).Required().Errors(),
			fields.Str("value", s.Value).Required().Errors(),
			fields.Int("status", s.Status).Required().Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Str("name", s.Name).Required().Errors(),
			fields.Str("title", s.Title).Required().Errors(),
			fields.Str("value", s.Value).Required().Errors(),
			fields.Int("status", s.Status).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Dict) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s Dict) Translates() validatex.Translates {
	return validatex.Translates{
		"id":       "ID",
		"pid":      "父级id",
		"name":     "标识",
		"title":    "名称",
		"value":    "映射值",
		"status":   "状态 1=启用 2=停用",
		"sort":     "排序",
		"extend":   "扩展字段",
		"desc":     "字段描述",
		"page":     "页码",
		"pageSize": "每页数量",
	}
}
