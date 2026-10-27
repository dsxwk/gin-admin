package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// RoleMenus 请求验证
type RoleMenus struct {
	base.Request
	ID     int64  `uri:"id" form:"id" label:"ID"`
	RoleId int64  `json:"role_id" form:"role_id" label:"角色id"`
	MenuId int64  `json:"menu_id" form:"menu_id" label:"菜单id"`
	Name   string `json:"name" form:"name" label:"角色名称"`
	PageListValidate
}

// ValidateWithScene 请求验证
func (s RoleMenus) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return validatex.Check(
			fields.Int("RoleId", s.RoleId).Required().Errors(),
			fields.Int("MenuId", s.MenuId).Required().Errors(),
			fields.Str("Name", s.Name).Required().Max(255).Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("ID", s.ID).Required().Gt(0).Errors(),
			fields.Int("RoleId", s.RoleId).Required().Errors(),
			fields.Int("MenuId", s.MenuId).Required().Errors(),
			fields.Str("Name", s.Name).Required().Max(255).Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s RoleMenus) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s RoleMenus) Translates() validatex.Translates {
	return validatex.Translates{
		"ID":                        "ID",
		"RoleId":                    "角色id",
		"MenuId":                    "菜单id",
		"Name":                      "角色名称",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
