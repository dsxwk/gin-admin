package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// Menu 菜单验证
type Menu struct {
	base.Request
	PageListValidate
	ID         int64        `json:"id" form:"id" label:"ID"`
	RoleId     string       `json:"roleId" form:"roleId" label:"角色id"`
	MenuId     int64        `json:"menuId" form:"menuId" label:"菜单id"`
	Pid        int64        `json:"pid" form:"pid" label:"父级id"`
	Name       string       `json:"name" form:"name" label:"路由名称|功能标识"`
	Type       int64        `json:"type" form:"type" label:"类型 1=菜单 2=功能"`
	Status     int64        `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort       int64        `json:"sort" form:"sort" label:"排序"`
	Meta       Meta         `json:"meta" form:"meta" label:"菜单元数据"`
	MenuAction ActionCreate `json:"menuAction" form:"menuAction" label:"菜单功能"`
	RoleMenus  []RoleMenu   `json:"roleMenus" form:"roleMenus" label:"角色菜单"`
}

// MenuCreate 菜单创建
type MenuCreate struct {
	Pid        int64        `json:"pid" form:"pid" label:"父级id"`
	Type       int64        `json:"type" form:"type" label:"类型 1=菜单 2=功能"`
	Name       string       `json:"name" form:"name" label:"路由名称|功能标识"`
	IsLink     int64        `json:"isLink" form:"isLink" label:"是否外链 1=是 2=否 默认=2"`
	Status     int64        `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort       int64        `json:"sort" form:"sort" label:"排序"`
	Meta       Meta         `json:"meta" form:"meta" label:"菜单元数据"`
	MenuAction ActionCreate `json:"menuAction" form:"menuAction" label:"菜单功能"`
	RoleMenus  []RoleMenu   `json:"roleMenus" form:"roleMenus" label:"角色菜单"`
}

// MenuUpdate 菜单更新
type MenuUpdate struct {
	Pid        int64        `json:"pid" form:"pid" label:"父级id"`
	Type       int64        `json:"type" form:"type" label:"类型 1=菜单 2=功能"`
	Name       string       `json:"name" form:"name" label:"路由名称|功能标识"`
	Status     int64        `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	Sort       int64        `json:"sort" form:"sort" label:"排序"`
	Meta       Meta         `json:"meta" form:"meta" label:"菜单元数据"`
	MenuAction ActionCreate `json:"menuAction" form:"menuAction" label:"菜单功能"`
	RoleMenus  []RoleMenu   `json:"roleMenus" form:"roleMenus" label:"角色菜单"`
}

// Meta 菜单元数据
type Meta struct {
	MenuId      int64  `json:"menuId" form:"menuId" label:"菜单id"`
	Title       string `json:"title" form:"title" label:"菜单名称"`
	Icon        string `json:"icon" form:"icon" label:"菜单图标"`
	TransKey    string `json:"transKey" form:"transKey" label:"翻译键"`
	Path        string `json:"path" form:"path" label:"路由路径"`
	Redirect    string `json:"redirect" form:"redirect" label:"重定向"`
	Component   string `json:"component" form:"component" label:"组件路径"`
	IsHide      int64  `json:"isHide" form:"isHide" label:"是否隐藏 1=是 2=否"`
	IsKeepAlive int64  `json:"isKeepAlive" form:"isKeepAlive" label:"是否缓存 1=是 2=否"`
	IsAffix     int64  `json:"isAffix" form:"isAffix" label:"是否固定 1=是 2=否"`
	IsLink      string `json:"isLink" form:"isLink" label:"外链地址"`
	IsIframe    int64  `json:"isIframe" form:"isIframe" label:"是否内嵌 1=是 2=否"`
}

// RoleMenu 角色菜单
type RoleMenu struct {
	RoleId int64  `json:"roleId" form:"roleId" label:"角色id"`
	MenuId int64  `json:"menuId" form:"menuId" label:"菜单id"`
	Name   string `json:"name" form:"name" label:"角色名称"`
}

// Validate 菜单元数据验证
func (s Meta) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Str("title", s.Title).Required().Errors(),
		fields.Str("icon", s.Icon).Required().Errors(),
		fields.Str("path", s.Path).Required().Errors(),
		fields.Str("component", s.Component).Required().Errors(),
		fields.Int("isHide", s.IsHide).Required().Errors(),
		fields.Int("isKeepAlive", s.IsKeepAlive).Required().Errors(),
		fields.Int("isAffix", s.IsAffix).Required().Errors(),
		fields.Int("isIframe", s.IsIframe).Required().Errors(),
	)
}

// Validate 角色菜单验证
func (s RoleMenu) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Str("name", s.Name).Required().Errors()
}

// ValidateWithScene 请求验证
func (s Menu) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "RoleMenu":
		return fields.Str("roleId", s.RoleId).Required().Errors()
	case "Create":
		return validatex.Check(
			fields.Int("type", s.Type).Required().Errors(),
			fields.Str("name", s.Name).Required().Errors(),
			fields.Int("status", s.Status).Required().Errors(),
			fields.Struct("meta", s.Meta),
		)
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Int("type", s.Type).Required().Errors(),
			fields.Str("name", s.Name).Required().Errors(),
			fields.Int("status", s.Status).Required().Errors(),
			fields.Struct("meta", s.Meta),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Menu) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s Menu) Translates() validatex.Translates {
	return validatex.Translates{
		"page":             "页码",
		"pageSize":         "每页数量",
		"id":               "ID",
		"roleId":           "角色ID",
		"menuId":           "菜单ID",
		"type":             "类型",
		"name":             "名称",
		"status":           "状态",
		"meta.title":       "菜单名称",
		"meta.icon":        "菜单图标",
		"meta.path":        "路由路径",
		"meta.component":   "组件路径",
		"meta.isHide":      "是否隐藏",
		"meta.isKeepAlive": "是否缓存",
		"meta.isAffix":     "是否固定",
		"meta.isIframe":    "是否内嵌",
	}
}
