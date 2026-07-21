package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// Dict 请求验证
type Dict struct {
	base.BaseRequest
	ID     int64       `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Pid    int64       `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name   string      `json:"name" form:"name" validate:"required" label:"标识"`
	Title  string      `json:"title" form:"title" validate:"required" label:"名称"`
	Value  string      `json:"value" form:"value" validate:"" label:"映射值"`
	Status int64       `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort   int64       `json:"sort" form:"sort" validate:"int" label:"排序"`
	Extend interface{} `json:"extend" form:"extend" validate:"" label:"扩展字段"`
	Desc   string      `json:"desc" form:"desc" validate:"" label:"字段描述"`
	PageListValidate
}

// DictCreate 字典创建验证
type DictCreate struct {
	Pid    int64       `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name   string      `json:"name" form:"name" validate:"required" label:"标识"`
	Title  string      `json:"title" form:"title" validate:"required" label:"名称"`
	Value  string      `json:"value" form:"value" validate:"required" label:"映射值"`
	Status int64       `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort   int64       `json:"sort" form:"sort" validate:"int" label:"排序"`
	Extend interface{} `json:"extend" form:"extend" validate:"" label:"扩展字段"`
	Desc   string      `json:"desc" form:"desc" validate:"" label:"字段描述"`
}

// DictUpdate 字典更新验证
type DictUpdate struct {
	Pid    int64       `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name   string      `json:"name" form:"name" validate:"required" label:"标识"`
	Title  string      `json:"title" form:"title" validate:"required" label:"名称"`
	Value  string      `json:"value" form:"value" validate:"required" label:"映射值"`
	Status int64       `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort   int64       `json:"sort" form:"sort" validate:"int" label:"排序"`
	Extend interface{} `json:"extend" form:"extend" validate:"" label:"扩展字段"`
	Desc   string      `json:"desc" form:"desc" validate:"" label:"字段描述"`
}

// Validate 请求验证
func (s Dict) Validate(data Dict, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Dict) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{"Pid", "Name", "Title", "Value", "Status", "Sort", "Extend", "Desc"},
		"Update": []string{"ID", "Pid", "Name", "Title", "Value", "Status", "Sort", "Extend", "Desc"},
		"Detail": []string{"ID"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Dict) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"gt":                           "字段 {field} 必须大于 0",
		"minLen":                       "{field} 长度不能少于 {min} 个字符",
		"maxLen":                       "{field} 长度不能超过 {max} 个字符",
		"PageListValidate.Page.gt":     "页码必须大于 0",
		"PageListValidate.PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s Dict) Translates() map[string]string {
	return validate.MS{
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
