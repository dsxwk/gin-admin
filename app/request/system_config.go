package request

import (
	"gin/common/base"

	"github.com/dsxwk/validatex"
)

// SystemConfig 请求验证
type SystemConfig struct {
	base.Request
	ID               int64                     `json:"id" form:"id" label:"id"`
	Key              string                    `json:"key" form:"key" label:"标识"`
	Name             string                    `json:"name" form:"name" label:"名称"`
	DefaultValue     string                    `json:"defaultValue" form:"defaultValue" label:"默认值"`
	OptionValue      string                    `json:"optionValue" form:"optionValue" label:"可选值"`
	Type             int64                     `json:"type" form:"type" label:"配置类型"`
	ConfigCategoryId int64                     `json:"configCategoryId" form:"configCategoryId" label:"配置分类Id"`
	List             []SystemConfigValueUpdate `json:"list" form:"list" label:"配置列表"`
	PageListValidate
}

// SystemConfigCreate 系统配置创建
type SystemConfigCreate struct {
	Key              string `json:"key" form:"key" label:"标识"`
	Name             string `json:"name" form:"name" label:"名称"`
	DefaultValue     string `json:"defaultValue" form:"defaultValue" label:"默认值"`
	OptionValue      string `json:"optionValue" form:"optionValue" label:"可选值"`
	Type             int64  `json:"type" form:"type" label:"配置类型"`
	ConfigCategoryId int64  `json:"configCategoryId" form:"configCategoryId" label:"配置分类Id"`
}

// SystemConfigUpdate 系统配置更新
type SystemConfigUpdate struct {
	Key              string `json:"key" form:"key" label:"标识"`
	Name             string `json:"name" form:"name" label:"名称"`
	DefaultValue     string `json:"defaultValue" form:"defaultValue" label:"默认值"`
	OptionValue      string `json:"optionValue" form:"optionValue" label:"可选值"`
	Type             int64  `json:"type" form:"type" label:"配置类型"`
	ConfigCategoryId int64  `json:"configCategoryId" form:"configCategoryId" label:"配置分类Id"`
}

// SystemConfigValueUpdate 系统配置批量更新
type SystemConfigValueUpdate struct {
	ID           int64  `json:"id" form:"id" label:"id"`
	Key          string `json:"key" form:"key" label:"标识"`
	DefaultValue string `json:"defaultValue" form:"defaultValue" label:"默认值"`
}

// SystemConfigUpdates 系统配置批量更新验证
type SystemConfigUpdates struct {
	base.Request
	List []SystemConfigValueUpdate `json:"list" form:"list" label:"配置列表"`
}

// Validate SystemConfigValueUpdate验证
func (s SystemConfigValueUpdate) Validate(fields validatex.Fields) validatex.Errors {
	return fields.Int("id", s.ID).Required().Gt(0).Errors()
}

// ValidateWithScene 请求验证
func (s SystemConfig) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "List":
		return s.PageListValidate.Validate(fields)
	case "Create":
		return validatex.Check(
			fields.Str("key", s.Key).Required().Errors(),
			fields.Str("name", s.Name).Required().Errors(),
			fields.Str("defaultValue", s.DefaultValue).Required().Errors(),
			fields.Str("optionValue", s.OptionValue).Required().Errors(),
			fields.Int("type", s.Type).Required().Errors(),
			fields.Int("configCategoryId", s.ConfigCategoryId).Required().Errors(),
		)
	case "Update":
		return validatex.Check(
			fields.Int("id", s.ID).Required().Gt(0).Errors(),
			fields.Str("key", s.Key).Required().Errors(),
			fields.Str("name", s.Name).Required().Errors(),
			fields.Str("defaultValue", s.DefaultValue).Required().Errors(),
			fields.Str("optionValue", s.OptionValue).Required().Errors(),
			fields.Int("type", s.Type).Required().Errors(),
			fields.Int("configCategoryId", s.ConfigCategoryId).Required().Errors(),
		)
	case "Detail", "Delete":
		return fields.Int("id", s.ID).Required().Gt(0).Errors()
	default:
		return nil
	}
}

// ValidateWithScene 系统配置批量更新验证
func (s SystemConfigUpdates) ValidateWithScene(scene string, fields validatex.Fields) validatex.Errors {
	switch scene {
	case "UpdateConfig":
		return validatex.Check(
			fields.Slice[struct{}]("list", make([]struct{}, len(s.List))).Required().Errors(),
			fields.SliceStruct("list", s.List),
		)
	default:
		return nil
	}
}

// Messages 验证器错误消息
func (s SystemConfig) Messages() validatex.Messages {
	return commonMessages()
}

// Messages 验证器错误消息
func (s SystemConfigUpdates) Messages() validatex.Messages {
	return commonMessages()
}

// Messages 验证器错误消息
func (s SystemConfigValueUpdate) Messages() validatex.Messages {
	return commonMessages()
}

// Translates 字段翻译
func (s SystemConfig) Translates() validatex.Translates {
	return validatex.Translates{
		"id":               "id",
		"key":              "标识",
		"name":             "名称",
		"defaultValue":     "默认值",
		"optionValue":      "可选值",
		"type":             "配置类型",
		"configCategoryId": "配置分类Id",
		"page":             "页码",
		"pageSize":         "每页数量",
	}
}

// Translates 字段翻译
func (s SystemConfigUpdates) Translates() validatex.Translates {
	return validatex.Translates{"list": "配置列表"}
}

// Translates 字段翻译
func (s SystemConfigValueUpdate) Translates() validatex.Translates {
	return validatex.Translates{
		"id":           "id",
		"key":          "标识",
		"defaultValue": "默认值",
	}
}
