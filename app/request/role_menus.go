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
			fields.Int("role_id", s.RoleId).Required().Errors(),
			fields.Int("menu_id", s.MenuId).Required().Errors(),
			fields.Str("name", s.Name).Required().Max(255).Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Int("role_id", s.RoleId).Required().Errors(),
			fields.Int("menu_id", s.MenuId).Required().Errors(),
			fields.Str("name", s.Name).Required().Max(255).Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
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
		"id":       "ID",
		"role_id":  "角色id",
		"menu_id":  "菜单id",
		"name":     "角色名称",
		"page":     "页码",
		"pageSize": "每页数量",
	}
}
