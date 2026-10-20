package request

import (
	"gin/app/errcode"
	"gin/common/base"

	"github.com/gookit/validate"
)

// File 文件请求验证
type File struct {
	base.Request
	ID      int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	GroupId *int64 `json:"groupId" form:"groupId" validate:"" label:"分组ID"`
	PageListValidate
}

// FileUpdateGroup 修改文件分组请求验证
type FileUpdateGroup struct {
	base.Request
	ID      int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	GroupId *int64 `json:"groupId" form:"groupId" validate:"" label:"分组ID"`
}

// Validate 请求验证
func (s File) Validate(data File, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
func (s File) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s File) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"gt":                           "字段 {field} 必须大于 0",
		"PageListValidate.Page.gt":     "页码必须大于 0",
		"PageListValidate.PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s File) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"GroupId":                   "分组ID",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}

// Validate 请求验证
func (s FileUpdateGroup) Validate(data FileUpdateGroup, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
func (s FileUpdateGroup) ConfigValidation(v *validate.Validation) {
	v.WithScenes(validate.SValues{
		"UpdateGroup": []string{"ID"},
	})
}

// Messages 验证器错误消息
func (s FileUpdateGroup) Messages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
		"gt":       "字段 {field} 必须大于 0",
	}
}

// Translates 字段翻译
func (s FileUpdateGroup) Translates() map[string]string {
	return validate.MS{
		"ID":      "ID",
		"GroupId": "分组ID",
	}
}
