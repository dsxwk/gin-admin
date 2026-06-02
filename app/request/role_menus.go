package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// RoleMenus 请求验证
type RoleMenus struct {
	base.BaseRequest
	ID     int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	RoleId int64  `json:"role_id" form:"role_id" validate:"required|int" label:"角色id"`
	MenuId int64  `json:"menu_id" form:"menu_id" validate:"required|int" label:"菜单id"`
	Name   string `json:"name" form:"name" validate:"required|max:255" label:"角色名称"`
	PageListValidate
}

// Validate 请求验证
func (s RoleMenus) Validate(data RoleMenus, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s RoleMenus) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"list":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"create": []string{"RoleId", "MenuId", "Name"},
		"update": []string{"ID", "RoleId", "MenuId", "Name"},
		"detail": []string{"ID"},
		"delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s RoleMenus) Messages() map[string]string {
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
func (s RoleMenus) Translates() map[string]string {
	return validate.MS{
		"ID":       "ID",
		"RoleId":   "角色id",
		"MenuId":   "菜单id",
		"Name":     "角色名称",
		"Page":     "页码",
		"PageSize": "每页数量",
	}
}
