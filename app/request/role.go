package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// Roles 请求验证
type Roles struct {
	base.BaseRequest
	ID        int64       `uri:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Name      string      `json:"name" form:"name" validate:"required|minLen:1|maxLen:20" label:"角色名称"`
	Desc      string      `json:"desc" form:"desc" validate:"maxLen:100" label:"角色描述"`
	Status    int64       `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	UserRoles []UserRoles `json:"UserRoles" form:"UserRoles" validate:"" label:"用户角色"`
	RoleMenus []RoleMenu  `json:"roleMenus" form:"roleMenus" validate:"" label:"角色菜单"`
	PageListValidate
}

// RoleCreate 角色创建验证
type RoleCreate struct {
	Name      string      `json:"name" validate:"required|minLen:1|maxLen:20" label:"角色名称"`
	Desc      string      `json:"desc" validate:"maxLen:100" label:"角色描述"`
	Status    int64       `json:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	UserRoles []UserRoles `json:"UserRoles" form:"UserRoles" validate:"" label:"用户角色"`
	RoleMenus []RoleMenu  `json:"roleMenus" form:"roleMenus" validate:"" label:"角色菜单"`
}

// UserRoles 用户角色
type UserRoles struct {
	RoleId int64  `json:"roleId" validate:"required|int|gt:0" label:"角色ID"`
	UserId int64  `json:"userId" validate:"required|int|gt:0" label:"用户ID"`
	Name   string `json:"name" validate:"required|minLen:1|maxLen:20" label:"角色名称"`
}

// RoleUpdate 角色更新验证
type RoleUpdate struct {
	ID        int64       `uri:"id" validate:"required|int|gt:0" label:"ID"`
	Name      string      `json:"name" validate:"required|minLen:1|maxLen:20" label:"角色名称"`
	Desc      string      `json:"desc" validate:"maxLen:100" label:"角色描述"`
	Status    int64       `json:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	UserRoles []UserRoles `json:"userRoles" form:"userRoles" validate:"" label:"用户角色"`
	RoleMenus []RoleMenu  `json:"roleMenus" form:"roleMenus" validate:"" label:"角色菜单"`
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
		"List":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{"Name", "Desc", "Status"},
		"Update": []string{"ID", "Name", "Desc", "Status"},
		"Detail": []string{"ID"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Roles) Messages() map[string]string {
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
func (s Roles) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Name":                      "角色名称",
		"Desc":                      "角色描述",
		"Status":                    "状态 1=启用 2=停用",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
