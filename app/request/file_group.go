package request

import (
	"gin/app/errcode"
	"gin/common/base"

	"github.com/gookit/validate"
)

// FileGroup 文件分组请求验证
type FileGroup struct {
	base.Request
	ID   int64  `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Pid  int64  `json:"pid" form:"pid" validate:"int" label:"父级ID"`
	Name string `json:"name" form:"name" validate:"required|maxLen:20" label:"分组名称"`
	Sort int64  `json:"sort" form:"sort" validate:"int" label:"排序"`
	PageListValidate
}

// Validate 请求验证
func (s FileGroup) Validate(data FileGroup, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
func (s FileGroup) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{"Pid", "Name", "Sort"},
		"Update": []string{"ID", "Pid", "Name", "Sort"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s FileGroup) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"gt":                           "字段 {field} 必须大于 0",
		"maxLen":                       "{field} 长度不能超过 {max} 个字符",
		"PageListValidate.Page.gt":     "页码必须大于 0",
		"PageListValidate.PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s FileGroup) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Pid":                       "父级ID",
		"Name":                      "分组名称",
		"Sort":                      "排序",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
