package request

import (
	"errors"
	"gin/common/base"
	"github.com/gookit/validate"
)

// ImportRecords 请求验证
type ImportRecords struct {
	base.BaseRequest
	ID   int64       `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	Type int64       `json:"type" form:"type" validate:"required|int" label:"导入类型"`
	Name string      `json:"name" form:"name" validate:"required" label:"类型名称"`
	Data interface{} `json:"data" form:"data" validate:"required" label:"导入数据"`
	PageListValidate
}

// Validate 请求验证
func (s ImportRecords) Validate(data ImportRecords, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s ImportRecords) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List":   []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s ImportRecords) Messages() map[string]string {
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
func (s ImportRecords) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"Type":                      "导入类型",
		"Name":                      "类型名称",
		"Data":                      "导入数据",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
