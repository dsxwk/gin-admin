package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// Roles 请求验证
type Roles struct {
	base.Request
	ID              int64            `uri:"id" form:"id" label:"ID"`
	Name            string           `json:"name" form:"name" label:"角色名称"`
	Desc            string           `json:"desc" form:"desc" label:"角色描述"`
	Status          int64            `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	UserRoles       []UserRoles      `json:"UserRoles" form:"UserRoles" label:"用户角色"`
	RoleMenus       []RoleMenu       `json:"roleMenus" form:"roleMenus" label:"角色菜单"`
	RolePermissions []RolePermission `json:"rolePermissions" form:"rolePermissions" label:"角色权限"`
	PageListValidate
}

// RoleCreate 角色创建验证
type RoleCreate struct {
	Name            string           `json:"name" form:"name" label:"角色名称"`
	Desc            string           `json:"desc" form:"desc" label:"角色描述"`
	Status          int64            `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	UserRoles       []UserRoles      `json:"UserRoles" form:"UserRoles" label:"用户角色"`
	RoleMenus       []RoleMenu       `json:"roleMenus" form:"roleMenus" label:"角色菜单"`
	RolePermissions []RolePermission `json:"rolePermissions" form:"rolePermissions" label:"角色权限"`
}

// UserRoles 用户角色
type UserRoles struct {
	RoleId int64  `json:"roleId" form:"roleId" label:"角色ID"`
	UserId int64  `json:"userId" form:"userId" label:"用户ID"`
	Name   string `json:"name" form:"name" label:"角色名称"`
}

// RolePermission 角色权限
type RolePermission struct {
	RoleId       int64 `json:"roleId" form:"roleId" label:"角色id"`
	PermissionId int64 `json:"permissionId" form:"permissionId" label:"权限id"`
}

// RoleUpdate 角色更新验证
type RoleUpdate struct {
	Name            string           `json:"name" form:"name" label:"角色名称"`
	Desc            string           `json:"desc" form:"desc" label:"角色描述"`
	Status          int64            `json:"status" form:"status" label:"状态 1=启用 2=停用"`
	UserRoles       []UserRoles      `json:"userRoles" form:"userRoles" label:"用户角色"`
	RoleMenus       []RoleMenu       `json:"roleMenus" form:"roleMenus" label:"角色菜单"`
	RolePermissions []RolePermission `json:"rolePermissions" form:"rolePermissions" label:"角色权限"`
}

// Validate 用户角色验证
func (s UserRoles) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Int("roleId", s.RoleId).Required().Gt(0).Errors(),
		fields.Int("userId", s.UserId).Required().Gt(0).Errors(),
		fields.Str("name", s.Name).Required().Max(20).Errors(),
	)
}

// Validate 角色权限验证
func (s RolePermission) Validate(fields validatex.Fields) validatex.Errors {
	return validatex.Check(
		fields.Int("roleId", s.RoleId).Required().Errors(),
		fields.Int("permissionId", s.PermissionId).Required().Errors(),
	)
}

// ValidateWithScene 请求验证
func (s Roles) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return validatex.Check(
			fields.Str("name", s.Name).Required().Min(1).Max(20).Errors(),
			fields.Str("desc", s.Desc).Max(100).Errors(),
			fields.Int("status", s.Status).Required().Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Str("name", s.Name).Required().Min(1).Max(20).Errors(),
			fields.Str("desc", s.Desc).Max(100).Errors(),
			fields.Int("status", s.Status).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s Roles) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s Roles) Translates() validatex.Translates {
	return validatex.Translates{
		"id":           "ID",
		"name":         "角色名称",
		"desc":         "角色描述",
		"status":       "状态 1=启用 2=停用",
		"roleId":       "角色ID",
		"userId":       "用户ID",
		"permissionId": "权限ID",
		"page":         "页码",
		"pageSize":     "每页数量",
	}
}
