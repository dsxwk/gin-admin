package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// Roles 请求验证
type Roles struct {
	base.BaseRequest
	ID     int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Name   string `json:"name" form:"name" validate:"required|max:255" label:"角色名称"`
	Desc   string `json:"desc" form:"desc" validate:"required|max:255" label:"角色描述"`
	Status int64  `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	PageListValidate
}

// Validate 请求验证
func (s Roles) Validate(data Roles, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Roles) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"list":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"create": []string{"Name", "Desc", "Status"},
		"update": []string{"ID", "Name", "Desc", "Status"},
		"detail": []string{"ID"},
		"delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Roles) Messages() map[string]string {
	return validate.MS{
		"required":    "字段 {field} 必填",
		"int":         "字段 {field} 必须为整数",
		"gt":          "字段 {field} 必须大于 0",
		"max":         "字段 {field} 长度不能超过 255",
		"Page.gt":     "页码必须大于 0",
		"PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s Roles) Translates() map[string]string {
	return validate.MS{
		"ID":       "ID",
		"Name":     "角色名称",
		"Desc":     "角色描述",
		"Status":   "状态 1=启用 2=停用",
		"Page":     "页码",
		"PageSize": "每页数量",
	}
}
