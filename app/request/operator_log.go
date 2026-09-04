package request

import (
	"gin/common/base"
	"gin/common/errcode"

	"github.com/gookit/validate"
)

// OperatorLog 操作日志请求验证
type OperatorLog struct {
	base.BaseRequest
	ID  int64   `json:"id" form:"id" validate:"required|int|gt:0" label:"ID"`
	IDs []int64 `json:"ids" validate:"required" label:"ID列表"`
	Ip  string  `json:"ip" form:"ip" validate:"" label:"ip地址"`
	PageListValidate
}

// Validate 请求验证
func (s OperatorLog) Validate(data OperatorLog, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errcode.ArgsError().WithMsg(v.Errors.One())
	}
	return nil
}

// ConfigValidation 配置验证场景
func (s OperatorLog) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List":        []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Detail":      []string{"ID"},
		"Delete":      []string{"ID"},
		"BatchDelete": []string{"IDs"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s OperatorLog) Messages() map[string]string {
	return validate.MS{
		"required":                     "字段 {field} 必填",
		"int":                          "字段 {field} 必须为整数",
		"gt":                           "字段 {field} 必须大于 0",
		"PageListValidate.Page.gt":     "页码必须大于 0",
		"PageListValidate.PageSize.gt": "每页数量必须大于 0",
	}
}

// Translates 字段翻译
func (s OperatorLog) Translates() map[string]string {
	return validate.MS{
		"ID":                        "ID",
		"IDs":                       "ID列表",
		"Ip":                        "ip地址",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
