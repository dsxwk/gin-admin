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
			fields.Str("Name", s.Name).Required().Errors(),
			fields.Str("Title", s.Title).Required().Errors(),
			fields.Str("Value", s.Value).Required().Errors(),
			fields.Int("Status", s.Status).Required().Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("ID", s.ID).Required().Gt(0).Errors(),
			fields.Str("Name", s.Name).Required().Errors(),
			fields.Str("Title", s.Title).Required().Errors(),
			fields.Str("Value", s.Value).Required().Errors(),
			fields.Int("Status", s.Status).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
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
		"ID":                        "ID",
		"Pid":                       "父级id",
		"Name":                      "标识",
		"Title":                     "名称",
		"Value":                     "映射值",
		"Status":                    "状态 1=启用 2=停用",
		"Sort":                      "排序",
		"Extend":                    "扩展字段",
		"Desc":                      "字段描述",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
