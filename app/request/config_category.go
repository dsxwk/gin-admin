package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// ConfigCategory 请求验证
type ConfigCategory struct {
	base.BaseRequest
	ID   int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Name string `json:"name" form:"name" validate:"required" label:"分类名称"`
	PageListValidate
}

// Validate 请求验证
func (s ConfigCategory) Validate(data ConfigCategory, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigCategoryCreate ConfigCategory创建验证
type ConfigCategoryCreate struct {
	Name string `json:"name" form:"name" validate:"required" label:"分类名称"`
}

// ConfigCategoryUpdate ConfigCategory更新验证
type ConfigCategoryUpdate struct {
	Name string `json:"name" form:"name" validate:"required" label:"分类名称"`
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s ConfigCategory) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List": []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{
			"Name",
		},
		"Update": []string{
			"ID",
			"Name",
		},
		"Detail": []string{"ID"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s ConfigCategory) Messages() map[string]string {
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
func (s ConfigCategory) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Name":                      "分类名称",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
