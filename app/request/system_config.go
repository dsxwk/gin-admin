package request

import (
	"errors"
	"fmt"
	"gin/common/base"
	"github.com/gookit/validate"
)

// SystemConfig 请求验证
type SystemConfig struct {
	base.BaseRequest
	ID               int64                     `json:"id" form:"id" validate:"required|int|gt:0" label:"id"`
	Key              string                    `json:"key" form:"key" validate:"required" label:"标识"`
	Name             string                    `json:"name" form:"name" validate:"required" label:"名称"`
	DefaultValue     string                    `json:"defaultValue" form:"defaultValue" validate:"required" label:"默认值"`
	OptionValue      string                    `json:"optionValue" form:"optionValue" validate:"required" label:"可选值"`
	Type             int64                     `json:"type" form:"type" validate:"required|int" label:"配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件"`
	ConfigCategoryId int64                     `json:"configCategoryId" form:"configCategoryId" validate:"required|int" label:"配置分类Id"`
	List             []SystemConfigValueUpdate `json:"list" validate:"required" label:"配置列表"`
	PageListValidate
}

// SystemConfigCreate 系统配置创建
type SystemConfigCreate struct {
	Key              string `json:"key" form:"key" validate:"required" label:"标识"`
	Name             string `json:"name" form:"name" validate:"required" label:"名称"`
	DefaultValue     string `json:"defaultValue" form:"defaultValue" validate:"required" label:"默认值"`
	OptionValue      string `json:"optionValue" form:"optionValue" validate:"required" label:"可选值"`
	Type             int64  `json:"type" form:"type" validate:"required|int" label:"配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件"`
	ConfigCategoryId int64  `json:"configCategoryId" form:"configCategoryId" validate:"required|int" label:"配置分类Id"`
}

// SystemConfigUpdate 系统配置更新
type SystemConfigUpdate struct {
	Key              string `json:"key" form:"key" validate:"required" label:"标识"`
	Name             string `json:"name" form:"name" validate:"required" label:"名称"`
	DefaultValue     string `json:"defaultValue" form:"defaultValue" validate:"required" label:"默认值"`
	OptionValue      string `json:"optionValue" form:"optionValue" validate:"required" label:"可选值"`
	Type             int64  `json:"type" form:"type" validate:"required|int" label:"配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件"`
	ConfigCategoryId int64  `json:"configCategoryId" form:"configCategoryId" validate:"required|int" label:"配置分类Id"`
}

// SystemConfigValueUpdate 系统配置批量更新
type SystemConfigValueUpdate struct {
	ID           int64  `json:"id" validate:"required|int|gt:0"`
	Key          string `json:"key" validate:""`
	DefaultValue string `json:"defaultValue" validate:""`
}

// SystemConfigUpdates 系统配置批量更新验证
type SystemConfigUpdates struct {
	List []SystemConfigValueUpdate `json:"list" validate:"required" label:"配置列表"`
}

// Validate 请求验证
func (s SystemConfig) Validate(data SystemConfig, scene string) error {
	v := validate.Struct(data, scene)
	if !v.Validate(scene) {
		return errors.New(v.Errors.One())
	}
	return nil
}

// Validate 系统配置批量更新请求验证
func (s SystemConfigUpdates) Validate() error {
	if len(s.List) == 0 {
		return errors.New("配置列表不能为空")
	}

	for i, item := range s.List {
		v := validate.Struct(item)
		if !v.Validate() {
			return fmt.Errorf("list[%d]项 %s", i, v.Errors.One())
		}
	}

	return nil
}

// Translates 字段翻译
func (s SystemConfigValueUpdate) Translates() map[string]string {
	return validate.MS{
		"ID":               "id",
		"Key":              "标识",
		"Name":             "名称",
		"DefaultValue":     "默认值",
		"OptionValue":      "可选值",
		"Type":             "配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件",
		"ConfigCategoryId": "配置分类Id",
	}
}

// Messages 验证器错误消息
func (s SystemConfigValueUpdate) Messages() map[string]string {
	return validate.MS{
		"required": "字段 {field} 必填",
		"int":      "字段 {field} 必须为整数",
		"gt":       "字段 {field} 必须大于 0",
	}
}

// ConfigValidation 配置验证
// - 定义验证场景
// - 也可以添加验证设置
func (s SystemConfig) ConfigValidation(v *validate.Validation) {
	scenes := validate.SValues{
		"List": []string{"PageListValidate.Page", "PageListValidate.PageSize"},
		"Create": []string{
			"Key",
			"Name",
			"DefaultValue",
			"OptionValue",
			"Type",
			"ConfigCategoryId",
		},
		"Update": []string{
			"ID",
			"Key",
			"Name",
			"DefaultValue",
			"OptionValue",
			"Type",
			"ConfigCategoryId",
		},
		"Detail": []string{"ID"},
		"Delete": []string{"ID"},
	}
	v.WithScenes(scenes)
}

// Messages 验证器错误消息
func (s SystemConfig) Messages() map[string]string {
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
func (s SystemConfig) Translates() map[string]string {
	return validate.MS{
		"ID":                        "id",
		"Key":                       "标识",
		"Name":                      "名称",
		"DefaultValue":              "默认值",
		"OptionValue":               "可选值",
		"Type":                      "配置类型 1=输入框 2=单选 3=复选 4=下拉菜单 5=文本域 6=附件",
		"ConfigCategoryId":          "配置分类Id",
		"PageListValidate.Page":     "页码",
		"PageListValidate.PageSize": "每页数量",
	}
}
