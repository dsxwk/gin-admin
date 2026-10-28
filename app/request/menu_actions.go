package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// MenuActions Validator
type MenuActions struct {
	base.Request
	ID        int64  `json:"id" form:"id" label:"ID"`
	MenuId    int64  `json:"menuId" form:"menuId" label:"菜单id"`
	Type      int64  `json:"type" form:"type" label:"类型 1=header 2=operation"`
	BtnType   string `json:"btnType" form:"btnType" label:"按钮类型 text|btn"`
	BtnStyle  string `json:"btnStyle" form:"btnStyle" label:"按钮样式"`
	BtnSize   string `json:"btnSize" form:"btnSize" label:"按钮尺寸"`
	IsConfirm int64  `json:"isConfirm" form:"isConfirm" label:"是否确认 1=是 2=否"`
	Label     string `json:"label" form:"label" label:"功能名称"`
	AuthValue string `json:"authValue" form:"authValue" label:"权限标识"`
	IsLink    int64  `json:"isLink" form:"isLink" label:"是否为链接 1=是 2=否"`
	Sort      int64  `json:"sort" form:"sort" label:"排序"`
	PageListValidate
}

// RoleAction 角色功能
type RoleAction struct {
	RoleId   int64  `json:"roleId" form:"roleId" label:"角色id"`
	ActionId int64  `json:"actionId" form:"actionId" label:"功能id"`
	Name     string `json:"name" form:"name" label:"角色名称"`
}

// ActionList 功能列表
type ActionList struct {
	MenuId int64 `json:"menuId" form:"menuId" label:"菜单id"`
}

// ActionCreate 功能创建
type ActionCreate struct {
	MenuId    int64  `json:"menuId" form:"menuId" label:"菜单id"`
	Type      int64  `json:"type" form:"type" label:"类型 1=header 2=operation"`
	BtnType   string `json:"btnType" form:"btnType" label:"按钮类型 text|btn"`
	BtnStyle  string `json:"btnStyle" form:"btnStyle" label:"按钮样式"`
	BtnSize   string `json:"btnSize" form:"btnSize" label:"按钮尺寸"`
	IsConfirm int64  `json:"isConfirm" form:"isConfirm" label:"是否确认 1=是 2=否"`
	Label     string `json:"label" form:"label" label:"功能名称"`
	TransKey  string `json:"transKey" form:"transKey" label:"翻译键"`
	AuthValue string `json:"authValue" form:"authValue" label:"权限标识"`
	IsLink    int64  `json:"isLink" form:"isLink" label:"是否为链接 1=是 2=否"`
}

// ActionUpdate 功能更新
type ActionUpdate struct {
	MenuId    int64  `json:"menuId" form:"menuId" label:"菜单id"`
	Type      int64  `json:"type" form:"type" label:"类型 1=header 2=operation"`
	BtnType   string `json:"btnType" form:"btnType" label:"按钮类型 text|btn"`
	BtnStyle  string `json:"btnStyle" form:"btnStyle" label:"按钮样式"`
	BtnSize   string `json:"btnSize" form:"btnSize" label:"按钮尺寸"`
	IsConfirm int64  `json:"isConfirm" form:"isConfirm" label:"是否确认 1=是 2=否"`
	Label     string `json:"label" form:"label" label:"功能名称"`
	TransKey  string `json:"transKey" form:"transKey" label:"翻译键"`
	AuthValue string `json:"authValue" form:"authValue" label:"权限标识"`
	IsLink    int64  `json:"isLink" form:"isLink" label:"是否为链接 1=是 2=否"`
}

// Validate 功能创建验证
func (s ActionCreate) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Int("type", s.Type).Required().Errors(),
		fields.Str("btnType", s.BtnType).Required().Errors(),
		fields.Str("btnStyle", s.BtnStyle).Required().Errors(),
		fields.Str("btnSize", s.BtnSize).Required().Errors(),
		fields.Int("isConfirm", s.IsConfirm).Required().Errors(),
		fields.Str("label", s.Label).Required().Errors(),
		fields.Str("authValue", s.AuthValue).Required().Errors(),
		fields.Int("isLink", s.IsLink).Required().Errors(),
	)
}

// Validate 功能更新验证
func (s ActionUpdate) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Int("type", s.Type).Required().Errors(),
		fields.Str("btnType", s.BtnType).Required().Errors(),
		fields.Str("btnStyle", s.BtnStyle).Required().Errors(),
		fields.Str("btnSize", s.BtnSize).Required().Errors(),
		fields.Int("isConfirm", s.IsConfirm).Required().Errors(),
		fields.Str("label", s.Label).Required().Errors(),
		fields.Str("authValue", s.AuthValue).Required().Errors(),
		fields.Int("isLink", s.IsLink).Required().Errors(),
	)
}

// ValidateWithScene 请求验证
func (s MenuActions) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create", "Update":
		return validatex.Check(
			fields.Int("type", s.Type).Required().Errors(),
			fields.Str("btnType", s.BtnType).Required().Errors(),
			fields.Str("btnStyle", s.BtnStyle).Required().Errors(),
			fields.Str("btnSize", s.BtnSize).Required().Errors(),
			fields.Int("isConfirm", s.IsConfirm).Required().Errors(),
			fields.Str("label", s.Label).Required().Errors(),
			fields.Str("authValue", s.AuthValue).Required().Errors(),
			fields.Int("isLink", s.IsLink).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s MenuActions) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s MenuActions) Translates() validatex.Translates {
	return validatex.Translates{
		"page":      "页码",
		"pageSize":  "每页数量",
		"id":        "ID",
		"menuId":    "菜单ID",
		"type":      "类型",
		"btnType":   "按钮类型",
		"btnStyle":  "按钮样式",
		"btnSize":   "按钮尺寸",
		"isConfirm": "是否确认",
		"label":     "功能名称",
		"authValue": "权限标识",
		"isLink":    "是否为链接",
	}
}
