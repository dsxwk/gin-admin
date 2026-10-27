package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// OperatorLog 操作日志请求验证
type OperatorLog struct {
	base.Request
	ID  int64   `json:"id" form:"id" label:"ID"`
	IDs []int64 `json:"ids" form:"ids" label:"ID列表"`
	Ip  string  `json:"ip" form:"ip" label:"ip地址"`
	PageListValidate
}

// ValidateWithScene 请求验证
func (s OperatorLog) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Detail", "Delete":
		return fields.Int("ID", s.ID).Required().Gt(0).Errors()
	case "BatchDelete":
		return fields.Slice("IDs", s.IDs).Required().
			Each("gt", "", func(value int64) bool {
				return value > 0
			}, int64(0)).
			Errors()
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s OperatorLog) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s OperatorLog) Translates() validatex.Translates {
	return validatex.Translates{
		"ID":                        "ID",
		"IDs":                       "ID列表",
		"Ip":                        "ip地址",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
