package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// Menu Validator
type Menu struct {
	base.BaseRequest
	PageListValidate
	Id        int64      `json:"id" validate:"required|int|gt:0" label:"ID"`
	RoleId    string     `json:"roleId" form:"roleId" validate:"required" label:"角色id"`
	MenuId    int64      `json:"menuId" form:"menuId" validate:"required|int|gt:0" label:"菜单id"`
	Pid       int64      `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name      string     `json:"name" form:"name" validate:"required" label:"路由名称"`
	Path      string     `json:"path" form:"path" validate:"required" label:"路由路径"`
	Redirect  string     `json:"redirect" form:"redirect" validate:"" label:"重定向"`
	Component string     `json:"component" form:"component" validate:"required" label:"组件路径"`
	IsLink    int64      `json:"isLink" form:"isLink" validate:"required|int" label:"是否外链 1=是 2=否 默认=2"`
	Status    int64      `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort      int64      `json:"sort" form:"sort" validate:"int" label:"排序"`
	Meta      Meta       `json:"meta" form:"meta" validate:"required" label:"菜单元数据"`
	RoleMenus []RoleMenu `json:"roleMenus" form:"roleMenus" validate:"" label:"角色菜单"`
}

// MenuCreate 菜单创建
type MenuCreate struct {
	Pid       int64      `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name      string     `json:"name" form:"name" validate:"required" label:"路由名称"`
	Path      string     `json:"path" form:"path" validate:"required" label:"路由路径"`
	Redirect  string     `json:"redirect" form:"redirect" validate:"" label:"重定向"`
	Component string     `json:"component" form:"component" validate:"required" label:"组件路径"`
	IsLink    int64      `json:"isLink" form:"isLink" validate:"required|int" label:"是否外链 1=是 2=否 默认=2"`
	Status    int64      `json:"status" form:"status" validate:"required|int" label:"状态 1=启用 2=停用"`
	Sort      int64      `json:"sort" form:"sort" validate:"int" label:"排序"`
	Meta      Meta       `json:"meta" form:"meta" validate:"required" label:"菜单元数据"`
	RoleMenus []RoleMenu `json:"roleMenus" form:"roleMenus" validate:"" label:"角色菜单"`
}

// MenuUpdate 菜单更新
type MenuUpdate struct {
	Id        int64      `json:"id" validate:"required|int|gt:0" label:"ID"`
	Pid       int64      `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Name      string     `json:"name" form:"name" validate:"required" label:"路由名称"`
	Path      string     `json:"path" form:"path" validate:"required" label:"路由路径"`
	Redirect  string     `json:"redirect" form:"redirect" validate:"" label:"重定向"`
	Component string     `json:"component" form:"component" validate:"required" label:"组件路径"`
	IsLink    int64      `json:"isLink" form:"isLink" validate:"required|int" label:"是否外链 1=是 2=否 默认=2"`
	Status    int64      `json:"status" form:"status" validate:"int" label:"状态 1=启用 2=停用"`
	Sort      int64      `json:"sort" form:"sort" validate:"int" label:"排序"`
	Meta      Meta       `json:"meta" form:"meta" validate:"required" label:"菜单元数据"`
	RoleMenus []RoleMenu `json:"roleMenus" form:"roleMenus" validate:"" label:"角色菜单"`
}

// Meta 菜单元数据
type Meta struct {
	MenuId      int64  `json:"menuId" form:"menuId" validate:"int" label:"菜单id"`
	Title       string `json:"title" form:"title" validate:"required" label:"菜单名称"`
	Icon        string `json:"icon" form:"icon" validate:"required" label:"菜单图标"`
	IsHide      int64  `json:"isHide" form:"isHide" validate:"required|int" label:"是否隐藏 1=是 2=否"`
	IsKeepAlive int64  `json:"isKeepAlive" form:"isKeepAlive" validate:"required|int" label:"是否缓存 1=是 2=否"`
	IsAffix     int64  `json:"isAffix" form:"isAffix" validate:"required|int" label:"是否固定 1=是 2=否"`
	IsLink      string `json:"isLink" form:"isLink" validate:"string" label:"外链/内嵌时链接地址(http:xxx.com),开启外链条件1 isLink:链接地址不为空"`
	IsIframe    int64  `json:"isIframe" form:"isIframe" validate:"required|int" label:"是否内嵌 1=是 2=否 开启条件1 isIframe:true 2 isLink:链接地址不为空"`
}

// RoleMenu 角色菜单
type RoleMenu struct {
	RoleId int64 `json:"roleId" form:"roleId" validate:"int" label:"角色id"`
	MenuId int64 `json:"menuId" form:"menuId" validate:"int" label:"菜单id"`
}

// Validate 请求验证
func (s Menu) Validate(data Menu, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}

	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Menu) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"List": []string{
			"PageListValidate.Page",
			"PageListValidate.PageSize",
		},
		"RoleMenu": []string{
			"RoleId",
		},
		"Create": []string{
			"Pid",
			"Name",
			"Path",
			"Redirect",
			"IsLink",
			"Status",
			"Sort",
			"Meta",
			"Meta.MenuId",
			"Meta.Title",
			"Meta.Icon",
			"Meta.IsHide",
			"Meta.IsAffix",
			"Meta.IsLink",
			"Meta.IsIframe",
		},
		"Update": []string{
			"Id",
			"Pid",
			"Name",
			"Path",
			"Redirect",
			"IsLink",
			"Status",
			"Sort",
			"Meta",
			"Meta.MenuId",
			"Meta.Title",
			"Meta.Icon",
			"Meta.IsHide",
			"Meta.IsAffix",
			"Meta.IsLink",
			"Meta.IsIframe",
		},
		"Detail": []string{
			"Id",
		},
		"Delete": []string{
			"Id",
		},
	})
}

// Messages 验证器错误消息
func (s Menu) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"PageListValidate.Page.gt":     "字段 {field} 需大于 0",
		"PageListValidate.PageSize.gt": "字段 {field} 需大于 0",
	}
}

// Translates 字段翻译
func (s Menu) Translates() map[string]string {
	return validate.MS{
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
		"ID":                        "ID",
	}
}
