package request

import (
	"errors"
	"github.com/gookit/validate"
)

type Meta struct {
	Title       string  `json:"title" validate:"required" label:"菜单名称"`
	Icon        string  `json:"icon" validate:"required" label:"菜单图标"`
	IsHide      bool    `json:"isHide" validate:"required" label:"是否隐藏"`
	IsKeepAlive bool    `json:"isKeepAlive" validate:"required" label:"是否缓存"`
	IsAffix     bool    `json:"isAffix" validate:"required" label:"是否固定"`
	IsLink      string  `json:"isLink" validate:"required" label:"外链/内嵌时链接地址"` // 外链/内嵌时链接地址（http:xxx.com），开启外链条件，`1、isLink: 链接地址不为空`
	IsIframe    bool    `json:"isIframe" validate:"required" label:"是否内嵌"`     // 是否内嵌，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
	Roles       []int64 `json:"roles" validate:"required" label:"菜单角色"`        // 权限标识，取角色管理
}

// MenuRequest 菜单请求验证
type MenuRequest struct {
	ID       int64  `json:"id" validate:"required" label:"ID"`
	PID      int64  `json:"pid" validate:"required" label:"父级ID"`
	Name     string `json:"name" validate:"required" label:"路由名称"`
	Path     string `json:"path" validate:"required" label:"路由路径"`
	Redirect string `json:"redirect" validate:"required" label:"重定向"`
	IsLink   bool   `json:"isLink" validate:"required" label:"是否外链"`
	Sort     int64  `json:"sort" validate:"required" label:"排序"`
	Meta     Meta   `json:"meta" validate:"required" label:"元数据"`
	Context
}

type MenuSearch struct {
	RoleIds string `form:"roleIds"`
}

// Validate 请求验证
func (s MenuRequest) Validate(data MenuRequest, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}

	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s MenuRequest) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"list":   []string{},
		"create": []string{"Name", "Path", "Meta.Title", "Meta.Roles"},
		"update": []string{"ID", "Name", "Path", "Meta.Title", "Meta.Roles"},
		"delete": []string{"ID"},
	})
}

// Messages 验证器错误消息
func (s MenuRequest) Messages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
	}
}

// Translates 字段翻译
func (s MenuRequest) Translates() map[string]string {
	return validate.MS{
		"ID":         "ID",
		"PID":        "父级ID",
		"Meta.Title": "菜单名称",
		"Name":       "路由名称",
		"Meta.Roles": "菜单角色",
	}
}
