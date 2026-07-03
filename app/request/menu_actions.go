package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// MenuActions Validator
type MenuActions struct {
	base.BaseRequest
	ID          int64        `json:"id" validate:"required|int|gt:0" label:"ID"`
	MenuId      int64        `json:"menuId" form:"menuId" validate:"int" label:"菜单id"`
	Pid         int64        `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Type        int64        `json:"type" form:"type" validate:"required|int" label:"类型 1=header 2=operation"`
	BtnType     string       `json:"btnType" form:"btnType" validate:"required|string" label:"按钮类型 text|btn"`
	BtnStyle    string       `json:"btnStyle" form:"btnStyle" validate:"required|string" label:"按钮样式"`
	BtnSize     string       `json:"btnSize" form:"btnSize" validate:"required|string" label:"按钮尺寸"`
	IsConfirm   int64        `json:"isConfirm" form:"isConfirm" validate:"required|int" label:"是否确认 1=是 2=否"`
	Label       string       `json:"label" form:"label" validate:"required|string" label:"功能名称"`
	AuthValue   string       `json:"authValue" form:"authValue" validate:"required|string" label:"权限标识"`
	IsLink      int64        `json:"isLink" form:"isLink" validate:"required|int" label:"是否为链接 1=是 2=否"`
	Sort        int64        `json:"sort" form:"sort" validate:"required|int" label:"排序"`
	RoleActions []RoleAction `json:"roleActions" form:"roleActions" validate:"" label:"角色功能"`
	PageListValidate
}

// RoleAction 角色功能
type RoleAction struct {
	RoleId   int64  `json:"roleId" form:"roleId" validate:"int" label:"角色id"`
	Name     string `json:"name" form:"name" validate:"int" label:"角色名称"`
	ActionId int64  `json:"actionId" form:"actionId" validate:"int" label:"功能id"`
}

// ActionList 功能列表
type ActionList struct {
	MenuId int64 `json:"menuId" form:"menuId" validate:"int" label:"菜单id"`
}

// ActionCreate 功能创建
type ActionCreate struct {
	Pid       int64  `json:"pid" form:"pid" validate:"int" label:"父级id"`
	MenuId    int64  `json:"menuId" form:"menuId" validate:"int" label:"菜单id"`
	Type      int64  `json:"type" form:"type" validate:"required|int" label:"类型 1=header 2=operation"`
	BtnType   string `json:"btnType" form:"btnType" validate:"required|string" label:"按钮类型 text|btn"`
	BtnStyle  string `json:"btnStyle" form:"btnStyle" validate:"required|string" label:"按钮样式"`
	BtnSize   string `json:"btnSize" form:"btnSize" validate:"required|string" label:"按钮尺寸"`
	IsConfirm int64  `json:"isConfirm" form:"isConfirm" validate:"required|int" label:"是否确认 1=是 2=否"`
	Label     string `json:"label" form:"label" validate:"required|string" label:"功能名称"`
	AuthValue string `json:"authValue" form:"authValue" validate:"required|string" label:"权限标识"`
	IsLink    int64  `json:"isLink" form:"isLink" validate:"required|int" label:"是否为链接 1=是 2=否"`
	Sort      int64  `json:"sort" form:"sort" validate:"required|int" label:"排序"`
}

// ActionUpdate 功能更新
type ActionUpdate struct {
	ID        int64  `json:"id" form:"id" validate:"required|int" label:"功能id"`
	MenuId    int64  `json:"menuId" form:"menuId" validate:"int" label:"菜单id"`
	Pid       int64  `json:"pid" form:"pid" validate:"int" label:"父级id"`
	Type      int64  `json:"type" form:"type" validate:"required|int" label:"类型 1=header 2=operation"`
	BtnType   string `json:"btnType" form:"btnType" validate:"required|string" label:"按钮类型 text|btn"`
	BtnStyle  string `json:"btnStyle" form:"btnStyle" validate:"required|string" label:"按钮样式"`
	BtnSize   string `json:"btnSize" form:"btnSize" validate:"required|string" label:"按钮尺寸"`
	IsConfirm int64  `json:"isConfirm" form:"isConfirm" validate:"required|int" label:"是否确认 1=是 2=否"`
	Label     string `json:"label" form:"label" validate:"required|string" label:"功能名称"`
	AuthValue string `json:"authValue" form:"authValue" validate:"required|string" label:"权限标识"`
	IsLink    int64  `json:"isLink" form:"isLink" validate:"required|int" label:"是否为链接 1=是 2=否"`
	Sort      int64  `json:"sort" form:"sort" validate:"required|int" label:"排序"`
}

// Validate 请求验证
func (s MenuActions) Validate(data MenuActions, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}

	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s MenuActions) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"List": []string{
			"ID",
		},
		"Create": []string{
			"Pid",
			"MenuId",
			"Type",
			"BtnType",
			"BtnStyle",
			"BtnSize",
			"IsConfirm",
			"Label",
			"AuthValue",
			"IsLink",
			"Sort",
		},
		"Update": []string{
			"ID",
			"Pid",
			"MenuId",
			"Type",
			"BtnType",
			"BtnStyle",
			"BtnSize",
			"IsConfirm",
			"Label",
			"AuthValue",
			"IsLink",
			"Sort",
		},
		"Detail": []string{
			"ID",
		},
		"Delete": []string{
			"ID",
		},
	})
}

// Messages 验证器错误消息
func (s MenuActions) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"PageListValidate.Page.gt":     "字段 {field} 需大于 0",
		"PageListValidate.PageSize.gt": "字段 {field} 需大于 0",
	}
}

// Translates 字段翻译
func (s MenuActions) Translates() map[string]string {
	return validate.MS{
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
		"ID":                        "ID",
		"Pid":                       "父级ID",
		"MenuId":                    "菜单ID",
		"Type":                      "类型",
		"BtnType":                   "按钮类型",
		"BtnStyle":                  "按钮样式",
		"BtnSize":                   "按钮尺寸",
		"IsConfirm":                 "是否确认",
		"Label":                     "功能名称",
		"AuthValue":                 "权限标识",
		"IsLink":                    "是否为链接",
		"Sort":                      "排序",
	}
}
