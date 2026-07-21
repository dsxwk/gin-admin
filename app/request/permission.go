package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// Permission 请求验证
type Permission struct {
	base.BaseRequest
	ID     int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Key    string `json:"key" form:"key" validate:"required" label:"权限标识"`
	Method string `json:"method" form:"method" validate:"required" label:"请求方式"`
	Uri    string `json:"uri" form:"uri" validate:"required" label:"路由地址"`
	PageListValidate
}

// Validate 请求验证
func (s Permission) Validate(data Permission, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s Permission) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List": []string{"PageListValidate.Page", "PageListValidate.PageSize"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s Permission) Messages() map[string]string {
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
func (s Permission) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Key":                       "权限标识",
		"Method":                    "请求方式",
		"Uri":                       "路由地址",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
